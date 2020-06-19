package domria

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rampart/internal/mining/config"
	"rampart/internal/mining/misc"
	"testing"
	"time"
)

func TestFetchSearchInvalidHousing(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.fetchFlats(misc.Secondary)
	if err == nil || err.Error() != "domria: fetcher doesn't accept secondary housing" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if flats != nil {
		t.Errorf("domria: non-nil flats, %v", flats)
	}
}

func newDefaultFetcher() *fetcher {
	return newTestFetcher("https://domria.ua/search/")
}

func newTestFetcher(searchURL string) *fetcher {
	return newFetcher(
		&config.Fetcher{
			Timeout:   misc.Timeout(100 * time.Millisecond),
			Portion:   10,
			Flags:     map[misc.Housing]string{misc.Primary: "pm_housing=1"},
			Headers:   map[string]string{"User-Agent": "domria-test-bot/1.0.0"},
			SearchURL: searchURL,
		},
	)
}

func TestFetchSearchWithReset(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				expected := "domria-test-bot/1.0.0"
				if actual := request.Header.Get("User-Agent"); actual != expected {
					t.Fatalf("domria: invalid user-agent, %s != %s", actual, expected)
				}
				expected = "/?pm_housing=1&page=2&limit=10"
				if actual := request.URL.String(); actual != expected {
					t.Fatalf("domria: invalid url, %s != %s", actual, expected)
				}
				if _, err := fmt.Fprint(writer, "{\"count\":0,\"items\":[]}\n"); err != nil {
					t.Fatalf("domria: unexpected error, %v", err)
				}
			},
		),
	)
	fetcher := newServerFetcher(server)
	fetcher.page = 2
	flats, err := fetcher.fetchFlats(misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if fetcher.page != 0 {
		t.Errorf("domria: fetcher left on page %d", fetcher.page)
	}
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: corrupted flats, %v", flats)
	}
	server.Close()
}

func newServerFetcher(server *httptest.Server) *fetcher {
	return newTestFetcher(server.URL + "/?%s&page=%d&limit=%d")
}

//nolint:funlen
func TestFetchSearchMultipleFlats(t *testing.T) {
	server := newServer(
		t,
		func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			if _, err := writer.Write(readAll(t, "multiple_flats")); err != nil {
				t.Fatalf("domria: unexpected error, %v", err)
			}
		},
	)
	fetcher := newServerFetcher(server)
	flats, err := fetcher.fetchFlats(misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if fetcher.page != 1 {
		t.Errorf("domria: fetcher didn't inc the page, %d", fetcher.page)
	}
	if len(flats) != 2 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-vinnitsa-blijnee-zamoste-16890016.html",
			"dom/photo/10990/1099054/109905467/109905467.jpg",
			time.Date(2020, time.June, 8, 6, 59, 6, 0, time.Local).UTC(),
			42333,
			66.76,
			0,
			0,
			3,
			8,
			9,
			misc.Primary,
			"",
			geom.NewPointFlat(geom.XY, []float64{28.4962815, 49.2410151}),
			"Вінницька",
			"Вінниця",
			"Ближнє замостя",
			"вул.Острозького",
			"",
		},
	)
	assertFlat(
		t,
		flats[1],
		&flat{
			"realty-prodaja-kvartira-vinnitsa-akademicheskiy-16892143.html",
			"dom/photo/10992/1099221/109922120/109922120.jpg",
			time.Date(2020, time.June, 8, 6, 58, 49, 0, time.Local).UTC(),
			36000,
			45.5,
			0,
			0,
			1,
			-7,
			7,
			misc.Primary,
			"Микрорайон «АКАДЕМІЧНИЙ»",
			geom.NewPointFlat(geom.XY, []float64{28.4269, 49.207109}),
			"Вінницька",
			"Вінниця",
			"Академічний",
			"вул. Миколаївська / вул. Тимофіївська",
			"",
		},
	)
	server.Close()
}

func newServer(t *testing.T, handler func(http.ResponseWriter, *http.Request)) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				expected := "domria-test-bot/1.0.0"
				if actual := request.Header.Get("User-Agent"); actual != expected {
					t.Fatalf("domria: invalid user-agent, %s != %s", actual, expected)
				}
				expected = "/?pm_housing=1&page=0&limit=10"
				if actual := request.URL.String(); actual != expected {
					t.Fatalf("domria: invalid url, %s != %s", actual, expected)
				}
				handler(writer, request)
			},
		),
	)
}

func readAll(t *testing.T, fixtureName string) []byte {
	bytes, err := ioutil.ReadFile("testdata/fetcher/" + fixtureName + ".json")
	if err != nil {
		t.Fatalf("domria: failed to read the file, %v", err)
	}
	return bytes
}

//nolint:gocognit,gocyclo,funlen
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
	if !actual.updateTime.Equal(expected.updateTime) {
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

func TestGetSearchEmptySearch(t *testing.T) {
	expected := "{\"count\":0,\"items\":[]}\n"
	server := newServer(
		t,
		func(writer http.ResponseWriter, _ *http.Request) {
			if _, err := fmt.Fprint(writer, expected); err != nil {
				t.Fatalf("domria: unexpected error, %v", err)
			}
		},
	)
	fetcher := newServerFetcher(server)
	bytes, err := fetcher.getSearch("pm_housing=1")
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if actual := string(bytes); actual != expected {
		t.Errorf("domria: invalid search, %s != %s", actual, expected)
	}
	server.Close()
}

func TestGetSearchWithTimeout(t *testing.T) {
	server := newServer(
		t,
		func(_ http.ResponseWriter, _ *http.Request) {
			time.Sleep(110 * time.Millisecond)
		},
	)
	fetcher := newServerFetcher(server)
	bytes, err := fetcher.getSearch("pm_housing=1")
	if err == nil || err.Error() != "domria: fetcher failed to perform a request, Get \""+
		server.URL+
		"/?pm_housing=1&page=0&limit=10\": context deadline exc"+
		"eeded (Client.Timeout exceeded while awaiting headers)" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if bytes != nil {
		t.Errorf("domria: non-nil bytes, %v", bytes)
	}
	server.Close()
}

func TestGetSearchNotFound(t *testing.T) {
	server := newServer(t, http.NotFound)
	fetcher := newServerFetcher(server)
	bytes, err := fetcher.getSearch("pm_housing=1")
	if err == nil || err.Error() != "domria: fetcher got response with status 404 Not Found" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
	if bytes != nil {
		t.Errorf("domria: non-nil bytes, %v", bytes)
	}
	server.Close()
}

func TestUnmarshalSearchEmptyString(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte(""), misc.Primary)
	if err == nil || err.Error() != "domria: fetcher failed t"+
		"o unmarshal the search, unexpected end of JSON input" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
	if flats != nil {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchInvalidJSON(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte("{"), misc.Primary)
	if err == nil || err.Error() != "domria: fetcher failed t"+
		"o unmarshal the search, unexpected end of JSON input" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
	if flats != nil {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchArrayInsteadOfObject(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte("[]"), misc.Primary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarshal the search"+
		", json: cannot unmarshal array into Go value of type domria.search" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
	if flats != nil {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchEmptyJSON(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte("{}"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchWithoutItems(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "without_items"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchEmptySearch(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "empty_search"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchEmptyItem(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "empty_item"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(t, flats[0], &flat{housing: misc.Primary})
}

func TestUnmarshalSearchValidItem(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "valid_item"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-rovno-schastlivoe-chernovola-vyacheslava-ulitsa-16818824.html",
			"dom/photo/10925/1092575/109257503/109257503.jpg",
			time.Date(2020, time.June, 6, 14, 57, 18, 0, time.Local).UTC(),
			27800,
			45,
			0,
			0,
			1,
			2,
			9,
			misc.Primary,
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

func TestUnmarshalSearchEmptyMainPhoto(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "empty_main_photo"), misc.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-ternopol-bam-17186701.html",
			"",
			time.Date(2020, time.June, 3, 16, 16, 26, 0, time.Local).UTC(),
			20797,
			53.1,
			0,
			12,
			2,
			2,
			4,
			misc.Secondary,
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

func TestUnmarshalSearchEmptyUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "empty_updated_at"), misc.Secondary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarsh"+
		"al the search, domria: moment string is too short, 2" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchTrashUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "trash_updated_at"), misc.Secondary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarshal the search, domria: mom"+
		"ent can't split date & timing, |@!|)  )0w23 8&Nu sho, pososesh huj?$@%@8182)( @" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchLeadingZerosUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "leading_zeros_updated_at"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-harkov-elizavetinskaya-ulitsa-17180614.html",
			"dom/photo/11270/1127013/112701340/112701340.jpg",
			time.Date(2020, time.June, 7, 7, 0, 4, 0, time.Local).UTC(),
			23000,
			42,
			21,
			12,
			1,
			7,
			16,
			misc.Primary,
			"ЖК Левада 2",
			geom.NewPointFlat(geom.XY, []float64{36.239501354492, 49.978100188645}),
			"Харківська",
			"Харків",
			"",
			"Лисаветинська вулиця",
			"2б",
		},
	)
}

func TestUnmarshalSearchMissingShapesUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "missing_shapes_updated_at"), misc.Primary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarshal "+
		"the search, domria: moment cannot split date, 2020- 07:53" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearch13MonthUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "13_month_updated_at"), misc.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-perevireno-prodaja-kvartira-vinnitsa-vishenka-vasiliya-porika-ulitsa-17073207.html",
			"dom/photo/11162/1116219/111621990/111621990.jpg",
			time.Date(2021, time.January, 7, 7, 7, 41, 0, time.Local).UTC(),
			27500,
			32.9,
			32.1,
			6,
			1,
			4,
			5,
			misc.Secondary,
			"",
			geom.NewPointFlat(geom.XY, []float64{28.4247279, 49.2291492}),
			"Вінницька",
			"Вінниця",
			"Вишенька",
			"Василя Порика вулиця",
			"1",
		},
	)
}

func TestUnmarshalSearchJustDateUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "just_date_updated_at"), misc.Secondary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarshal t"+
		"he search, domria: moment cannot split timing, 2020-06-07 " {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchJustTimeUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "just_time_updated_at"), misc.Secondary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarshal"+
		" the search, domria: moment cannot split date,  07:47:11" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchEmptyPriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "empty_price_arr"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-chernovtsy-fastovskaya-ruska-17169204.html",
			"dom/photo/11259/1125975/112597577/112597577.jpg",
			time.Date(2020, time.May, 31, 12, 44, 7, 0, time.Local).UTC(),
			0,
			86,
			50,
			15,
			3,
			6,
			9,
			misc.Primary,
			"",
			geom.NewPointFlat(geom.XY, []float64{25.9820274, 48.2831323}),
			"Чернівецька",
			"Чернівці",
			"Фастівська",
			"Руська",
			"223Д",
		},
	)
}

func TestUnmarshalSearchNoUSDPriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "no_usd_price_arr"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-harkov-shevchenkovskiy-16798175.html",
			"dom/photo/10906/1090623/109062364/109062364.jpg",
			time.Date(2020, time.June, 7, 7, 20, 30, 0, time.Local).UTC(),
			0,
			51,
			20,
			12,
			1,
			14,
			16,
			misc.Primary,
			"",
			geom.NewPointFlat(geom.XY, []float64{36.2245388, 49.9974272}),
			"Харківська",
			"Харків",
			"Шевченківський",
			"пр Ботаническая",
			"2",
		},
	)
}

func TestUnmarshalSearchEmptyPricePriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "empty_price_price_arr"), misc.Secondary)
	if err == nil || err.Error() != "domria: fetcher failed to unmars"+
		"hal the search, domria: price string is too short, 2" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchWhitespacePricePriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "whitespace_price_price_arr"), misc.Secondary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarshal th"+
		"e search, strconv.ParseFloat: parsing \"\": invalid syntax" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchTrashPricePriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "trash_price_price_arr"), misc.Primary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarshal the "+
		"search, strconv.ParseFloat: parsing \"Suck\": invalid syntax" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchNegativePricePriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "negative_price_price_arr"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-ternopol-bam-saharova-andreya-akademika-ulitsa-16349831.html",
			"dom/photo/10507/1050708/105070868/105070868.jpg",
			time.Date(2020, time.June, 7, 7, 55, 45, 0, time.Local).UTC(),
			-38225,
			73.2,
			0,
			7.57,
			3,
			9,
			11,
			misc.Primary,
			"",
			geom.NewPointFlat(geom.XY, []float64{25.644687974235, 49.550329822762}),
			"Тернопільська",
			"Тернопіль",
			"Бам",
			"Сахарова Андрія Академіка вулиця",
			"10",
		},
	)
}

func TestUnmarshalSearchTrashTotalSquareMeters(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "trash_total_square_meters"), misc.Secondary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarshal the"+
		" search, invalid character '-' after object key:value pair" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchSupremeKitchenSquareMeters(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "supreme_kitchen_square_meters"), misc.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-ujgorod-tsentr-voloshina-ulitsa-15559098.html",
			"dom/photo/9751/975170/97517018/97517018.jpg",
			time.Date(2020, time.June, 7, 7, 43, 22, 0, time.Local).UTC(),
			90000,
			96,
			0,
			112,
			4,
			3,
			3,
			misc.Secondary,
			"",
			geom.NewPointFlat(geom.XY, []float64{22.301875199999998, 48.621579}),
			"Закарпатська",
			"Ужгород",
			"Центр",
			"Волошина вулиця",
			"",
		},
	)
}

func TestUnmarshalSearchNegativeFloor(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "negative_floor"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-kiev-solomenskiy-petra-radchenko-ulitsa-16760338.html",
			"dom/photo/10873/1087329/108732937/108732937.jpg",
			time.Date(2020, time.June, 7, 9, 49, 9, 0, time.Local).UTC(),
			48510,
			59.32,
			0,
			0,
			2,
			-1,
			26,
			misc.Primary,
			"ЖК Медовий-2",
			geom.NewPointFlat(geom.XY, []float64{30.4760253, 50.4128865}),
			"Київська",
			"Київ",
			"Солом'янський",
			"Петра Радченко улица",
			"27",
		},
	)
}

func TestUnmarshalSearchSupremeFloor(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "supreme_floor"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-kiev-shevchenkovskiy-zlatoustovskaya-ulitsa-16489927.html",
			"dom/photo/10621/1062170/106217048/106217048.jpg",
			time.Date(2020, time.June, 1, 10, 42, 16, 0, time.Local).UTC(),
			159300,
			114,
			57,
			16,
			3,
			116,
			18,
			misc.Primary,
			"ЖК «Шевченківський»",
			geom.NewPointFlat(geom.XY, []float64{30.487440507934, 50.450000744175}),
			"Київська",
			"Київ",
			"Шевченківський",
			"Золотоустівська вулиця",
			"27",
		},
	)
}

func TestUnmarshalSearchJustLongitude(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "just_longitude"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-vinnitsa-tsentr-lva-tolstogo-ulitsa-17203089.html",
			"dom/photo/11289/1128911/112891158/112891158.jpg",
			time.Date(2020, time.June, 7, 13, 8, 45, 0, time.Local).UTC(),
			195000,
			286,
			0,
			0,
			5,
			17,
			17,
			misc.Primary,
			"",
			geom.NewPointFlat(geom.XY, []float64{28.4607622, 0}),
			"Вінницька",
			"Вінниця",
			"Центр",
			"Льва Толстого вулиця",
			"9",
		},
	)
}

func TestUnmarshalSearchJustLatitude(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "just_latitude"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-perevireno-prodaja-kvartira-hmelnitskiy-vyistavka-starokostyantinovskoe-shosse-16982542.html",
			"dom/photo/11243/1124301/112430139/112430139.jpg",
			time.Date(2020, time.June, 5, 22, 36, 29, 0, time.Local).UTC(),
			44000,
			50,
			0,
			17.5,
			1,
			4,
			10,
			misc.Primary,
			"",
			geom.NewPointFlat(geom.XY, []float64{0, 49.431359}),
			"Хмельницька",
			"Хмельницький",
			"Виставка",
			"Старокостянтинівське шосе",
			"20/7",
		},
	)
}

func TestUnmarshalSearchStringCoordinates(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "string_coordinates"), misc.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-vinnitsa-podole-17135787.html",
			"dom/photo/11226/1122631/112263193/112263193.jpg",
			time.Date(2020, time.June, 2, 9, 38, 5, 0, time.Local).UTC(),
			129000,
			108,
			62,
			13.4,
			4,
			3,
			9,
			misc.Secondary,
			"Микрорайон Поділля",
			geom.NewPointFlat(geom.XY, []float64{28.4489892, 49.2173192}),
			"Вінницька",
			"Вінниця",
			"Поділля",
			"вул. Зодчих / вул. Бортняка",
			"1",
		},
	)
}

func TestUnmarshalSearchEmptyStringCoordinates(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "empty_string_coordinates"), misc.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-odessa-primorskiy-nekrasova-pereulok-16179973.html",
			"dom/photo/10370/1037099/103709962/103709962.jpg",
			time.Date(2020, time.June, 8, 5, 39, 24, 0, time.Local).UTC(),
			199000,
			145,
			78,
			27,
			4,
			2,
			2,
			misc.Secondary,
			"",
			nil,
			"Одеська",
			"Одеса",
			"Приморський",
			"Некрасова провулок",
			"",
		},
	)
}

func TestUnmarshalSearchTrashCoordinates(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "trash_coordinates"), misc.Secondary)
	if err == nil || err.Error() != "domria: fetcher failed to unmarshal the sear"+
		"ch, strconv.ParseFloat: parsing \"982jd293jd)J\": invalid syntax" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchSupremeCoordinates(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "supreme_coordinates"), misc.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-herson-suvorovskiy-17165402.html",
			"dom/photo/11256/1125653/112565321/112565321.jpg",
			time.Date(2020, time.June, 8, 10, 9, 58, 0, time.Local).UTC(),
			55000,
			72,
			0,
			0,
			3,
			2,
			7,
			misc.Secondary,
			"",
			geom.NewPointFlat(geom.XY, []float64{-183.839023, 2931.000183399}),
			"Херсонська",
			"Херсон",
			"Суворовський",
			"200 Лет Херсона пр.",
			"",
		},
	)
}

func TestUnmarshalSearchEmptyStreets(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "empty_streets"), misc.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-odessa-kievskiy-ilfa-i-petrova-ulitsa-17120761.html",
			"dom/photo/11211/1121120/112112031/112112031.jpg",
			time.Date(2020, time.June, 8, 6, 7, 59, 0, time.Local).UTC(),
			37500,
			63,
			38,
			10,
			3,
			4,
			9,
			misc.Secondary,
			"",
			nil,
			"Одеська",
			"Одеса",
			"Київський",
			"",
			"",
		},
	)
}

func TestUnmarshalSearchJustRUStreet(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "just_ru_street"), misc.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-lvov-galitskiy-17148133.html",
			"dom/photo/11238/1123874/112387482/112387482.jpg",
			time.Date(2020, time.June, 8, 7, 30, 43, 0, time.Local).UTC(),
			79000,
			59,
			41,
			9,
			2,
			3,
			3,
			misc.Secondary,
			"",
			nil,
			"Львівська",
			"Львів",
			"Галицький",
			"С. Томашівського",
			"",
		},
	)
}

//nolint:funlen
func TestUnmarshalSearchMultipleItems(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll(t, "multiple_items"), misc.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 3 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-vinnitsa-podole-generala-yakova-gandzyuka-ulitsa-17150263.html",
			"dom/photo/11241/1124150/112415070/112415070.jpg",
			time.Date(2020, time.June, 8, 6, 59, 13, 0, time.Local).UTC(),
			42000,
			63,
			0,
			10,
			2,
			2,
			9,
			misc.Primary,
			"ЖК Перлина Поділля",
			geom.NewPointFlat(geom.XY, []float64{28.437752173707, 49.214143792302}),
			"Вінницька",
			"Вінниця",
			"Поділля",
			"генерала Якова Гандзюка вулиця",
			"6",
		},
	)
	assertFlat(
		t,
		flats[1],
		&flat{
			"realty-prodaja-kvartira-dnepr-slobojanskoe-slobojanskiy-prospekt-16927270.html",
			"dom/photo/11025/1102580/110258034/110258034.jpg",
			time.Date(2018, time.June, 8, 10, 7, 18, 0, time.Local).UTC(),
			31928,
			67.4,
			0,
			0,
			1,
			8,
			10,
			misc.Primary,
			"ЖК Дніпровська Брама 2",
			geom.NewPointFlat(geom.XY, []float64{35.085059977507, 48.536070034556}),
			"Дніпропетровська",
			"Дніпро",
			"Слобожанське",
			"Слобожанский проспект",
			"",
		},
	)
	assertFlat(
		t,
		flats[2],
		&flat{
			"realty-prodaja-kvartira-dnepr-slobojanskoe-slobojanskiy-prospekt-16927282.html",
			"dom/photo/11025/1102580/110258071/110258071.jpg",
			time.Date(2020, time.June, 8, 10, 7, 18, 0, time.Local).UTC(),
			21168,
			45.4,
			0,
			0,
			1,
			6,
			10,
			misc.Primary,
			"ЖК Дніпровська Брама 2",
			geom.NewPointFlat(geom.XY, []float64{35.085059977507, 48.536070034556}),
			"Дніпропетровська",
			"Дніпро",
			"Слобожанське",
			"Слобожанский проспект",
			"",
		},
	)
}
