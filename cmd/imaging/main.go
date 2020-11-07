package main

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io"
	"net/http"
	"os"
	"sync"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "imaging")
	cfg, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	client := &http.Client{Timeout: cfg.Imaging.Timeout}
	records := make(chan []string, cfg.Imaging.ThreadNumber)
	group := &sync.WaitGroup{}
	group.Add(cfg.Imaging.ThreadNumber)
	for i := 0; i < cfg.Imaging.ThreadNumber; i++ {
		go work(client, records, group, cfg.Imaging, entry)
	}
	file, err := os.Open(cfg.Imaging.InputPath)
	if err != nil {
		entry.Fatalf("main: imaging failed to open the input file, %v", err)
	}
	err = run(file, records)
	group.Wait()
	if err != nil {
		_ = file.Close()
		entry.Fatal(err)
	}
	if err := file.Close(); err != nil {
		entry.Fatalf("main: imaging failed to close the input file, %v", err)
	}
}

func work(
	client *http.Client,
	records <-chan []string,
	group *sync.WaitGroup,
	config config.Imaging,
	logger log.FieldLogger,
) {
	for record := range records {
		if err := load(client, config, record); err != nil {
			logger.WithField("url", record[0]).Error(err)
		}
	}
	group.Done()
}

func load(client *http.Client, config config.Imaging, record []string) error {
	request, err := http.NewRequest(http.MethodGet, record[0], nil)
	if err != nil {
		return fmt.Errorf("main: imaging failed to make a request, %v", err)
	}
	config.Headers.Inject(request)
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("main: imaging failed to send a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return fmt.Errorf("main: imaging got a non-ok status %s", response.Status)
	}
	file, err := os.Create(
		misc.ResolvePath(
			fmt.Sprintf(config.OutputFormat, sha1.Sum([]byte(record[0])), record[1]),
		),
	)
	if err != nil {
		_ = response.Body.Close()
		return fmt.Errorf("main: imaging failed to create the file, %v", err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		_ = file.Close()
		_ = response.Body.Close()
		return fmt.Errorf("main: imaging failed to copy, %v", err)
	}
	if err := file.Close(); err != nil {
		_ = response.Body.Close()
		return fmt.Errorf("main: imaging failed to close the file, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("main: imaging failed to close the response body, %v", err)
	}
	return nil
}

func run(input io.Reader, records chan<- []string) error {
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
			close(records)
			return nil
		} else if err != nil {
			close(records)
			return fmt.Errorf("main: imaging failed to read a row of the input file, %v", err)
		}
		records <- record
	}
}
