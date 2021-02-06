package domria

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io"
	"io/ioutil"
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
		config.RetryLimit,
		config.Portion,
		flags,
		config.SearchFormat,
		drain,
		logger,
	}
}

type Fetcher struct {
	client       *http.Client
	page         int
	retryLimit   int
	portion      int
	flags        map[misc.Housing]string
	searchFormat string
	drain        *metrics.Drain
	logger       log.FieldLogger
}

func (f *Fetcher) FetchFlats(housing misc.Housing) []Flat {
	f.drain.DrainPage(f.page)
	flag, ok := f.flags[housing]
	if !ok {
		f.logger.WithField("housing", housing).Error("domria: fetcher doesn't accept the housing")
		return make([]Flat, 0)
	}
	entry := f.logger.WithField("page", f.page)
	start := time.Now()
	search, err := f.getSearch(flag, entry)
	f.drain.DrainDuration(metrics.FetchingDuration, start)
	if err != nil {
		f.drain.DrainNumber(metrics.FailedFetchingNumber)
		entry.Error(err)
		return make([]Flat, 0)
	}
	flats := f.getFlats(search, housing)
	if len(flats) > 0 {
		f.page++
	} else {
		f.page = 0
	}
	return flats
}

func (f *Fetcher) getSearch(flag string, logger log.FieldLogger) (search, error) {
	url := fmt.Sprintf(f.searchFormat, flag, f.page, f.portion)
	bytes, err := make([]byte, 0), io.EOF
	for retry := 1; retry <= f.retryLimit && err != nil; retry++ {
		if bytes, err = f.trySearch(url); err != nil {
			logger.WithField("retry", retry).Error(err)
		}
	}
	var s search
	if err != nil {
		return s, fmt.Errorf("domria: fetcher exhausted retry limit")
	}
	if err := json.Unmarshal(bytes, &s); err != nil {
		return s, fmt.Errorf("domria: fetcher failed to unmarshal the search, %v", err)
	}
	return s, nil
}

func (f *Fetcher) trySearch(url string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to construct a request, %v", err)
	}
	request.Header.Set("User-Agent", misc.UserAgent)
	response, err := f.client.Do(request)
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
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to close the response body, %v", err)
	}
	return bytes, nil
}

func (f *Fetcher) getFlats(s search, housing misc.Housing) []Flat {
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
			HouseNumber: string(i.BuildingNumberStr),
		}
	}
	return flats
}
