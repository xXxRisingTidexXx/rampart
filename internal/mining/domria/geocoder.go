package domria

import (
	log "github.com/sirupsen/logrus"
	"time"
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
	geocodedNumber, locatedNumber, duration := 0, 0, 0.0
	newFlats := make([]*flat, 0, expectedLength)
	for _, flat := range flats {
		if flat.point != nil {
			newFlats = append(newFlats, flat)
		} else {
			geocodedNumber++
			start := time.Now()
			newFlat, err := geocoder.geocodeFlat(flat)
			duration += time.Since(start).Seconds()
			if err != nil {
				log.Error(err)
			}
			if newFlat != nil {
				locatedNumber++
				newFlats = append(newFlats, newFlat)
			}
		}
	}
	actualLength := len(newFlats)
	log.Debugf(
		"domria: geocoder passed %d flats; %d geocoded, %d located (%.3fs)",
		actualLength,
		geocodedNumber,
		locatedNumber,
		duration/float64(geocodedNumber),
	)
	if lookup := float64(locatedNumber) / float64(geocodedNumber); lookup < geocoder.minLookup {
		log.Warningf("domria: geocoder met low lookup (%.3f)", lookup)
	}
	return newFlats
}

func (geocoder *geocoder) geocodeFlat(f *flat) (*flat, error) {
	return f, nil
}
