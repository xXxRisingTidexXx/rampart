package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type Locality struct {
	City    string   `json:"cityName"`
	State   string   `json:"stateName"`
	Payload *Payload `json:"payload"`
}

func (locality *Locality) String() string {
	return fmt.Sprintf("{%s %s %v}", locality.City, locality.State, locality.Payload)
}

type Payload struct {
	CityID  int `json:"cityId"`
	StateID int `json:"stateId"`
}

func (payload *Payload) String() string {
	return fmt.Sprintf("{%d %d}", payload.CityID, payload.StateID)
}

type Housing int

const (
	Primary Housing = iota
	Secondary
)

func NewFetcher(timeout time.Duration) *Fetcher {
	return &Fetcher{
		&http.Client{Timeout: timeout},
		map[Housing]string{Primary: "newbuildings=1", Secondary: "secondary=1"},
		"prospector/1.0 (rampart/prospector)",
		"https://dom.ria.com/node/api/autocompleteCities?langId=4",
		"https://dom.ria.com/searchEngine/?category=1&realty_type=2&operation_type=1" +
			"&fullCategoryOperation=1_2_1&%s&page=%d&state_id=%d&city_id=%d&limit=%d",
	}
}

type Fetcher struct {
	client       *http.Client
	housingFlags map[Housing]string
	userAgent    string
	localityURL  string
	flatURL      string
}

func (fetcher *Fetcher) Fetch() ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (fetcher *Fetcher) fetchJSON(url string, target interface{}) error {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", fetcher.userAgent)
	response, err := fetcher.client.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, target); err != nil {
		return err
	}
	if err := response.Body.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	log.SetLevel(log.InfoLevel)
	log.Info("prospector started")

	citiesURL := "https://dom.ria.com/node/api/autocompleteCities?langId=4"
	client := http.Client{Timeout: time.Second * 2}
	request, err := http.NewRequest(http.MethodGet, citiesURL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("User-Agent", "prospector/1.0 (rampart/prospector)")
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	localities := make([]*Locality, 24)
	if err := json.Unmarshal(body, &localities); err != nil {
		log.Fatalln(err)
	}
	if err := response.Body.Close(); err != nil {
		log.Fatalln(err)
	}
	log.Info(localities)
	log.Info("prospector finished")
}
