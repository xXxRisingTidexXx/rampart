package domria

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom"
	"io/ioutil"
	"net/http"
	"rampart/internal/mining/config"
	"rampart/internal/misc"
	"strings"
	"time"
)

func newGeocoder(config *config.Geocoder) *geocoder {
	return &geocoder{
		&http.Client{Timeout: time.Duration(config.Timeout)},
		config.Headers,
		config.StatelessCities,
		config.SearchURL,
		config.MinLookup,
	}
}

type geocoder struct {
	client          *http.Client
	headers         map[string]string
	statelessCities *misc.Set
	searchURL       string
	minLookup       float64
}

func (geocoder *geocoder) geocodeFlats(flats []*flat) []*flat {
	expectedLength := len(flats)
	if expectedLength == 0 {
		log.Debug("domria: geocoder skipped flats")
		return flats
	}
	geocodedNumber, locatedNumber, duration := 0.0, 0.0, 0.0
	newFlats := make([]*flat, 0, expectedLength)
	for i := range flats {
		if flats[i].point != nil {
			newFlats = append(newFlats, flats[i])
		} else if flats[i].district != "" && flats[i].street != "" && flats[i].houseNumber != "" {
			geocodedNumber++
			start := time.Now()
			bytes, err := geocoder.getLocations(flats[i])
			duration += time.Since(start).Seconds()
			var newFlat *flat
			if err == nil {
				newFlat, err = geocoder.locateFlat(flats[i], bytes)
			}
			if err != nil {
				log.Error(err)
			} else if newFlat != nil {
				locatedNumber++
				newFlats = append(newFlats, newFlat)
			}
		}
	}
	lookup := 0.0
	if geocodedNumber != 0 {
		duration /= geocodedNumber
		lookup = locatedNumber / geocodedNumber
	}
	log.Debugf("domria: geocoder passed %d flats (%.3fs)", len(newFlats), duration)
	if lookup < geocoder.minLookup {
		log.Warningf("domria: geocoder met low lookup (%.3f)", lookup)
	}
	return newFlats
}

func (geocoder *geocoder) getLocations(flat *flat) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, geocoder.getSearchURL(flat), nil)
	if err != nil {
		return nil, fmt.Errorf("domria: geocoder failed to construct a request, %v", err)
	}
	for key, value := range geocoder.headers {
		request.Header.Set(key, value)
	}
	response, err := geocoder.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("domria: geocoder failed to perform a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: geocoder got response with status %s", response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: geocoder failed to read the response body, %v", err)
	}
	if err = response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: geocoder failed to close the response body, %v", err)
	}
	return bytes, nil
}

func (geocoder *geocoder) getSearchURL(flat *flat) string {
	state := ""
	whitespace, plus := " ", "+"
	if !geocoder.statelessCities.Contains(flat.city) {
		state = strings.ReplaceAll(flat.state, whitespace, plus)
	}
	return fmt.Sprintf(
		geocoder.searchURL,
		state,
		strings.ReplaceAll(flat.city, whitespace, plus),
		strings.ReplaceAll(flat.district, whitespace, plus),
		strings.ReplaceAll(flat.street, whitespace, plus),
		strings.ReplaceAll(flat.houseNumber, whitespace, plus),
	)
}

func (geocoder *geocoder) locateFlat(f *flat, bytes []byte) (*flat, error) {
	var locations []*location
	if err := json.Unmarshal(bytes, &locations); err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to unmarshal the locations, %v", err)
	}
	if len(locations) == 0 {
		return nil, nil
	}
	return &flat{
		f.originURL,
		f.imageURL,
		f.updateTime,
		f.price,
		f.totalArea,
		f.livingArea,
		f.kitchenArea,
		f.roomNumber,
		f.floor,
		f.totalFloor,
		f.housing,
		f.complex,
		geom.NewPointFlat(
			geom.XY,
			[]float64{float64(locations[0].Lon), float64(locations[0].Lat)},
		).SetSRID(4326),
		f.state,
		f.city,
		f.district,
		f.street,
		f.houseNumber,
	}, nil
}
