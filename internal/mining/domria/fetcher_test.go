package domria

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"rampart/internal/mining"
	"rampart/internal/mining/configs"
	"testing"
	"time"
)

func TestFetcherUnmarshalSearchEmptyString(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte(""), mining.Primary)
	if flats != nil {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
	if err == nil || err.Error() != "domria: failed to unmarshal the search, unexpected end of JSON input" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
}

func newDefaultFetcher() *fetcher {
	return newFetcher(
		"domria-test-bot/v1.0.0",
		&configs.Fetcher{
			Timeout:         2 * time.Second,
			Portion:         10,
			Flags:           map[mining.Housing]string{mining.Primary: "pm_housing=1"},
			SearchURL:       "https://domria.ua/search/",
			OriginURLPrefix: "https://domria.ua/uk/",
			ImageURLPrefix:  "https://cdn.riastatic.ua/photos/",
			StateEnding:     "ька",
			StateSuffix:     " область",
			DistrictLabel:   "Район",
			DistrictEnding:  "ий",
			DistrictSuffix:  " район",
		},
	)
}

func TestFetcherUnmarshalSearchInvalidJSON(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte("{"), mining.Primary)
	if flats != nil {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
	if err == nil || err.Error() != "domria: failed to unmarshal the search, unexpected end of JSON input" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
}

func TestFetcherUnmarshalSearchArrayInsteadOfObject(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte("[]"), mining.Primary)
	if flats != nil {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
	if err == nil || err.Error() != "domria: failed to unmarshal the search"+
		", json: cannot unmarshal array into Go value of type domria.search" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
}

func TestFetcherUnmarshalSearchEmptyJSON(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte("{}"), mining.Primary)
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
	if err != nil {
		t.Errorf("domria: unexpected error, %v", err)
	}
}

func TestFetcherUnmarshalSearchWithoutItems(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("without_items"), mining.Primary)
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
	if err != nil {
		t.Errorf("domria: unexpected error, %v", err)
	}
}

func readAll(fixtureName string) []byte {
	file, err := os.Open(filepath.Join("testdata", fixtureName+".json"))
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	if err = file.Close(); err != nil {
		panic(err)
	}
	return bytes
}

func TestFetcherUnmarshalSearchEmptySearch(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("empty_search"), mining.Primary)
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
	if err != nil {
		t.Errorf("domria: unexpected error, %v", err)
	}
}

//func TestFetcherUnmarshalSearchInvalidItem(t *testing.T) {}

//func TestFetcherUnmarshalSearchEmptyItem(t *testing.T) {}
//
//func TestFetcherUnmarshalSearchValidItem(t *testing.T) {}
