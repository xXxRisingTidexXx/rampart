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
		&configs.Fetcher{
			Timeout:         2 * time.Second,
			Portion:         10,
			Flags:           map[mining.Housing]string{mining.Primary: "pm_housing=1"},
			Headers:         map[string]string{"User-Agent": "domria-test-bot/v1.0.0"},
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

func TestFetcherUnmarshalSearchEmptyItem(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("empty_item"), mining.Primary)
	if err != nil {
		t.Errorf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	if flats[0].originURL != "" || flats[0].imageURL != "" {
		t.Errorf("domria: unexpected flat, %v", flats[0])
	}
}

//func TestFetcherUnmarshalSearchValidItem(t *testing.T) {}
