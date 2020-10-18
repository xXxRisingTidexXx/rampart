package domria

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewFetcher(
	config config.Fetcher,
	drain *metrics.Drain,
	logger log.FieldLogger,
) *Fetcher {
	flags := make(map[misc.Housing]string, len(config.Flags))
	for key, value := range config.Flags {
		flags[key] = value
	}
	return &Fetcher{
		&http.Client{Timeout: config.Timeout},
		0,
		config.Portion,
		flags,
		config.Headers,
		config.SearchFormat,
		drain,
		logger,
	}
}

type Fetcher struct {
	client       *http.Client
	page         int
	portion      int
	flags        map[misc.Housing]string
	headers      misc.Headers
	searchFormat string
	drain        *metrics.Drain
	logger       log.FieldLogger
}

func (fetcher *Fetcher) FetchFlats(housing misc.Housing) []Flat {
	flag, ok := fetcher.flags[housing]
	if !ok {
		fetcher.logger.WithField("housing", housing).Error(
			"domria: fetcher doesn't accept the housing",
		)
		return make([]Flat, 0)
	}
	start := time.Now()
	search, err := fetcher.getSearch(flag)
	fetcher.drain.DrainDuration(metrics.FetchingDuration, start)
	if err != nil {
		fetcher.drain.DrainNumber(metrics.FailedFetchingNumber)
		fetcher.logger.Error(err)
		return make([]Flat, 0)
	}
	flats := fetcher.getFlats(search, housing)
	if len(flats) > 0 {
		fetcher.page++
	} else {
		fetcher.page = 0
	}
	return flats
}

func (fetcher *Fetcher) getSearch(flag string) (search, error) {
	var s search
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(fetcher.searchFormat, flag, fetcher.page, fetcher.portion),
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

func (fetcher *Fetcher) getFlats(s search, housing misc.Housing) []Flat {
	flats := make([]Flat, len(s.Items))
	for j, i := range s.Items {
		photos := make([]string, 0, len(i.Photos))
		for id := range i.Photos {
			photos = append(photos, id)
		}
		panoramas := make([]string, len(i.Panoramas))
		for k := range i.Panoramas {
			panoramas[k] = i.Panoramas[k].Img
		}
		street := i.StreetNameUK
		if street == "" && i.StreetName != "" {
			street = i.StreetName
		}
		flats[j] = Flat{
			Source:      i.Source,
			URL:         i.BeautifulURL,
			Photos:      photos,
			Panoramas:   panoramas,
			UpdateTime:  time.Time(i.UpdatedAt),
			IsSold:      i.SaleDate != "",
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
