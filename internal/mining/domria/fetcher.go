package domria

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
	"net/http"
	"time"
)

func NewFetcher(config *config.Fetcher, gatherer *metrics.Gatherer) *Fetcher {
	return &Fetcher{
		&http.Client{Timeout: time.Duration(config.Timeout)},
		0,
		config.Portion,
		config.Flags,
		config.Headers,
		config.SearchURL,
		gatherer,
	}
}

type Fetcher struct {
	client    *http.Client
	page      int
	portion   int
	flags     map[misc.Housing]string
	headers   map[string]string
	searchURL string
	gatherer  *metrics.Gatherer
}

func (fetcher *Fetcher) FetchFlats(housing misc.Housing) ([]*Flat, error) {
	flag, ok := fetcher.flags[housing]
	if !ok {
		return nil, fmt.Errorf("domria: fetcher doesn't accept %v housing", housing)
	}
	start := time.Now()
	bytes, err := fetcher.getSearch(flag)
	fetcher.gatherer.GatherFetchingDuration(start)
	if err != nil {
		return nil, err
	}
	flats, err := fetcher.unmarshalSearch(bytes, housing)
	if err != nil {
		return nil, err
	}
	if len(flats) > 0 {
		fetcher.page++
	} else {
		fetcher.page = 0
	}
	return flats, nil
}

func (fetcher *Fetcher) getSearch(flag string) ([]byte, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(fetcher.searchURL, flag, fetcher.page, fetcher.portion),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to construct a request, %v", err)
	}
	for key, value := range fetcher.headers {
		request.Header.Set(key, value)
	}
	response, err := fetcher.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to perform a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: fetcher got response with status %s", response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: fetcher failed to read the response body, %v", err)
	}
	if err = response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to close the response body, %v", err)
	}
	return bytes, nil
}

func (fetcher *Fetcher) unmarshalSearch(bytes []byte, housing misc.Housing) ([]*Flat, error) {
	var search search
	if err := json.Unmarshal(bytes, &search); err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to unmarshal the search, %v", err)
	}
	flats := make([]*Flat, len(search.Items))
	for i, item := range search.Items {
		price := 0.0
		if item.PriceArr != nil {
			price = float64(item.PriceArr.USD)
		}
		street := item.StreetNameUK
		if street == "" && item.StreetName != "" {
			street = item.StreetName
		}
		flats[i] = &Flat{
			item.BeautifulURL,
			item.MainPhoto,
			time.Time(item.UpdatedAt),
			price,
			item.TotalSquareMeters,
			item.LivingSquareMeters,
			item.KitchenSquareMeters,
			item.RoomsCount,
			item.Floor,
			item.FloorsCount,
			housing,
			item.UserNewbuildNameUK,
			orb.Point{float64(item.Longitude), float64(item.Latitude)},
			item.StateNameUK,
			item.CityNameUK,
			item.DistrictNameUK,
			street,
			item.BuildingNumberStr,
		}
	}
	return flats, nil
}
