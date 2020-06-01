package domria

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom"
	"io/ioutil"
	"net/http"
	"rampart/pkg/mining"
	"rampart/pkg/mining/configs"
	"strings"
	"time"
)

func newFetcher(userAgent string, config *configs.Fetcher) *fetcher {
	return &fetcher{
		userAgent,
		&http.Client{Timeout: config.Timeout},
		0,
		config.Portion,
		config.Flags,
		config.SearchURL,
		config.OriginURLPrefix,
		config.ImageURLPrefix,
		config.StateEnding,
		config.StateSuffix,
		config.DistrictLabel,
		config.DistrictEnding,
		config.DistrictSuffix,
	}
}

type fetcher struct {
	userAgent       string
	client          *http.Client
	page            int
	portion         int
	flags           map[mining.Housing]string
	searchURL       string
	originURLPrefix string
	imageURLPrefix  string
	stateEnding     string
	stateSuffix     string
	districtLabel   string
	districtEnding  string
	districtSuffix  string
}

func (fetcher *fetcher) fetchFlats(housing mining.Housing) ([]*flat, error) {
	flag, ok := fetcher.flags[housing]
	if !ok {
		return nil, fmt.Errorf("domria: %v housing isn't acceptable", housing)
	}
	bytes, err := fetcher.getSearch(flag)
	if err != nil {
		return nil, err
	}
	flats, err := fetcher.unmarshalSearch(bytes, housing)
	if err != nil {
		return nil, err
	}
	if length := len(flats); length > 0 {
		log.Debugf("domria: %s housing fetcher on %d page fetched %d flats", housing, fetcher.page, length)
		fetcher.page++
	} else {
		log.Debugf("domria: %s housing fetcher on %d page reset", housing, fetcher.page)
		fetcher.page = 0
	}
	return flats, nil
}

func (fetcher *fetcher) getSearch(flag string) ([]byte, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(fetcher.searchURL, flag, fetcher.page, fetcher.portion),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("domria: failed to construct a request, %v", err)
	}
	request.Header.Set("User-Agent", fetcher.userAgent)
	response, err := fetcher.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("domria: failed to perform a request, %v", err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("domria: failed to read the response body, %v", err)
	}
	if err = response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: failed to close the response body, %v", err)
	}
	return bytes, nil
}

func (fetcher *fetcher) unmarshalSearch(bytes []byte, housing mining.Housing) ([]*flat, error) {
	var search search
	if err := json.Unmarshal(bytes, &search); err != nil {
		return nil, fmt.Errorf("domria: failed to unmarshal the search, %v", err)
	}
	flats := make([]*flat, len(search.Items))
	for i, item := range search.Items {
		originURL := item.BeautifulURL
		if originURL != "" {
			originURL = fetcher.originURLPrefix + originURL
		}
		imageURL := item.MainPhoto
		if imageURL != "" {
			imageURL = fetcher.imageURLPrefix + imageURL
		}
		var point *geom.Point
		if item.Longitude != 0 && item.Latitude != 0 {
			point = geom.NewPointFlat(geom.XY, []float64{float64(item.Longitude), float64(item.Latitude)})
		}
		state := item.StateNameUK
		if state != "" && strings.HasSuffix(state, fetcher.stateEnding) {
			state += fetcher.stateSuffix
		}
		district := item.DistrictNameUK
		if district != "" &&
			item.DistrictTypeName == fetcher.districtLabel &&
			strings.HasSuffix(district, fetcher.districtEnding) {
			district += fetcher.districtSuffix
		}
		street := item.StreetNameUK
		if street == "" && item.StreetName != "" {
			street = item.StreetName
		}
		flats[i] = &flat{
			originURL,
			imageURL,
			(*time.Time)(item.UpdatedAt),
			float64(item.PriceArr.USD),
			item.TotalSquareMeters,
			item.LivingSquareMeters,
			item.KitchenSquareMeters,
			item.RoomsCount,
			item.Floor,
			item.FloorsCount,
			housing,
			item.UserNewbuildNameUK,
			point,
			state,
			item.CityNameUK,
			district,
			street,
			item.BuildingNumberStr,
		}
	}
	return flats, nil
}
