package domria

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"strings"
)

func NewSanitizer(
	config config.Sanitizer,
	drain *metrics.Drain,
	logger log.FieldLogger,
) *Sanitizer {
	return &Sanitizer{
		config.URLPrefix,
		config.PhotoFormat,
		config.PanoramaPrefix,
		config.PanoramaSuffix,
		config.StateMap,
		config.StateSuffix,
		config.CityMap,
		config.DistrictMap,
		config.DistrictCitySwaps,
		config.DistrictEnding,
		config.DistrictSuffix,
		strings.NewReplacer(config.StreetReplacements...),
		strings.NewReplacer(config.HouseNumberReplacements...),
		config.HouseNumberMaxLength,
		drain,
		logger,
	}
}

type Sanitizer struct {
	urlPrefix            string
	photoFormat          string
	panoramaPrefix       string
	panoramaSuffix       string
	stateMap             map[string]string
	stateSuffix          string
	cityMap              map[string]string
	districtMap          map[string]string
	districtCitySwaps    misc.Set
	districtEnding       string
	districtSuffix       string
	streetReplacer       *strings.Replacer
	houseNumberReplacer  *strings.Replacer
	houseNumberMaxLength int
	drain                *metrics.Drain
	logger               log.FieldLogger
}

func (s *Sanitizer) SanitizeFlats(flats []Flat) []Flat {
	newFlats := make([]Flat, len(flats))
	for i, flat := range flats {
		newFlats[i] = s.sanitizeFlat(flat)
	}
	return newFlats
}

func (s *Sanitizer) sanitizeFlat(flat Flat) Flat {
	url := flat.URL
	if url != "" {
		url = s.urlPrefix + url
	} else {
		s.logger.WithField("source", flat.Source).Warning(
			"domria: sanitizer found a flat without an url",
		)
	}
	photos := make([]string, 0)
	if index := strings.LastIndex(flat.URL, "-"); index != -1 {
		slug := flat.URL[:index]
		for _, p := range flat.Photos {
			photos = append(photos, fmt.Sprintf(s.photoFormat, slug, p))
		}
	} else if flat.URL != "" {
		s.logger.WithField("source", flat.Source).Warning(
			"domria: sanitizer found a flat without a dash in the url",
		)
	}
	panoramas := make([]string, len(flat.Panoramas))
	for i, p := range flat.Panoramas {
		panoramas[i] = s.panoramaPrefix +
			strings.ReplaceAll(p, ".jpg", s.panoramaSuffix)
	}
	state := strings.TrimSpace(flat.State)
	if value, ok := s.stateMap[state]; ok {
		state = value
		s.drain.DrainNumber(metrics.StateSanitationNumber)
	}
	if state != "" {
		state += s.stateSuffix
	}
	city := strings.TrimSpace(flat.City)
	if value, ok := s.cityMap[city]; ok {
		city = value
		s.drain.DrainNumber(metrics.CitySanitationNumber)
	}
	district := strings.TrimSpace(flat.District)
	if value, ok := s.districtMap[district]; ok {
		district = value
		s.drain.DrainNumber(metrics.DistrictSanitationNumber)
	}
	if s.districtCitySwaps.Contains(city) {
		city, district = district, ""
		s.drain.DrainNumber(metrics.SwapSanitationNumber)
	}
	if strings.HasSuffix(district, s.districtEnding) {
		district += s.districtSuffix
	}
	street, houseNumber := flat.Street, s.sanitizeHouseNumber(flat.HouseNumber)
	if index := strings.Index(flat.Street, ","); index != -1 {
		street = flat.Street[:index]
		s.drain.DrainNumber(metrics.StreetSanitationNumber)
		extraNumber := s.sanitizeHouseNumber(flat.Street[index+1:])
		if houseNumber == "" &&
			extraNumber != "" &&
			extraNumber[0] >= '0' &&
			extraNumber[0] <= '9' {
			houseNumber = extraNumber
			s.drain.DrainNumber(metrics.HouseNumberSanitationNumber)
		}
	}
	if runes := []rune(houseNumber); len(runes) > s.houseNumberMaxLength {
		houseNumber = string(runes[:s.houseNumberMaxLength])
	}
	return Flat{
		Source:      flat.Source,
		URL:         url,
		Photos:      photos,
		Panoramas:   panoramas,
		UpdateTime:  flat.UpdateTime,
		Price:       flat.Price,
		TotalArea:   flat.TotalArea,
		LivingArea:  flat.LivingArea,
		KitchenArea: flat.KitchenArea,
		RoomNumber:  flat.RoomNumber,
		Floor:       flat.Floor,
		TotalFloor:  flat.TotalFloor,
		Housing:     flat.Housing,
		Complex:     flat.Complex,
		Point:       flat.Point,
		State:       state,
		City:        city,
		District:    district,
		Street:      strings.TrimSpace(s.streetReplacer.Replace(street)),
		HouseNumber: houseNumber,
	}
}

func (s *Sanitizer) sanitizeHouseNumber(houseNumber string) string {
	if houseNumber == "" {
		return houseNumber
	}
	newHouseNumber := s.houseNumberReplacer.Replace(houseNumber)
	if index := strings.Index(newHouseNumber, ","); index != -1 {
		return newHouseNumber[:index]
	}
	return newHouseNumber
}
