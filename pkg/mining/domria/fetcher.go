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

func newFetcher(portion int, timeout time.Duration) *fetcher {
	if portion < 1 {
		panic(fmt.Sprintf("domria: fetcher got invalid portion, %d", portion))
	}
	return &fetcher{
		&http.Client{Timeout: timeout},
		"prospector/1.0 (rampart/prospector)",
		0,
		portion,
		map[mining.Housing]string{mining.Primary: "newbuildings=1", mining.Secondary: "secondary=1"},
		"https://dom.ria.com/searchEngine/?category=1&realty_type=2&opera" +
			"tion_type=1&fullCategoryOperation=1_2_1&%s&page=%d&limit=%d",
		"https://dom.ria.com/uk/%s",
		"https://cdn.riastatic.com/photos/%s",
	}
}

type fetcher struct {
	client       *http.Client
	userAgent    string
	page         int
	portion      int
	housingFlags map[mining.Housing]string
	searchURL    string
	originURL    string
	imageURL     string
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
		fetcher.page = 0
		return nil, nil
	}
	fetcher.page++
	flats := make([]*flat, 0, length)
	for _, item := range search.Items {
		if flat, err := fetcher.mapItem(item); err != nil {
			log.Error(err)
		} else {
			flats = append(flats, flat)
		}
	}
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
		return nil, fmt.Errorf("domria: failed to unmarshal the search, %v", search)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: failed to close the response body, %v", err)
	}
	return &search, nil
}

func (fetcher *fetcher) mapItem(item *item) (*flat, error) {
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
	if item.TotalSquareMeters <= 0 {
		return nil, fmt.Errorf("domria: non-positive total area at %s", originURL)
	}
	if item.LivingSquareMeters < 0 {
		return nil, fmt.Errorf("domria: negative living area at %s", originURL)
	}
	if item.KitchenSquareMeters < 0 {
		return nil, fmt.Errorf("domria: negative kitchen area at %s", originURL)
	}
	
	return &flat{
		originURL: originURL,
		imageURL: imageURL,
		updatedAt: (*time.Time)(item.UpdatedAt),
		price: price,
	}, nil
}
