package domria

import (
	log "github.com/sirupsen/logrus"
)

func newGeocoder() *geocoder {
	return &geocoder{
		"https://nominatim.openstreetmap.org/search?city=%s&c" +
			"ounty=%s&street=%s+%s&format=json&countrycodes=ua",
		0.7,
	}
}

type geocoder struct {
	searchURL string
	minLookup float64
}

func (geocoder *geocoder) geocodeFlats(flats []*flat) []*flat {
	expectedLength := len(flats)
	if expectedLength == 0 {
		log.Debug("domria: geocoder skipped flats")
		return flats
	}
	geocodedNumber, locatedNumber, avgDuration := 0, 0, 0.0
	newFlats := make([]*flat, 0, expectedLength)

	actualLength := len(newFlats)
	log.Debugf(
		"domria: geocoder passed %d flats; %d geocoded, %d located (%.3f)",
		actualLength,
		geocodedNumber,
		locatedNumber,
		avgDuration / float64(geocodedNumber),
	)
	if lookup := float64(locatedNumber) / float64(geocodedNumber); lookup < geocoder.minLookup {
		log.Warningf("domria: geocoder met low lookup (%.3f)", lookup)
	}
	return newFlats
}
