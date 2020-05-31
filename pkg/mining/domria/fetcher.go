package domria

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"rampart/pkg/mining"
	"rampart/pkg/mining/configs"
	"strconv"
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
	search, err := fetcher.fetchSearch(flag)
	if err != nil {
		return nil, err
	}
	length := len(search.Items)
	if length == 0 {
		log.Debugf("domria: %s housing fetcher on %d page reset", housing, fetcher.page)
		fetcher.page = 0
		return nil, nil
	}
	flats := make([]*flat, 0, length)
	for _, item := range search.Items {
		if flat, err := fetcher.mapItem(item, housing); err != nil {
			log.Error(err)
		} else {
			flats = append(flats, flat)
		}
	}
	log.Debugf("domria: %s housing fetcher on %d page fetched %d flats", housing, fetcher.page, len(flats))
	fetcher.page++
	return flats, nil
}

func (fetcher *fetcher) fetchSearch(flag string) (*search, error) {
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
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("domria: failed to read the response body, %v", err)
	}
	var search search
	if err := json.Unmarshal(body, &search); err != nil {
		return nil, fmt.Errorf("domria: failed to unmarshal the search, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: failed to close the response body, %v", err)
	}
	return &search, nil
}

func (fetcher *fetcher) mapItem(item *item, housing mining.Housing) (*flat, error) {
	if item.BeautifulURL == "" {
		return nil, fmt.Errorf("domria: item url can't be empty")
	}
	originURL := fmt.Sprintf(fetcher.originURL, item.BeautifulURL)
	imageURL := item.MainPhoto
	if imageURL != "" {
		imageURL = fmt.Sprintf(fetcher.imageURL, imageURL)
	}
	rawPrice, ok := item.PriceArr[fetcher.usdLabel]
	if !ok {
		return nil, fmt.Errorf("domria: absent USD price at %s", originURL)
	}
	price, err := strconv.ParseFloat(strings.ReplaceAll(rawPrice, " ", ""), 64)
	if err != nil {
		return nil, fmt.Errorf("domria: invalid USD price at %s, %v", originURL, err)
	}
	if item.TotalSquareMeters < fetcher.minTotalArea || item.TotalSquareMeters > fetcher.maxTotalArea {
		return nil, fmt.Errorf("domria: total area is out of range at %s", originURL)
	}
	if item.LivingSquareMeters < fetcher.minLivingArea ||
		item.LivingSquareMeters >= item.TotalSquareMeters {
		return nil, fmt.Errorf("domria: living area is out of range at %s", originURL)
	}
	if item.KitchenSquareMeters < fetcher.minKitchenArea ||
		item.KitchenSquareMeters >= item.TotalSquareMeters {
		return nil, fmt.Errorf("domria: kitchen area is out of range at %s", originURL)
	}
	if item.RoomsCount < fetcher.minRoomNumber || item.RoomsCount > fetcher.maxRoomNumber {
		return nil, fmt.Errorf("domria: room number is out of range at %s", originURL)
	}
	specificArea := item.TotalSquareMeters / float64(item.RoomsCount)
	if specificArea < fetcher.minSpecificArea || specificArea > fetcher.maxSpecificArea {
		return nil, fmt.Errorf("domria: specific area is out of range at %s", originURL)
	}
	if item.FloorsCount < fetcher.minTotalFloor || item.FloorsCount > fetcher.maxTotalFloor {
		return nil, fmt.Errorf("domria: total floor is out of range, %s", originURL)
	}
	if item.Floor < fetcher.minFloor || item.Floor > item.FloorsCount {
		return nil, fmt.Errorf("domria: floor is out of range, %s", originURL)
	}
	longitude := float64(item.Longitude)
	if longitude < fetcher.minLongitude || longitude > fetcher.maxLongitude {
		return nil, fmt.Errorf("domria: longitude is out of range at %s", originURL)
	}
	latitude := float64(item.Latitude)
	if latitude < fetcher.minLatitude || latitude > fetcher.maxLatitude {
		return nil, fmt.Errorf("domria: latitude is out of range at %s", originURL)
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
	return &flat{
		originURL,
		imageURL,
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
		latitude,
		longitude,
		state,
		item.CityNameUK,
		district,
		street,
		item.BuildingNumberStr,
	}, nil
}
