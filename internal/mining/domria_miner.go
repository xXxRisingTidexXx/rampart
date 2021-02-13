package mining

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewDomriaMiner(config config.DomriaMiner) Miner {
	return &domriaMiner{
		config.Name,
		config.Spec,
		&http.Client{Timeout: config.Timeout},
		config.Page - 1,
		config.RetryLimit,
		config.SearchPrefix,
		config.UserAgent,
		config.URLPrefix,
		config.ImageURLFormat,
		config.MaxTotalArea,
		config.MaxRoomNumber,
		config.MaxTotalFloor,
		map[int]misc.Housing{1: misc.SecondaryHousing, 2: misc.PrimaryHousing},
		config.Swaps,
		config.Cities,
		strings.NewReplacer(config.StreetReplacements...),
		strings.NewReplacer(config.HouseNumberReplacements...),
		config.MaxHouseNumberLength,
	}
}

type domriaMiner struct {
	name                 string
	spec                 string
	client               *http.Client
	page                 int
	retryLimit           int
	searchPrefix         string
	userAgent            string
	urlPrefix            string
	imageURLFormat       string
	maxTotalArea         float64
	maxRoomNumber        int
	maxTotalFloor        int
	housings             map[int]misc.Housing
	swaps                misc.Set
	cities               map[string]string
	streetReplacer       *strings.Replacer
	houseNumberReplacer  *strings.Replacer
	maxHouseNumberLength int
}

func (m *domriaMiner) Name() string {
	return m.name
}

func (m *domriaMiner) Spec() string {
	return m.spec
}

func (m *domriaMiner) MineFlat() (Flat, error) {
	m.page++
	s, err := search{}, io.EOF
	for retry := 0; retry < m.retryLimit && err != nil; retry++ {
		now := time.Now()
		s, err = m.trySearch()
		metrics.MessisMiningDuration.WithLabelValues(m.name).Observe(time.Since(now).Seconds())
		metrics.MessisMiningRetries.WithLabelValues(m.name).Inc()
	}
	if err != nil {
		metrics.MessisMinings.WithLabelValues(m.name, "failure").Inc()
		return Flat{}, err
	}
	if len(s.Items) == 0 {
		m.page = -1
		metrics.MessisMinings.WithLabelValues(m.name, "reset").Inc()
		return Flat{}, io.EOF
	}
	if err := m.validateItem(s.Items[0]); err != nil {
		metrics.MessisMinings.WithLabelValues(m.name, "validation").Inc()
		return Flat{}, err
	}
	url := m.urlPrefix + s.Items[0].BeautifulURL
	urls, slug := make([]string, 0, len(s.Items[0].Photos)), url[:strings.LastIndex(url, "-")]
	for id := range s.Items[0].Photos {
		urls = append(urls, fmt.Sprintf(m.imageURLFormat, slug, id))
	}
	city := strings.TrimSpace(s.Items[0].CityNameUK)
	if m.swaps.Contains(city) {
		city = strings.TrimSpace(s.Items[0].DistrictNameUK)
		metrics.MessisMiningSanitations.WithLabelValues(m.name, "swap").Inc()
	}
	if value, ok := m.cities[city]; ok {
		city = value
		metrics.MessisMiningSanitations.WithLabelValues(m.name, "city").Inc()
	}
	initialStreet := s.Items[0].StreetNameUK
	if initialStreet == "" {
		initialStreet = s.Items[0].StreetName
	}
	street := initialStreet
	houseNumber := m.sanitizeHouseNumber(string(s.Items[0].BuildingNumberStr))
	if index := strings.Index(initialStreet, ","); index != -1 {
		street = initialStreet[:index]
		metrics.MessisMiningSanitations.WithLabelValues(m.name, "street").Inc()
		extraNumber := m.sanitizeHouseNumber(initialStreet[index+1:])
		if houseNumber == "" &&
			extraNumber != "" &&
			extraNumber[0] >= '0' &&
			extraNumber[0] <= '9' {
			houseNumber = extraNumber
			metrics.MessisMiningSanitations.WithLabelValues(m.name, "house-number").Inc()
		}
	}
	if runes := []rune(houseNumber); len(runes) > m.maxHouseNumberLength {
		houseNumber = string(runes[:m.maxHouseNumberLength])
	}
	metrics.MessisMinings.WithLabelValues(m.name, "success").Inc()
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
		Housing:     m.housings[s.Items[0].RealtySaleType],
		Point:       orb.Point{float64(s.Items[0].Longitude), float64(s.Items[0].Latitude)},
		City:        city,
		Street:      strings.TrimSpace(m.streetReplacer.Replace(street)),
		HouseNumber: houseNumber,
	}, nil
}

func (m *domriaMiner) trySearch() (search, error) {
	var s search
	request, err := http.NewRequest(http.MethodGet, m.searchPrefix+strconv.Itoa(m.page), nil)
	if err != nil {
		return s, fmt.Errorf("mining: miner failed to construct a request, %v", err)
	}
	request.Header.Set("User-Agent", m.userAgent)
	response, err := m.client.Do(request)
	if err != nil {
		return s, fmt.Errorf("mining: miner failed to make a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return s, fmt.Errorf("mining: miner got response with status %s", response.Status)
	}
	if err := json.NewDecoder(response.Body).Decode(&s); err != nil {
		_ = response.Body.Close()
		return s, fmt.Errorf("mining: miner failed to unmarshal a search, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return s, fmt.Errorf("mining: miner failed to close a response body, %v", err)
	}
	return s, nil
}

func (m *domriaMiner) validateItem(i item) error {
	if i.BeautifulURL == "" {
		return fmt.Errorf("mining: miner ignored an item without an url on page %d", m.page)
	}
	url := m.urlPrefix + i.BeautifulURL
	if strings.LastIndex(url, "-") == -1 {
		return fmt.Errorf("mining: miner ignored an item without a dash in url, %s", url)
	}
	if len(i.Photos) == 0 {
		return fmt.Errorf("mining: miner ignored an item without images, %s", url)
	}
	if i.SaleDate != "" {
		return fmt.Errorf("mining: miner ignored an item sold on %s, %s", i.SaleDate, url)
	}
	if i.PriceArr.USD <= 0 {
		return fmt.Errorf("mining: miner ignored an item with price %f, %s", i.PriceArr.USD, url)
	}
	if i.TotalSquareMeters <= 0 || i.TotalSquareMeters > m.maxTotalArea {
		return fmt.Errorf(
			"mining: miner ignored an item with total area %f, %s",
			i.TotalSquareMeters,
			url,
		)
	}
	if i.LivingSquareMeters < 0 || i.LivingSquareMeters > i.TotalSquareMeters {
		return fmt.Errorf(
			"mining: miner ignored an item with living area %f, %s",
			i.LivingSquareMeters,
			url,
		)
	}
	if i.KitchenSquareMeters < 0 || i.KitchenSquareMeters > i.TotalSquareMeters {
		return fmt.Errorf(
			"mining: miner ignored an item with kitchen area %f, %s",
			i.KitchenSquareMeters,
			url,
		)
	}
	if i.RoomsCount <= 0 || i.RoomsCount > m.maxRoomNumber {
		return fmt.Errorf(
			"mining: miner ignored an item with room number %d, %s",
			i.RoomsCount,
			url,
		)
	}
	if i.FloorsCount <= 0 || i.FloorsCount > m.maxTotalFloor {
		return fmt.Errorf(
			"mining: miner ignored an item with total floor %d, %s",
			i.FloorsCount,
			url,
		)
	}
	if i.Floor <= 0 || i.Floor > i.FloorsCount {
		return fmt.Errorf("mining: miner ignored an item with floor %d, %s", i.Floor, url)
	}
	if _, ok := m.housings[i.RealtySaleType]; !ok {
		return fmt.Errorf(
			"mining: miner ignored an item with housing %d, %s",
			i.RealtySaleType,
			url,
		)
	}
	if i.Longitude < -180 || i.Longitude > 180 {
		return fmt.Errorf("mining: miner ignored an item with longitude %f, %s", i.Longitude, url)
	}
	if i.Latitude < -90 || i.Latitude > 90 {
		return fmt.Errorf("mining: miner ignored an item with latitude %f, %s", i.Latitude, url)
	}
	return nil
}

func (m *domriaMiner) sanitizeHouseNumber(houseNumber string) string {
	if houseNumber == "" {
		return houseNumber
	}
	newHouseNumber := m.houseNumberReplacer.Replace(houseNumber)
	if index := strings.Index(newHouseNumber, ","); index != -1 {
		return newHouseNumber[:index]
	}
	return newHouseNumber
}
