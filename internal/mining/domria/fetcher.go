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

func NewFetcher(config config.Fetcher, drain *metrics.Drain) *Fetcher {
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
		drain,
	}
}

type Fetcher struct {
	client    *http.Client
	page      int
	portion   int
	flags     map[string]string
	headers   misc.Headers
	searchURL string
	drain     *metrics.Drain
}

func (fetcher *Fetcher) FetchFlats(housing string) ([]Flat, error) {
	flag, ok := fetcher.flags[housing]
	if !ok {
		return nil, fmt.Errorf("domria: fetcher doesn't accept housing %s", housing)
	}
	start := time.Now()
	search, err := fetcher.getSearch(flag)
	fetcher.drain.GatherFetchingDuration(start)
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

func (fetcher *Fetcher) getSearch(flag string) (search, error) {
	s := search{}
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(fetcher.searchURL, flag, fetcher.page, fetcher.portion),
		nil,
	)
	if err != nil {
		return s, fmt.Errorf("domria: fetcher failed to construct a request, %v", err)
	}
	fetcher.headers.Inject(request)
	response, err := fetcher.client.Do(request)
	if err != nil {
		return s, fmt.Errorf("domria: fetcher failed to perform a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return s, fmt.Errorf("domria: fetcher got response with status %s", response.Status)
	}
	if err := json.NewDecoder(response.Body).Decode(&s); err != nil {
		_ = response.Body.Close()
		return s, fmt.Errorf("domria: fetcher failed to unmarshal the search, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return s, fmt.Errorf("domria: fetcher failed to close the response body, %v", err)
	}
	return s, nil
}

func (fetcher *Fetcher) getFlats(s search, housing string) []Flat {
	flats := make([]Flat, len(s.Items))
	for j, i := range s.Items {
		street := i.StreetNameUK
		if street == "" && i.StreetName != "" {
			street = i.StreetName
		}
		flats[j] = Flat{
			Source:      i.Source,
			OriginURL:   i.BeautifulURL,
			ImageURL:    i.MainPhoto,
			MediaCount:  len(i.Photos) + len(i.Panoramas),
			UpdateTime:  time.Time(i.UpdatedAt),
			IsInspected: i.Inspected == 1,
			Price:       float64(i.PriceArr.USD),
			TotalArea:   i.TotalSquareMeters,
			LivingArea:  i.LivingSquareMeters,
			KitchenArea: i.KitchenSquareMeters,
			RoomNumber:  i.RoomsCount,
			Floor:       i.Floor,
			TotalFloor:  i.FloorsCount,
			Housing:     housing,
			Complex:     i.UserNewbuildNameUK,
			Point:       orb.Point{float64(i.Longitude), float64(i.Latitude)},
			State:       i.StateNameUK,
			City:        i.CityNameUK,
			District:    i.DistrictNameUK,
			Street:      street,
			HouseNumber: i.BuildingNumberStr,
		}
	}
	return flats
}
