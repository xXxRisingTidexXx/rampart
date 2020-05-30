package domria

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"rampart/pkg"
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
		"https://dom.ria.com/node/api/autocompleteCities?langId=4",
		"https://dom.ria.com/searchEngine/?category=1&realty_type=2&operation_type=1" +
			"&fullCategoryOperation=1_2_1&%s&page=%d&state_id=%d&city_id=%d&limit=%d",
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
	localityURL  string
	searchURL    string
	originURL    string
	imageURL     string
}

func (fetcher *fetcher) fetchFlats(state, city string, housing mining.Housing) ([]*flat, error) {
	var localities []*locality
	if err := fetcher.fetchJSON(fetcher.localityURL, &localities); err != nil {
		return nil, err
	}
	stateID, cityID, err := fetcher.findLocalityIDs(localities, state, city)
	if err != nil {
		return nil, err
	}
	housingFlag, ok := fetcher.housingFlags[housing]
	if !ok {
		return nil, fmt.Errorf("domria: couldn't find flag with housing %v", housing)
	}
	var search search
	err = fetcher.fetchJSON(
		fmt.Sprintf(fetcher.searchURL, housingFlag, fetcher.page, stateID, cityID, fetcher.portion),
		&search,
	)
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

func (fetcher *fetcher) fetchJSON(url string, target pkg.Any) error {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", fetcher.userAgent)
	response, err := fetcher.client.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, target); err != nil {
		return err
	}
	if err := response.Body.Close(); err != nil {
		return err
	}
	return nil
}

func (fetcher *fetcher) findLocalityIDs(localities []*locality, city, state string) (int, int, error) {
	for _, locality := range localities {
		if locality.State == state && locality.City == city {
			return locality.Payload.StateID, locality.Payload.CityID, nil
		}
	}
	return 0, 0, fmt.Errorf("domria: unknown locality with state %s and city %s", state, city)
}

func (fetcher *fetcher) mapItem(item *item) (*flat, error) {
	if item.BeautifulURL == "" {
		return nil, fmt.Errorf("domria: item url can't be empty")
	}
	originURL := fmt.Sprintf(fetcher.originURL, item.BeautifulURL)
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
	return &flat{

	}, nil
}
