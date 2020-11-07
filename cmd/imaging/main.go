package main

import (
	"encoding/csv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"io"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "imaging")
	cfg, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	file, err := os.Open(cfg.Imaging.InputPath)
	if err != nil {
		entry.Fatalf("main: imaging failed to open the input file, %v", err)
	}
	if err := run(file); err != nil {
		entry.Error(err)
	}
	if err := file.Close(); err != nil {
		entry.Fatalf("main: imaging failed to close the input file, %v", err)
	}
}

func run(input io.Reader) error {
	reader := csv.NewReader(input)
	if _, err := reader.Read(); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("main: imaging failed to read header of the input file, %v", err)
	}
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("main: imaging failed to read a row of the input file, %v", err)
		}
		fmt.Println(fields)
	}
}
