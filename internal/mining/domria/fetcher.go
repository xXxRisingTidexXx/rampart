package domria

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom"
	"io/ioutil"
	"net/http"
	"rampart/internal/mining"
	"rampart/internal/mining/configs"
	"time"
)

func newFetcher(config *configs.Fetcher) *fetcher {
	return &fetcher{
		&http.Client{Timeout: time.Duration(config.Timeout)},
		0,
		config.Portion,
		config.Flags,
		config.Headers,
		config.SearchURL,
	}
}

type fetcher struct {
	client    *http.Client
	page      int
	portion   int
	flags     map[mining.Housing]string
	headers   map[string]string
	searchURL string
}

func (fetcher *fetcher) fetchFlats(housing mining.Housing) ([]*flat, error) {
	flag, ok := fetcher.flags[housing]
	if !ok {
		return nil, fmt.Errorf("domria: fetcher doesn't accept %v housing", housing)
	}
	start := time.Now()
	bytes, err := fetcher.getSearch(flag)
	if err != nil {
		return nil, err
	}
	duration := time.Since(start).Seconds()
	flats, err := fetcher.unmarshalSearch(bytes, housing)
	if err != nil {
		return nil, err
	}
	if length := len(flats); length > 0 {
		log.Debugf("domria: fetcher on %d page received %d flats (%.3fs)", fetcher.page, length, duration)
		fetcher.page++
	} else {
		log.Debugf("domria: fetcher on %d page reset (%.3fs)", fetcher.page, duration)
		fetcher.page = 0
	}
	return flats, nil
}

func (fetcher *fetcher) getSearch(flag string) (bytes []byte, err error) {
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
	defer func() {
		if closingErr := response.Body.Close(); closingErr != nil {
			closingErr = fmt.Errorf("domria: fetcher failed to close the response body, %v", closingErr)
			if err == nil {
				err = closingErr
			} else {
				log.Error(closingErr)
			}
		}
		if err != nil {
			bytes = nil
		}
	}()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("domria: fetcher got response with status %s", response.Status)
	}
	return ioutil.ReadAll(response.Body)
}

func (fetcher *fetcher) unmarshalSearch(bytes []byte, housing mining.Housing) ([]*flat, error) {
	var search search
	if err := json.Unmarshal(bytes, &search); err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to unmarshal the search, %v", err)
	}
	flats := make([]*flat, len(search.Items))
	for i, item := range search.Items {
		price := 0.0
		if item.PriceArr != nil {
			price = float64(item.PriceArr.USD)
		}
		var point *geom.Point
		if item.Longitude != 0 || item.Latitude != 0 {
			point = geom.NewPointFlat(geom.XY, []float64{float64(item.Longitude), float64(item.Latitude)})
		}
		street := item.StreetNameUK
		if street == "" && item.StreetName != "" {
			street = item.StreetName
		}
		flats[i] = &flat{
			item.BeautifulURL,
			item.MainPhoto,
			(*time.Time)(item.UpdatedAt),
			price,
			item.TotalSquareMeters,
			item.LivingSquareMeters,
			item.KitchenSquareMeters,
			item.RoomsCount,
			item.Floor,
			item.FloorsCount,
			housing,
			item.UserNewbuildNameUK,
			point,
			item.StateNameUK,
			item.CityNameUK,
			item.DistrictNameUK,
			street,
			item.BuildingNumberStr,
		}
	}
	return flats, nil
}
