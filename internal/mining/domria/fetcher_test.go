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

func TestUnmarshalSearchEmptyString(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte(""), mining.Primary)
	if err == nil || err.Error() != "domria: failed to unmarshal the search, unexpected end of JSON input" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
	if flats != nil {
		t.Errorf("domria: non-empty flats, %v", flats)
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

func TestUnmarshalSearchInvalidJSON(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte("{"), mining.Primary)
	if err == nil || err.Error() != "domria: failed to unmarshal the search, unexpected end of JSON input" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
	if flats != nil {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchArrayInsteadOfObject(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte("[]"), mining.Primary)
	if err == nil || err.Error() != "domria: failed to unmarshal the search"+
		", json: cannot unmarshal array into Go value of type domria.search" {
		t.Errorf("domria: absent or invalid error, %v", err)
	}
	if flats != nil {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchEmptyJSON(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch([]byte("{}"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchWithoutItems(t *testing.T) {
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

func TestUnmarshalSearchEmptySearch(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("empty_search"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if flats == nil || len(flats) != 0 {
		t.Errorf("domria: nil/non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchEmptyItem(t *testing.T) {
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

func TestUnmarshalSearchValidItem(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("valid_item"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 6, 14, 57, 18, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-rovno-schastlivoe-chernovola-vyacheslava-ulitsa-16818824.html",
			"dom/photo/10925/1092575/109257503/109257503.jpg",
			&updateTime,
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

func TestUnmarshalSearchEmptyMainPhoto(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("empty_main_photo"), mining.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 3, 16, 16, 26, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-ternopol-bam-17186701.html",
			"",
			&updateTime,
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

func TestUnmarshalSearchEmptyUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("empty_updated_at"), mining.Secondary)
	if err == nil || err.Error() != "domria: failed to unmarsh"+
		"al the search, domria: moment string is too short, 2" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchTrashUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("trash_updated_at"), mining.Secondary)
	if err == nil || err.Error() != "domria: failed to unmarshal the search, domria: mom"+
		"ent can't split date & timing, |@!|)  )0w23 8&Nu sho, pososesh huj?$@%@8182)( @" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchLeadingZerosUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("leading_zeros_updated_at"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 7, 7, 0, 4, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-harkov-elizavetinskaya-ulitsa-17180614.html",
			"dom/photo/11270/1127013/112701340/112701340.jpg",
			&updateTime,
			23000,
			42,
			21,
			12,
			1,
			7,
			16,
			mining.Primary,
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
	flats, err := fetcher.unmarshalSearch(readAll("missing_shapes_updated_at"), mining.Primary)
	if err == nil || err.Error() != "domria: failed to unmarshal "+
		"the search, domria: moment cannot split date, 2020- 07:53" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearch13MonthUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("13_month_updated_at"), mining.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2021, time.January, 7, 7, 7, 41, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-perevireno-prodaja-kvartira-vinnitsa-vishenka-vasiliya-porika-ulitsa-17073207.html",
			"dom/photo/11162/1116219/111621990/111621990.jpg",
			&updateTime,
			27500,
			32.9,
			32.1,
			6,
			1,
			4,
			5,
			mining.Secondary,
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
	flats, err := fetcher.unmarshalSearch(readAll("just_date_updated_at"), mining.Secondary)
	if err == nil || err.Error() != "domria: failed to unmarshal t"+
		"he search, domria: moment cannot split timing, 2020-06-07 " {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchJustTimeUpdatedAt(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("just_time_updated_at"), mining.Secondary)
	if err == nil || err.Error() != "domria: failed to unmarshal"+
		" the search, domria: moment cannot split date,  07:47:11" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchEmptyPriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("empty_price_arr"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.May, 31, 12, 44, 7, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-chernovtsy-fastovskaya-ruska-17169204.html",
			"dom/photo/11259/1125975/112597577/112597577.jpg",
			&updateTime,
			0,
			86,
			50,
			15,
			3,
			6,
			9,
			mining.Primary,
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
	flats, err := fetcher.unmarshalSearch(readAll("no_usd_price_arr"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 7, 7, 20, 30, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-harkov-shevchenkovskiy-16798175.html",
			"dom/photo/10906/1090623/109062364/109062364.jpg",
			&updateTime,
			0,
			51,
			20,
			12,
			1,
			14,
			16,
			mining.Primary,
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
	flats, err := fetcher.unmarshalSearch(readAll("empty_price_price_arr"), mining.Secondary)
	if err == nil || err.Error() != "domria: failed to unmars"+
		"hal the search, domria: price string is too short, 2" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchWhitespacePricePriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("whitespace_price_price_arr"), mining.Secondary)
	if err == nil || err.Error() != "domria: failed to unmarshal th"+
		"e search, strconv.ParseFloat: parsing \"\": invalid syntax" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchTrashPricePriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("trash_price_price_arr"), mining.Primary)
	if err == nil || err.Error() != "domria: failed to unmarshal the "+
		"search, strconv.ParseFloat: parsing \"Suck\": invalid syntax" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchNegativePricePriceArr(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("negative_price_price_arr"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 7, 7, 55, 45, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-ternopol-bam-saharova-andreya-akademika-ulitsa-16349831.html",
			"dom/photo/10507/1050708/105070868/105070868.jpg",
			&updateTime,
			-38225,
			73.2,
			0,
			7.57,
			3,
			9,
			11,
			mining.Primary,
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
	flats, err := fetcher.unmarshalSearch(readAll("trash_total_square_meters"), mining.Secondary)
	if err == nil || err.Error() != "domria: failed to unmarshal the"+
		" search, invalid character '-' after object key:value pair" {
		t.Fatalf("domria: absent or invalid error, %v", err)
	}
	if len(flats) != 0 {
		t.Errorf("domria: non-empty flats, %v", flats)
	}
}

func TestUnmarshalSearchSupremeKitchenSquareMeters(t *testing.T) {
	fetcher := newDefaultFetcher()
	flats, err := fetcher.unmarshalSearch(readAll("supreme_kitchen_square_meters"), mining.Secondary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 7, 7, 43, 22, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-ujgorod-tsentr-voloshina-ulitsa-15559098.html",
			"dom/photo/9751/975170/97517018/97517018.jpg",
			&updateTime,
			90000,
			96,
			0,
			112,
			4,
			3,
			3,
			mining.Secondary,
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
	flats, err := fetcher.unmarshalSearch(readAll("negative_floor"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 7, 9, 49, 9, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-kiev-solomenskiy-petra-radchenko-ulitsa-16760338.html",
			"dom/photo/10873/1087329/108732937/108732937.jpg",
			&updateTime,
			48510,
			59.32,
			0,
			0,
			2,
			-1,
			26,
			mining.Primary,
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
	flats, err := fetcher.unmarshalSearch(readAll("supreme_floor"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 1, 10, 42, 16, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-kiev-shevchenkovskiy-zlatoustovskaya-ulitsa-16489927.html",
			"dom/photo/10621/1062170/106217048/106217048.jpg",
			&updateTime,
			159300,
			114,
			57,
			16,
			3,
			116,
			18,
			mining.Primary,
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
	flats, err := fetcher.unmarshalSearch(readAll("just_longitude"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 7, 13, 8, 45, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-prodaja-kvartira-vinnitsa-tsentr-lva-tolstogo-ulitsa-17203089.html",
			"dom/photo/11289/1128911/112891158/112891158.jpg",
			&updateTime,
			195000,
			286,
			0,
			0,
			5,
			17,
			17,
			mining.Primary,
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
	flats, err := fetcher.unmarshalSearch(readAll("just_latitude"), mining.Primary)
	if err != nil {
		t.Fatalf("domria: unexpected error, %v", err)
	}
	if len(flats) != 1 {
		t.Fatalf("domria: corrupted flats, %v", flats)
	}
	updateTime := time.Date(2020, time.June, 5, 22, 36, 29, 0, time.Local).UTC()
	assertFlat(
		t,
		flats[0],
		&flat{
			"realty-perevireno-prodaja-kvartira-hmelnitskiy-vyistavka-starokostyantinovskoe-shosse-16982542.html",
			"dom/photo/11243/1124301/112430139/112430139.jpg",
			&updateTime,
			44000,
			50,
			0,
			17.5,
			1,
			4,
			10,
			mining.Primary,
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

//func TestUnmarshalSearchStringCoordinates(t *testing.T) {}
//
//func TestUnmarshalSearchEmptyStringCoordinates(t *testing.T) {}
//
//func TestUnmarshalSearchTrashCoordinates(t *testing.T) {}
//
// Coordinates are out of the allowed range.
//func TestUnmarshalSearchSupremeCoordinates(t *testing.T) {}
//
//func TestUnmarshalSearchEmptyStreets(t *testing.T) {}
//
//func TestUnmarshalSearchJustRUStreet(t *testing.T) {}
