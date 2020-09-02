package domria

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewFetcher(config *config.Fetcher, gatherer *metrics.Gatherer) *Fetcher {
	flags := make(map[string]string, len(config.Flags))
	for key, value := range config.Flags {
		flags[string(key)] = value
	}
	return &Fetcher{
		&http.Client{Timeout: config.Timeout},
		0,
		config.Portion,
		flags,
		config.Headers,
		config.SearchURL,
		gatherer,
	}
}

type Fetcher struct {
	client    *http.Client
	page      int
	portion   int
	flags     map[string]string
	headers   misc.Headers
	searchURL string
	gatherer  *metrics.Gatherer
}

func (fetcher *Fetcher) FetchFlats(housing string) ([]*Flat, error) {
	flag, ok := fetcher.flags[housing]
	if !ok {
		return nil, fmt.Errorf("domria: fetcher doesn't accept housing %s", housing)
	}
	start := time.Now()
	search, err := fetcher.getSearch(flag)
	fetcher.gatherer.GatherFetchingDuration(start)
	if err != nil {
		return nil, err
	}
	flats := fetcher.getFlats(search, housing)
	if len(flats) > 0 {
		fetcher.page++
	} else {
		fetcher.page = 0
	}
	return flats, nil
}

func (fetcher *Fetcher) getSearch(flag string) (*search, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(fetcher.searchURL, flag, fetcher.page, fetcher.portion),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to construct a request, %v", err)
	}
	fetcher.headers.Inject(request)
	response, err := fetcher.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to perform a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: fetcher got response with status %s", response.Status)
	}
	search := search{}
	if err := json.NewDecoder(response.Body).Decode(&search); err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: fetcher failed to unmarshal the search, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to close the response body, %v", err)
	}
	return &search, nil
}

func (fetcher *Fetcher) getFlats(search *search, housing string) []*Flat {
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
			OriginURL:   item.BeautifulURL,
			ImageURL:    item.MainPhoto,
			MediaCount:  len(item.Photos) + len(item.Panoramas),
			UpdateTime:  time.Time(item.UpdatedAt),
			IsInspected: item.Inspected == 1,
			Price:       price,
			TotalArea:   item.TotalSquareMeters,
			LivingArea:  item.LivingSquareMeters,
			KitchenArea: item.KitchenSquareMeters,
			RoomNumber:  item.RoomsCount,
			Floor:       item.Floor,
			TotalFloor:  item.FloorsCount,
			Housing:     housing,
			Complex:     item.UserNewbuildNameUK,
			Point:       orb.Point{float64(item.Longitude), float64(item.Latitude)},
			State:       item.StateNameUK,
			City:        item.CityNameUK,
			District:    item.DistrictNameUK,
			Street:      street,
			HouseNumber: item.BuildingNumberStr,
			Source:      item.Source,
		}
	}
	return flats
}
