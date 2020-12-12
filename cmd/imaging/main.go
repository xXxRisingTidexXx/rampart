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
	entry := log.WithField("app", "imaging")
	c, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	file, err := os.Open(c.Imaging.InputPath)
	if err != nil {
		entry.Fatalf("main: imaging failed to open the input file, %v", err)
	}
	records := make(chan []string, c.Imaging.ThreadNumber)
	raws := make(chan imaging.Raw, (c.Imaging.ThreadNumber+1)*len(imaging.Effects))
	dumpCount := c.Imaging.ThreadNumber + runtime.NumCPU()<<1
	assets := make(chan imaging.Asset, dumpCount)
	client := &http.Client{Timeout: c.Imaging.Timeout}
	loadGroup := &sync.WaitGroup{}
	loadGroup.Add(c.Imaging.ThreadNumber)
	for i := 0; i < c.Imaging.ThreadNumber; i++ {
		go load(records, raws, assets, client, c.Imaging, entry, loadGroup)
	}
	processGroup := &sync.WaitGroup{}
	processGroup.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go process(raws, assets, entry, processGroup)
	}
	dumpGroup := &sync.WaitGroup{}
	dumpGroup.Add(dumpCount)
	for i := 0; i < dumpCount; i++ {
		go dump(assets, c.Imaging.OutputFormat, entry, dumpGroup)
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
		entry.Fatalf("main: imaging failed to close the input file, %v", err)
	}
}

func load(
	records <-chan []string,
	raws chan<- imaging.Raw,
	assets chan<- imaging.Asset,
	client *http.Client,
	config config.Imaging,
	logger log.FieldLogger,
	group *sync.WaitGroup,
) {
	for record := range records {
		for retry, err := 1, io.EOF; retry <= config.RetryLimit && err != nil; retry++ {
			if err = pipe(record, client, config.Headers, raws, assets); err != nil {
				logger.WithFields(log.Fields{"url": record[0], "retry": retry}).Error(err)
			}
		}
	}
	group.Done()
}

func pipe(
	record []string,
	client *http.Client,
	headers misc.Headers,
	raws chan<- imaging.Raw,
	assets chan<- imaging.Asset,
) error {
	request, err := http.NewRequest(http.MethodGet, record[0], nil)
	if err != nil {
		return fmt.Errorf("main: imaging failed to make a request, %v", err)
	}
	headers.Inject(request)
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("main: imaging failed to send a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return fmt.Errorf("main: imaging got a non-ok status %s", response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return fmt.Errorf("main: imaging failed to read the response body, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("main: imaging failed to close the response body, %v", err)
	}
	source, _, err := image.Decode(gobytes.NewBuffer(bytes))
	if err != nil {
		return fmt.Errorf("main: imaging failed to decode the source, %v", err)
	}
	hash := sha1.Sum([]byte(record[0]))
	assets <- imaging.Asset{Hash: hash, Label: record[1], Effect: "origin", Bytes: bytes}
	for _, effect := range imaging.Effects {
		raws <- imaging.Raw{Hash: hash, Label: record[1], Effect: effect, Source: source}
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
				"hash":   hex.EncodeToString(raw.Hash[:]),
				"effect": raw.Effect.Name(),
				"label":  raw.Label,
			}
			logger.WithFields(fields).Error(err)
		} else {
			assets <- imaging.Asset{
				Hash:   raw.Hash,
				Effect: raw.Effect.Name(),
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
			misc.ResolvePath(fmt.Sprintf(format, asset.Hash, asset.Effect, asset.Label)),
			asset.Bytes,
			0644,
		)
		if err != nil {
			fields := log.Fields{
				"hash":   hex.EncodeToString(asset.Hash[:]),
				"effect": asset.Effect,
				"label":  asset.Label,
			}
			logger.WithFields(fields).Errorf("main: imaging failed to write the file, %v", err)
		}
	}
	group.Done()
}

func read(input io.Reader, records chan<- []string) error {
	reader := csv.NewReader(input)
	if _, err := reader.Read(); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("main: imaging failed to read header of the input file, %v", err)
	}
	if reader.FieldsPerRecord != 2 {
		return fmt.Errorf(
			"main: imaging got invalid field number, %d != 2",
			reader.FieldsPerRecord,
		)
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("main: imaging failed to read a row of the input file, %v", err)
		}
		records <- record
	}
}
