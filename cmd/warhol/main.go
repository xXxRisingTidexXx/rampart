package main

import (
	gobytes "bytes"
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/imaging"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "warhol")
	c, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	file, err := os.Open(c.Warhol.InputPath)
	if err != nil {
		entry.Fatalf("main: warhol failed to open the input file, %v", err)
	}
	records := make(chan []string, c.Warhol.ThreadNumber)
	raws := make(chan imaging.Raw, (c.Warhol.ThreadNumber+1)*len(imaging.Effects))
	dumpCount := c.Warhol.ThreadNumber + runtime.NumCPU()<<1
	assets := make(chan imaging.Asset, dumpCount)
	client := &http.Client{Timeout: c.Warhol.Timeout}
	loadGroup := &sync.WaitGroup{}
	loadGroup.Add(c.Warhol.ThreadNumber)
	for i := 0; i < c.Warhol.ThreadNumber; i++ {
		go load(c.Warhol, records, raws, assets, client, entry, loadGroup)
	}
	processGroup := &sync.WaitGroup{}
	processGroup.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go process(raws, assets, entry, processGroup)
	}
	dumpGroup := &sync.WaitGroup{}
	dumpGroup.Add(dumpCount)
	for i := 0; i < dumpCount; i++ {
		go dump(assets, c.Warhol.OutputFormat, entry, dumpGroup)
	}
	err = read(file, records)
	close(records)
	loadGroup.Wait()
	close(raws)
	processGroup.Wait()
	close(assets)
	dumpGroup.Wait()
	if err != nil {
		_ = file.Close()
		entry.Fatal(err)
	}
	if err := file.Close(); err != nil {
		entry.Fatalf("main: warhol failed to close the input file, %v", err)
	}
}

func load(
	config config.Warhol,
	records <-chan []string,
	raws chan<- imaging.Raw,
	assets chan<- imaging.Asset,
	client *http.Client,
	logger log.FieldLogger,
	group *sync.WaitGroup,
) {
	for record := range records {
		for retry, err := 1, io.EOF; retry <= config.RetryLimit && err != nil; retry++ {
			if err = pipe(record, client, config.UserAgent, raws, assets); err != nil {
				logger.WithFields(log.Fields{"url": record[0], "retry": retry}).Error(err)
			}
		}
	}
	group.Done()
}

func pipe(
	record []string,
	client *http.Client,
	userAgent string,
	raws chan<- imaging.Raw,
	assets chan<- imaging.Asset,
) error {
	request, err := http.NewRequest(http.MethodGet, record[0], nil)
	if err != nil {
		return fmt.Errorf("main: warhol failed to make a request, %v", err)
	}
	request.Header.Set("User-Agent", userAgent)
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("main: warhol failed to send a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return fmt.Errorf("main: warhol got a non-ok status %s", response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return fmt.Errorf("main: warhol failed to read the response body, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("main: warhol failed to close the response body, %v", err)
	}
	source, _, err := image.Decode(gobytes.NewBuffer(bytes))
	if err != nil {
		return fmt.Errorf("main: warhol failed to decode the source, %v", err)
	}
	sum := sha1.Sum([]byte(record[0]))
	hash := hex.EncodeToString(sum[:])
	assets <- imaging.Asset{
		Hash:   hash,
		Group:  record[1],
		Label:  record[2],
		Effect: "origin",
		Bytes:  bytes,
	}
	for _, effect := range imaging.Effects {
		raws <- imaging.Raw{
			Hash:   hash,
			Group:  record[1],
			Label:  record[2],
			Effect: effect,
			Source: source,
		}
	}
	return nil
}

func process(
	raws <-chan imaging.Raw,
	assets chan<- imaging.Asset,
	logger log.FieldLogger,
	group *sync.WaitGroup,
) {
	for raw := range raws {
		bytes, err := raw.Effect.Apply(raw.Source)
		if err != nil {
			fields := log.Fields{
				"hash":   raw.Hash,
				"effect": raw.Effect.Name(),
				"group":  raw.Group,
				"label":  raw.Label,
			}
			logger.WithFields(fields).Error(err)
		} else {
			assets <- imaging.Asset{
				Hash:   raw.Hash,
				Effect: raw.Effect.Name(),
				Group:  raw.Group,
				Label:  raw.Label,
				Bytes:  bytes,
			}
		}
	}
	group.Done()
}

func dump(
	assets <-chan imaging.Asset,
	format string,
	logger log.FieldLogger,
	group *sync.WaitGroup,
) {
	for asset := range assets {
		err := ioutil.WriteFile(
			misc.ResolvePath(
				fmt.Sprintf(format, asset.Hash, asset.Effect, asset.Group, asset.Label),
			),
			asset.Bytes,
			0644,
		)
		if err != nil {
			fields := log.Fields{
				"hash":   asset.Hash,
				"effect": asset.Effect,
				"group":  asset.Group,
				"label":  asset.Label,
			}
			logger.WithFields(fields).Errorf("main: warhol failed to write the file, %v", err)
		}
	}
	group.Done()
}

func read(input io.Reader, records chan<- []string) error {
	reader := csv.NewReader(input)
	if _, err := reader.Read(); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("main: warhol failed to read header of the input file, %v", err)
	}
	if expected := 3; reader.FieldsPerRecord != expected {
		return fmt.Errorf(
			"main: warhol got invalid field number, %d != %d",
			reader.FieldsPerRecord,
			expected,
		)
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("main: warhol failed to read a row of the input file, %v", err)
		}
		records <- record
	}
}
