package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func NewFetcher(portion int, timeout time.Duration) *Fetcher {
	return &Fetcher{
		portion,
		&http.Client{Timeout: timeout},
		map[Housing]string{Primary: "newbuildings=1", Secondary: "secondary=1"},
		"prospector/1.0 (rampart/prospector)",
		"https://dom.ria.com/node/api/autocompleteCities?langId=4",
		"https://dom.ria.com/searchEngine/?category=1&realty_type=2&operation_type=1" +
			"&fullCategoryOperation=1_2_1&%s&page=%d&state_id=%d&city_id=%d&limit=%d",
	}
}

type Fetcher struct {
	portion      int
	client       *http.Client
	housingFlags map[Housing]string
	userAgent    string
	localityURL  string
	flatURL      string
}

func (fetcher *Fetcher) Fetch(state, city string, housing Housing) ([]map[string]interface{}, error) {
	localities := make([]*Locality, 24)
	if err := fetcher.fetchJSON(fetcher.localityURL, &localities); err != nil {
		return nil, err
	}
	stateID, cityID := -1, -1
	for _, locality := range localities {
		if locality.State == state && locality.City == city {
			stateID, cityID = locality.Payload.StateID, locality.Payload.CityID
			break
		}
	}
	if stateID == -1 || cityID == -1 {
		return nil, fmt.Errorf("couldn't find locality with state %s and city %s", state, city)
	}
	log.Info(stateID, cityID)
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

type Locality struct {
	State   string   `json:"stateName"`
	City    string   `json:"cityName"`
	Payload *Payload `json:"payload"`
}

func (locality *Locality) String() string {
	return fmt.Sprintf("{%s %s %v}", locality.State, locality.City, locality.Payload)
}

type Payload struct {
	StateID int `json:"stateId"`
	CityID  int `json:"cityId"`
}

func (payload *Payload) String() string {
	return fmt.Sprintf("{%d %d}", payload.StateID, payload.CityID)
}

func main() {
	log.SetLevel(log.InfoLevel)
	log.Info("prospector started")
	fetcher := NewFetcher(5, 5*time.Second)
	_, err := fetcher.Fetch("Київська", "Київ", Primary)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("prospector finished")
}
