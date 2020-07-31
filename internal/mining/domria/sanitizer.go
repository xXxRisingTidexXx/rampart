package domria

import (
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"strings"
)

func NewSanitizer(config *config.Sanitizer, gatherer *metrics.Gatherer) *Sanitizer {
	return &Sanitizer{
		config.OriginURLPrefix,
		config.ImageURLPrefix,
		config.StateDictionary,
		config.StateSuffix,
		config.CityDictionary,
		config.DistrictDictionary,
		config.DistrictCitySwaps,
		config.DistrictEnding,
		config.DistrictSuffix,
		strings.NewReplacer(config.StreetReplacements...),
		strings.NewReplacer(config.HouseNumberReplacements...),
		",",
		gatherer,
	}
}

type Sanitizer struct {
	originURLPrefix     string
	imageURLPrefix      string
	stateDictionary     map[string]string
	stateSuffix         string
	cityDictionary      map[string]string
	districtDictionary  map[string]string
	districtCitySwaps   *misc.Set
	districtEnding      string
	districtSuffix      string
	streetReplacer      *strings.Replacer
	houseNumberReplacer *strings.Replacer
	comma               string
	gatherer            *metrics.Gatherer
}

func (sanitizer *Sanitizer) SanitizeFlats(flats []*Flat) []*Flat {
	newFlats := make([]*Flat, len(flats))
	for i, flat := range flats {
		newFlats[i] = sanitizer.sanitizeFlat(flat)
	}
	return newFlats
}

func (sanitizer *Sanitizer) sanitizeFlat(flat *Flat) *Flat {
	originURL := flat.OriginURL
	if originURL != "" {
		originURL = sanitizer.originURLPrefix + originURL
	}
	imageURL := flat.ImageURL
	if imageURL != "" {
		imageURL = sanitizer.imageURLPrefix + imageURL
	}
	state := strings.TrimSpace(flat.State)
	if value, ok := sanitizer.stateDictionary[state]; ok {
		state = value
		sanitizer.gatherer.GatherStateSanitization()
	}
	if state != "" {
		state += sanitizer.stateSuffix
	}
	city := strings.TrimSpace(flat.City)
	if value, ok := sanitizer.cityDictionary[city]; ok {
		city = value
		sanitizer.gatherer.GatherCitySanitization()
	}
	district := strings.TrimSpace(flat.District)
	if value, ok := sanitizer.districtDictionary[district]; ok {
		district = value
		sanitizer.gatherer.GatherDistrictSanitization()
	}
	if sanitizer.districtCitySwaps.Contains(city) {
		city, district = district, city
		sanitizer.gatherer.GatherSwapSanitization()
	}
	if strings.HasSuffix(district, sanitizer.districtEnding) {
		district += sanitizer.districtSuffix
	}
	street, houseNumber := flat.Street, sanitizer.sanitizeHouseNumber(flat.HouseNumber)
	if index := strings.Index(flat.Street, sanitizer.comma); index != -1 {
		street = flat.Street[:index]
		sanitizer.gatherer.GatherStreetSanitization()
		extraNumber := sanitizer.sanitizeHouseNumber(flat.Street[index+1:])
		if houseNumber == "" && extraNumber != "" && extraNumber[0] >= '0' && extraNumber[0] <= '9' {
			houseNumber = extraNumber
			sanitizer.gatherer.GatherHouseNumberSanitization()
		}
	}
	return &Flat{
		originURL,
		imageURL,
		flat.UpdateTime,
		flat.Price,
		flat.TotalArea,
		flat.LivingArea,
		flat.KitchenArea,
		flat.RoomNumber,
		flat.Floor,
		flat.TotalFloor,
		flat.Housing,
		flat.Complex,
		flat.Point,
		0,
		state,
		city,
		district,
		strings.TrimSpace(sanitizer.streetReplacer.Replace(street)),
		houseNumber,
	}
}

func (sanitizer *Sanitizer) sanitizeHouseNumber(houseNumber string) string {
	if houseNumber == "" {
		return houseNumber
	}
	newHouseNumber := sanitizer.houseNumberReplacer.Replace(houseNumber)
	if index := strings.Index(newHouseNumber, sanitizer.comma); index != -1 {
		return newHouseNumber[:index]
	}
	return newHouseNumber
}
