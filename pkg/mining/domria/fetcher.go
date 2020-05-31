package domria

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"rampart/pkg/mining"
	"strconv"
	"strings"
	"time"
)

func newFetcher() *fetcher {
	return &fetcher{
		&http.Client{Timeout: 20 * time.Second},
		"prospector/1.0 (rampart/prospector)",
		0,
		100,
		map[mining.Housing]string{mining.Primary: "newbuildings=1", mining.Secondary: "secondary=1"},
		"https://dom.ria.com/searchEngine/?category=1&realty_type=2&opera" +
			"tion_type=1&fullCategoryOperation=1_2_1&%s&page=%d&limit=%d",
		"https://dom.ria.com/uk/%s",
		"https://cdn.riastatic.com/photos/%s",
		10,
		560,
		8,
		2,
		1,
		9,
		12,
		130,
		1,
		2,
		50,
	}
}

type fetcher struct {
	client          *http.Client
	userAgent       string
	page            int
	portion         int
	housingFlags    map[mining.Housing]string
	searchURL       string
	originURL       string
	imageURL        string
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
}

func (fetcher *fetcher) fetchFlats(housing mining.Housing) ([]*flat, error) {
	housingFlag, ok := fetcher.housingFlags[housing]
	if !ok {
		return nil, fmt.Errorf("domria: %v housing isn't acceptable", housing)
	}
	search, err := fetcher.fetchSearch(housingFlag)
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

func (fetcher *fetcher) fetchSearch(housingFlag string) (*search, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(fetcher.searchURL, housingFlag, fetcher.page, fetcher.portion),
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
	rawPrice, ok := item.PriceArr["1"]
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
	if item.LivingSquareMeters != 0 &&
		(item.LivingSquareMeters < fetcher.minLivingArea ||
			item.LivingSquareMeters >= item.TotalSquareMeters) {
		return nil, fmt.Errorf("domria: living area is out of range at %s", originURL)
	}
	if item.KitchenSquareMeters != 0 &&
		(item.KitchenSquareMeters < fetcher.minKitchenArea ||
			item.KitchenSquareMeters >= item.TotalSquareMeters) {
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
	if item.Longitude < -180 || item.Longitude > 180 {
		return nil, fmt.Errorf("domria: longitude is out of range at %s", originURL)
	}
	if item.Latitude < -90 || item.Latitude > 90 {
		return nil, fmt.Errorf("domria: latitude is out of range at %s", originURL)
	}
	district := item.DistrictNameUK
	if district != "" && item.DistrictTypeName == "Район" && strings.HasSuffix(district, "ий") {
		district += " район"
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
		float64(item.Longitude),
		float64(item.Latitude),
		item.StateNameUK + " область",
		item.CityNameUK,
		district,
		street,
		item.BuildingNumberStr,
	}, nil
}
