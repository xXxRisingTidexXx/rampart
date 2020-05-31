package domria

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
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
		config.OriginURL,
		config.ImageURL,
		config.USDLabel,
		config.MinTotalArea,
		config.MaxTotalArea,
		config.MinLivingArea,
		config.MinKitchenArea,
		config.MinRoomNumber,
		config.MaxRoomNumber,
		config.MinSpecificArea,
		config.MaxSpecificArea,
		config.MinFloor,
		config.MinTotalFloor,
		config.MaxTotalFloor,
		config.MinLongitude,
		config.MaxLongitude,
		config.MinLatitude,
		config.MaxLatitude,
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
	originURL       string
	imageURL        string
	usdLabel        string
	minTotalArea    float64
	maxTotalArea    float64
	minLivingArea   float64
	minKitchenArea  float64
	minRoomNumber   int
	maxRoomNumber   int
	minSpecificArea float64
	maxSpecificArea float64
	minFloor        int
	minTotalFloor   int
	maxTotalFloor   int
	minLongitude    float64
	maxLongitude    float64
	minLatitude     float64
	maxLatitude     float64
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
			originURL = fmt.Sprintf(fetcher.originURL, originURL)
		}
		imageURL := item.MainPhoto
		if imageURL != "" {
			imageURL = fmt.Sprintf(fetcher.imageURL, imageURL)
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
			float64(item.Longitude),
			float64(item.Latitude),
			state,
			item.CityNameUK,
			district,
			street,
			item.BuildingNumberStr,
		}
	}
	return flats, nil
}
