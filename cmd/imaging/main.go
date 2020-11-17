package main

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io"
	"io/ioutil"
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
		for retry, err := 1, io.EOF; retry <= config.RetryLimit && err != nil; retry++ {
			if err = load(client, config, record); err != nil {
				logger.WithFields(log.Fields{"url": record[0], "retry": retry}).Error(err)
			}
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
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return fmt.Errorf("main: imaging failed to read the response body, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("main: imaging failed to close the response body, %v", err)
	}
	err = ioutil.WriteFile(
		misc.ResolvePath(
			fmt.Sprintf(config.OutputFormat, sha1.Sum([]byte(record[0])), record[1]),
		),
		bytes,
		0644,
	)
	if err != nil {
		return fmt.Errorf("main: imaging failed to write the file, %v", err)
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
