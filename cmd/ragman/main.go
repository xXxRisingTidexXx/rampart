package main

import (
	gobytes "bytes"
	"encoding/json"
	"fmt"
	gocron "github.com/robfig/cron/v3"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/domria"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"io/ioutil"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bytes, err := ioutil.ReadFile("datasets/outdated.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	urls := map[string]struct{}{}
	for _, line := range gobytes.Split(bytes, []byte{'\n'}) {
		urls[string(line)] = struct{}{}
	}
	cron := gocron.New()
	_, _ = cron.AddJob(
		"0-59/1 * * * *",
		NewRagman("primary", cfg.Mining.DomriaPrimaryMiner.Fetcher, urls),
	)
	_, _ = cron.AddJob(
		"1-59/1 * * * *",
		NewRagman("secondary", cfg.Mining.DomriaSecondaryMiner.Fetcher, urls),
	)
	cron.Run()
}

func NewRagman(housing string, config *config.Fetcher, urls map[string]struct{}) *Ragman {
	return &Ragman{
		housing,
		domria.NewFetcher(config, metrics.NewGatherer(housing, nil)),
		urls,
	}
}

type Ragman struct {
	housing string
	fetcher *domria.Fetcher
	urls    map[string]struct{}
}

func (ragman *Ragman) Run() {
	flats, err := ragman.fetcher.FetchFlats(ragman.housing)
	if err != nil {
		fmt.Println(err)
	} else {
		buffer := &gobytes.Buffer{}
		for _, flat := range flats {
			if _, ok := ragman.urls["https://dom.ria.com/uk/" + flat.OriginURL]; ok {
				if err := json.Indent(buffer, []byte(flat.Source), "", "  "); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(buffer.String())
				}
			}
		}
	}
}
