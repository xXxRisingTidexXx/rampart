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
			Timeout:   2 * time.Second,
			Portion:   10,
			Flags:     map[mining.Housing]string{mining.Primary: "pm_housing=1"},
			Headers:   map[string]string{"User-Agent": "domria-test-bot/v1.0.0"},
			SearchURL: "https://domria.ua/search/",
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
	if flats[0].originURL != "" {
		t.Errorf("domria: non-empty origin url, %v", flats[0])
	}
	if flats[0].imageURL != "" {
		t.Errorf("domria: non-empty image url, %v", flats[0])
	}
	if flats[0].updateTime != nil {
		t.Errorf("domria: non-empty update time, %v", flats[0])
	}
	if flats[0].price != 0 {
		t.Errorf("domria: non-zero price, %v", flats[0])
	}
	if flats[0].totalArea != 0 {
		t.Errorf("domria: non-zero total area, %v", flats[0])
	}
	if flats[0].livingArea != 0 {
		t.Errorf("domria: non-zero living area, %v", flats[0])
	}
	if flats[0].kitchenArea != 0 {
		t.Errorf("domria: non-zero kitchen area, %v", flats[0])
	}
	if flats[0].roomNumber != 0 {
		t.Errorf("domria: non-zero room number, %v", flats[0])
	}
	if flats[0].floor != 0 {
		t.Errorf("domria: non-zero floor, %v", flats[0])
	}
	if flats[0].totalFloor != 0 {
		t.Errorf("domria: non-zero total floor, %v", flats[0])
	}
	if flats[0].housing != mining.Primary {
		t.Errorf("domria: invalid housing, %v", flats[0])
	}
	if flats[0].complex != "" {
		t.Errorf("domria: non-empty complex, %v", flats[0])
	}
	if flats[0].point != nil {
		t.Errorf("domria: non-empty point, %v", flats[0])
	}
	if flats[0].state != "" {
		t.Errorf("domria: non-empty state, %v", flats[0])
	}
	if flats[0].city != "" {
		t.Errorf("domria: non-empty city, %v", flats[0])
	}
	if flats[0].district != "" {
		t.Errorf("domria: non-empty district, %v", flats[0])
	}
	if flats[0].street != "" {
		t.Errorf("domria: non-empty street, %v", flats[0])
	}
	if flats[0].houseNumber != "" {
		t.Errorf("domria: non-empty house number, %v", flats[0])
	}
}

func TestFetcherUnmarshalSearchValidItem(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("valid_item"), mining.Primary)
	if err != nil {
		t.Errorf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}

}
