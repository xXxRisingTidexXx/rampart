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
	for i := range flats {
		var newFlat *flat
		if flats[i].point != nil {
			newFlat = flats[i]
		} else {
			geocodedNumber++
			start := time.Now()
			bytes, err := geocoder.getLocations(flats[i])
			duration += time.Since(start).Seconds()
			if err == nil {
				newFlat, err = geocoder.locateFlat(flats[i], bytes)
			}
		}
		if newFlat != nil {
			locatedNumber++
			newFlats = append(newFlats, newFlat)
			newFlat = nil // TODO: is redundant?
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

func (geocoder *geocoder) getLocations(flat *flat) ([]byte, error) {
	return nil, nil
}

func (geocoder *geocoder) locateFlat(f *flat, bytes []byte) (*flat, error) {
	return nil, nil
}
