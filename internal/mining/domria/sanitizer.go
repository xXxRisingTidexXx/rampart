package domria

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/config"
	"strings"
)

func newSanitizer(config *config.Sanitizer) *sanitizer {
	return &sanitizer{
		config.OriginURLPrefix,
		config.ImageURLPrefix,
		config.StateEnding,
		config.StateSuffix,
		config.DistrictEnding,
		config.DistrictSuffix,
	}
}

// TODO: add street split (some streets contain house numbers).
// TODO: add street purity (replace shorts and russian suffixes and endings).
// TODO: add complex purity (replace "ЖК " prefixes if needed).
type sanitizer struct {
	originURLPrefix string
	imageURLPrefix  string
	stateEnding     string
	stateSuffix     string
	districtEnding  string
	districtSuffix  string
}

func (sanitizer *sanitizer) sanitizeFlats(flats []*flat) []*flat {
	length := len(flats)
	if length == 0 {
		log.Debug("domria: sanitizer skipped flats")
		return flats
	}
	newFlats := make([]*flat, length)
	for i, flat := range flats {
		newFlats[i] = sanitizer.sanitizeFlat(flat)
	}
	log.Debugf("domria: sanitizer sanitized %d flats", length)
	return newFlats
}

func (sanitizer *sanitizer) sanitizeFlat(f *flat) *flat {
	originURL := f.originURL
	if originURL != "" {
		originURL = sanitizer.originURLPrefix + originURL
	}
	imageURL := f.imageURL
	if imageURL != "" {
		imageURL = sanitizer.imageURLPrefix + imageURL
	}
	state := f.state
	if strings.HasSuffix(state, sanitizer.stateEnding) {
		state += sanitizer.stateSuffix
	}
	district := f.district
	if strings.HasSuffix(district, sanitizer.districtEnding) {
		district += sanitizer.districtSuffix
	}
	return &flat{
		originURL,
		imageURL,
		f.updateTime,
		f.price,
		f.totalArea,
		f.livingArea,
		f.kitchenArea,
		f.roomNumber,
		f.floor,
		f.totalFloor,
		f.housing,
		f.complex,
		f.point,
		state,
		f.city,
		district,
		f.street,
		f.houseNumber,
	}
}
