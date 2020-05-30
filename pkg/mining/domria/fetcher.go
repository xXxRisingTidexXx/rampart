package domria

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rampart/pkg"
	"rampart/pkg/mining"
	"time"
)

func newFetcher(portion int, timeout time.Duration) *fetcher {
	return &fetcher{
		&http.Client{Timeout: timeout},
		0,
		portion,
		map[mining.Housing]string{mining.Primary: "newbuildings=1", mining.Secondary: "secondary=1"},
		"prospector/1.0 (rampart/prospector)",
		"https://dom.ria.com/node/api/autocompleteCities?langId=4",
		"https://dom.ria.com/searchEngine/?category=1&realty_type=2&operation_type=1" +
			"&fullCategoryOperation=1_2_1&%s&page=%d&state_id=%d&city_id=%d&limit=%d",
	}
}

type fetcher struct {
	client       *http.Client
	page         int
	portion      int
	housingFlags map[mining.Housing]string
	userAgent    string
	localityURL  string
	searchURL    string
}

func (fetcher *fetcher) fetchFlats(state, city string, housing mining.Housing) ([]*flat, error) {
	var localities []*locality
	if err := fetcher.fetchJSON(fetcher.localityURL, &localities); err != nil {
		return nil, err
	}
	stateID, cityID := -1, -1
	for _, locality := range localities {
		if locality.State == state && locality.City == city {
			stateID, cityID = locality.Payload.StateID, locality.Payload.CityID
			break
		}
	}
	if stateID == -1 || cityID == -1 {
		return nil, fmt.Errorf("domria: couldn't find locality with state %s and city %s", state, city)
	}
	housingFlag, ok := fetcher.housingFlags[housing]
	if !ok {
		return nil, fmt.Errorf("domria: couldn't find flag with housing %v", housing)
	}
	var search search
	err := fetcher.fetchJSON(
		fmt.Sprintf(fetcher.searchURL, housingFlag, 1, stateID, cityID, fetcher.portion),
		&search,
	)
	if err != nil {
		return nil, err
	}
	flats := make([]*flat, 0, len(search.Items))
	
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
