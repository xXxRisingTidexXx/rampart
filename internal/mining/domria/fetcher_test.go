package domria

import (
	"github.com/twpayne/go-geom"
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
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
}

func TestFetcherUnmarshalSearchWithoutItems(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("without_items"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
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
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
}

func TestFetcherUnmarshalSearchEmptyItem(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("empty_item"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(t, flats[0], &flat{housing: mining.Primary})
}

//nolint:gocognit,gocyclo
func assertFlat(t *testing.T, actual *flat, expected *flat) {
	if actual == nil {
		t.Fatal("domria: empty actual")
	}
	if expected == nil {
		t.Fatal("domria: empty expected")
	}
	if actual.originURL != expected.originURL {
		t.Errorf("domria: invalid origin url, %s != %s", actual.originURL, expected.originURL)
	}
	if actual.imageURL != expected.imageURL {
		t.Errorf("domria: invalid image url, %s != %s", actual.imageURL, expected.imageURL)
	}
	if expected.updateTime == nil && actual.updateTime != nil {
		t.Errorf("domria: non-nil update time, %v", actual.updateTime)
	}
	if expected.updateTime != nil &&
		(actual.updateTime == nil || !actual.updateTime.Equal(*expected.updateTime)) {
		t.Errorf("domria: invalid update time, %v != %v", actual.updateTime, expected.updateTime)
	}
	if actual.price != expected.price {
		t.Errorf("domria: invalid price, %.1f != %.1f", actual.price, expected.price)
	}
	if actual.totalArea != expected.totalArea {
		t.Errorf("domria: invalid total area, %1.f != %.1f", actual.totalArea, expected.totalArea)
	}
	if actual.livingArea != expected.livingArea {
		t.Errorf("domria: invalid living area, %.1f != %.1f", actual.livingArea, expected.livingArea)
	}
	if actual.kitchenArea != expected.kitchenArea {
		t.Errorf("domria: invalid kitchen area, %.1f != %.1f", actual.kitchenArea, expected.kitchenArea)
	}
	if actual.roomNumber != expected.roomNumber {
		t.Errorf("domria: invalid room number, %d != %d", actual.roomNumber, expected.roomNumber)
	}
	if actual.floor != expected.floor {
		t.Errorf("domria: invalid floor, %d != %d", actual.floor, expected.floor)
	}
	if actual.totalFloor != expected.totalFloor {
		t.Errorf("domria: invalid total floor, %d != %d", actual.totalFloor, expected.totalFloor)
	}
	if actual.housing != expected.housing {
		t.Errorf("domria: invalid housing, %s != %s", actual.housing, expected.housing)
	}
	if actual.complex != expected.complex {
		t.Errorf("domria: invalid complex, %s != %s", actual.complex, expected.complex)
	}
	if expected.point == nil && actual.point != nil {
		t.Errorf("domria: non-nil point, %v", actual.point)
	}
	if expected.point != nil &&
		(actual.point == nil ||
			actual.point.Layout() != expected.point.Layout() ||
			actual.point.X() != expected.point.X() ||
			actual.point.Y() != expected.point.Y()) {
		t.Errorf("domria: invalid point, %v != %v", actual.point, expected.point)
	}
	if actual.state != expected.state {
		t.Errorf("domria: invalid state, %s != %s", actual.state, expected.state)
	}
	if actual.city != expected.city {
		t.Errorf("domria: invalid city, %s != %s", actual.city, expected.city)
	}
	if actual.district != expected.district {
		t.Errorf("domria: invalid district, %s != %s", actual.district, expected.district)
	}
	if actual.street != expected.street {
		t.Errorf("domria: invalid street, %s != %s", actual.street, expected.street)
	}
	if actual.houseNumber != expected.houseNumber {
		t.Errorf("domria: invalid house number, %s != %s", actual.houseNumber, expected.houseNumber)
	}
}

func TestFetcherUnmarshalSearchValidItem(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("valid_item"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updatedAt := time.Date(2020, time.June, 6, 14, 57, 18, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-rovno-schastlivoe-chernovola-vyacheslava-ulitsa-16818824.html",
			"dom/photo/10925/1092575/109257503/109257503.jpg",
			&updatedAt,
			27800,
			45,
			0,
			0,
			1,
			2,
			9,
			mining.Primary,
			"ЖК На Щасливому, будинок 27",
			geom.NewPointFlat(geom.XY, []float64{26.267247115344, 50.59766586795}),
			"Рівненська",
			"Рівне",
			"Щасливе",
			"Черновола Вячеслава улица",
			"91-Ф",
		},
	)
}

func TestFetcherUnmarshalSearchEmptyMainPhoto(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("empty_main_photo"), mining.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updatedAt := time.Date(2020, time.June, 3, 16, 16, 26, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-ternopol-bam-17186701.html",
			"",
			&updatedAt,
			20797,
			53.1,
			0,
			12,
			2,
			2,
			4,
			mining.Secondary,
			"",
			geom.NewPointFlat(geom.XY, []float64{25.594767, 49.553517}),
			"Тернопільська",
			"Тернопіль",
			"Бам",
			"",
			"",
		},
	)
}

//func TestFetcherUnmarshalSearchEmptyUpdatedAt(t *testing.T) {}
//
//func TestFetcherUnmarshalSearchTrashUpdatedAt(t *testing.T) {}
//
//func TestFetcherUnmarshalSearchLeadingZerosUpdatedAt(t *testing.T) {}
//
//func TestFetcherUnmarshalSearchMissingShapesUpdatedAt(t *testing.T) {}
//
//func TestFetcherUnmarshalSearch13MonthUpdatedAt(t *testing.T) {}
//
//func TestFetcherUnmarshalSearchShortDateTimingUpdatedAt(t *testing.T) {}
//
//func TestFetcherUnmarshalSearchShortDateUpdatedAt(t *testing.T) {}
//
//func TestFetcherUnmarshalSearchLongTimeUpdatedAt(t *testing.T) {}
