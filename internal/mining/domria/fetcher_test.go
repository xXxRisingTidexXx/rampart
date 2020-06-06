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
	assertFlat(t, flats[0], &flat{})
}

func assertFlat(t *testing.T, actual *flat, expected *flat) {
	if actual == nil {
		t.Fatal("domria: empty actual")
	}
	if expected == nil {
		t.Fatal("domria: empty expected")
	}
	if actual.originURL != expected.originURL {
		t.Errorf("domria: invalid origin url, %s", actual.originURL)
	}
	if actual.imageURL != expected.imageURL {
		t.Errorf("domria: invalid image url, %s", actual.imageURL)
	}
	if expected.updateTime == nil && actual.updateTime != nil {
		t.Errorf("domria: non-nil update time, %v", actual.updateTime)
	}
	if expected.updateTime != nil &&
		(actual.updateTime == nil || !actual.updateTime.Equal(*expected.updateTime)) {
		t.Errorf("domria: invalid update time, %v", actual.updateTime)
	}
	if actual.price != expected.price {
		t.Errorf("domria: invalid price, %.1f", actual.price)
	}
	if actual.totalArea != expected.totalArea {
		t.Errorf("domria: invalid total area, %1.f", actual.totalArea)
	}
	if actual.livingArea != expected.livingArea {
		t.Errorf("domria: invalid living area, %.1f", actual.livingArea)
	}
	if actual.kitchenArea != expected.kitchenArea {
		t.Errorf("domria: invalid kitchen area, %.1f", actual.kitchenArea)
	}
	if actual.roomNumber != expected.roomNumber {
		t.Errorf("domria: invalid room number, %d", actual.roomNumber)
	}
	if actual.floor != expected.floor {
		t.Errorf("domria: invalid floor, %d", actual.floor)
	}
	if actual.totalFloor != expected.totalFloor {
		t.Errorf("domria: invalid total floor, %d", actual.totalFloor)
	}
	if actual.housing != mining.Primary {
		t.Errorf("domria: invalid housing, %s", actual.housing)
	}
	if actual.complex != expected.complex {
		t.Errorf("domria: invalid complex, %s", actual.complex)
	}
	if expected.point == nil && actual.point != nil {
		t.Errorf("domria: non-nil point, %v", actual.point)
	}
	if expected.point != nil &&
		(actual.point == nil ||
			actual.point.X() != expected.point.X() ||
			actual.point.Y() != expected.point.Y()) {
		t.Errorf("domria: invalid point, %v", actual.point)
	}
	if actual.state != expected.state {
		t.Errorf("domria: invalid state, %s", actual.state)
	}
	if actual.city != expected.city {
		t.Errorf("domria: invalid city, %s", actual.city)
	}
	if actual.district != expected.district {
		t.Errorf("domria: invalid district, %s", actual.district)
	}
	if actual.street != expected.street {
		t.Errorf("domria: invalid street, %s", actual.street)
	}
	if actual.houseNumber != expected.houseNumber {
		t.Errorf("domria: invalid house number, %s", actual.houseNumber)
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
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-rovno-schastlivoe-chernovola-vyacheslava-ulitsa-16818824.html"
		},
	)
	if flats[0].originURL != "realty-prodaja-kvartira-rovno-schastlivoe-chernovola-vyacheslava-ulitsa-16818824.html" {
		t.Errorf("domria: invalid origin url, %v", flats[0])
	}
	if flats[0].imageURL != "dom/photo/10925/1092575/109257503/109257503.jpg" {
		t.Errorf("domria: invalid image url, %v", flats[0])
	}
	if flats[0].updateTime == nil ||
		flats[0].updateTime.Unix() != time.Date(2020, time.June, 6, 14, 57, 18, 0, time.Local).Unix() {
		t.Errorf("domria: invalid update time, %v", flats[0])
	}
	if flats[0].price != 27800 {
		t.Errorf("domria: invalid price, %v", flats[0])
	}
	if flats[0].totalArea != 45 {
		t.Errorf("domria: invalid total area, %v", flats[0])
	}
	if flats[0].livingArea != 0 {
		t.Errorf("domria: invalid living area, %v", flats[0])
	}
	if flats[0].kitchenArea != 0 {
		t.Errorf("domria: invalid kitchen area, %v", flats[0])
	}
	if flats[0].roomNumber != 1 {
		t.Errorf("domria: invalid room number, %v", flats[0])
	}
	if flats[0].floor != 2 {
		t.Errorf("domria: invalid floor, %v", flats[0])
	}
	if flats[0].totalFloor != 9 {
		t.Errorf("domria: invalid total floor, %v", flats[0])
	}
	if flats[0].housing != mining.Primary {
		t.Errorf("domria: invalid housing, %v", flats[0])
	}
	if flats[0].complex != "ЖК На Щасливому, будинок 27" {
		t.Errorf("domria: invalid complex, %v", flats[0])
	}
	if flats[0].point == nil || flats[0].point.X() != 26.267247115344 || flats[0].point.Y() != 50.59766586795 {
		t.Errorf("domria: invalid point, %v", flats[0])
	}
	if flats[0].state != "Рівненська" {
		t.Errorf("domria: invalid state, %v", flats[0])
	}
	if flats[0].city != "Рівне" {
		t.Errorf("domria: invalid city, %v", flats[0])
	}
	if flats[0].district != "Щасливе" {
		t.Errorf("domria: invalid district, %v", flats[0])
	}
	if flats[0].street != "Черновола Вячеслава улица" {
		t.Errorf("domria: invalid street, %v", flats[0])
	}
	if flats[0].houseNumber != "91-Ф" {
		t.Errorf("domria: invalid house number, %v", flats[0])
	}
}
