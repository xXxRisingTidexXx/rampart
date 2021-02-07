package mining

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func NewDomriaMiner(config config.DomriaMiner) Miner {
	return &domriaMiner{
		config.Name,
		config.Spec,
		&http.Client{Timeout: config.Timeout},
		-1,
		config.RetryLimit,
		config.SearchPrefix,
		config.UserAgent,
		config.URLPrefix,
		config.ImageURLFormat,
		map[int]misc.Housing{1: misc.SecondaryHousing, 2: misc.PrimaryHousing},
		config.Swaps,
		config.CityOrthography,
		config.StreetOrthography,
		config.HouseNumberOrthography,
		config.HouseNumberMaxLength,
	}
}

type domriaMiner struct {
	name                   string
	spec                   string
	client                 *http.Client
	page                   int
	retryLimit             int
	searchPrefix           string
	userAgent              string
	urlPrefix              string
	imageURLFormat         string
	housings               map[int]misc.Housing
	swaps                  misc.Set
	cityOrthography        map[string]string
	streetOrthography      []string
	houseNumberOrthography []string
	houseNumberMaxLength   int
}

func (m *domriaMiner) Name() string {
	return m.name
}

func (m *domriaMiner) Spec() string {
	return m.spec
}

// TODO: retry metric.
func (m *domriaMiner) MineFlat() (Flat, error) {
	m.page++
	bytes, err := make([]byte, 0), io.EOF
	for retry := 0; retry < m.retryLimit && err != nil; retry++ {
		bytes, err = m.trySearch()
	}
	if err != nil {
		return Flat{}, err
	}
	var s search
	if err := json.Unmarshal(bytes, &s); err != nil {
		return Flat{}, fmt.Errorf("mining: miner failed to unmarshal the search, %v", err)
	}
	if len(s.Items) == 0 {
		m.page = -1
		return Flat{}, io.EOF
	}
	if s.Items[0].BeautifulURL == "" {
		return Flat{}, fmt.Errorf(
			"mining: miner ignored an item without an url on page %d",
			m.page,
		)
	}
	url := m.urlPrefix + s.Items[0].BeautifulURL
	index := strings.LastIndex(url, "-")
	if index == -1 {
		return Flat{}, fmt.Errorf("mining: miner ignored an item without a dash in url, %s", url)
	}
	urls, slug := make([]string, 0, len(s.Items[0].Photos)), url[:index]
	for id := range s.Items[0].Photos {
		urls = append(urls, fmt.Sprintf(m.imageURLFormat, slug, id))
	}
	if s.Items[0].SaleDate != "" {
		return Flat{}, fmt.Errorf(
			"mining: miner ignored a sold on %s item, %s",
			s.Items[0].SaleDate,
			url,
		)
	}
	housing, ok := m.housings[s.Items[0].RealtySaleType]
	if !ok {
		return Flat{}, fmt.Errorf(
			"mining: miner ignored an item with housing %d, %s",
			s.Items[0].RealtySaleType,
			url,
		)
	}
	if s.Items[0].Longitude < -180 || s.Items[0].Longitude > 180 {
		return Flat{}, fmt.Errorf(
			"mining: miner ignored an outlier longitude %f item, %s",
			s.Items[0].Longitude,
			url,
		)
	}
	if s.Items[0].Latitude < -90 || s.Items[0].Latitude > 90 {
		return Flat{}, fmt.Errorf(
			"mining: miner ignored an outlier latitude %f item, %s",
			s.Items[0].Latitude,
			url,
		)
	}
	street := s.Items[0].StreetNameUK
	if street == "" {
		street = s.Items[0].StreetName
	}
	return Flat{
		URL:         url,
		ImageURLs:   urls,
		Price:       float64(s.Items[0].PriceArr.USD),
		TotalArea:   s.Items[0].TotalSquareMeters,
		LivingArea:  s.Items[0].LivingSquareMeters,
		KitchenArea: s.Items[0].KitchenSquareMeters,
		RoomNumber:  s.Items[0].RoomsCount,
		Floor:       s.Items[0].Floor,
		TotalFloor:  s.Items[0].FloorsCount,
		Housing:     housing,
		Point:       orb.Point{float64(s.Items[0].Longitude), float64(s.Items[0].Latitude)},
		City:        s.Items[0].CityNameUK,
		Street:      street,
		HouseNumber: string(s.Items[0].BuildingNumberStr),
	}, nil
}

func (m *domriaMiner) trySearch() ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, m.searchPrefix+strconv.Itoa(m.page), nil)
	if err != nil {
		return nil, fmt.Errorf("mining: miner failed to construct a request, %v", err)
	}
	request.Header.Set("User-Agent", m.userAgent)
	response, err := m.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("mining: miner failed to make a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("mining: miner got response with status %s", response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("mining: miner failed to read the response body, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("mining: miner failed to close the response body, %v", err)
	}
	return bytes, nil
}
