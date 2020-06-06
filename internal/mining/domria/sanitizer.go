package domria

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/mining/configs"
	"strings"
)

func newSanitizer(config *configs.Sanitizer) *sanitizer {
	return &sanitizer{
		config.OriginURLPrefix,
		config.ImageURLPrefix,
		config.StateEnding,
		config.StateSuffix,
		config.DistrictEnding,
		config.DistrictSuffix,
	}
}

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
	newFlats := make([]*flat, length)
	for i := range flats {
		originURL := flats[i].originURL
		if originURL != "" {
			originURL = sanitizer.originURLPrefix + originURL
		}
		imageURL := flats[i].imageURL
		if imageURL != "" {
			imageURL = sanitizer.imageURLPrefix + imageURL
		}
		state := flats[i].state
		if strings.HasSuffix(state, sanitizer.stateEnding) {
			state += sanitizer.stateSuffix
		}
		district := flats[i].district
		if strings.HasSuffix(district, sanitizer.districtEnding) {
			district += sanitizer.stateSuffix
		}
		newFlats[i] = &flat{
			originURL,
			imageURL,
			flats[i].updateTime,
			flats[i].price,
			flats[i].totalArea,
			flats[i].livingArea,
			flats[i].kitchenArea,
			flats[i].roomNumber,
			flats[i].floor,
			flats[i].totalFloor,
			flats[i].housing,
			flats[i].complex,
			flats[i].point,
			state,
			flats[i].city,
			district,
			flats[i].street,
			flats[i].houseNumber,
		}
	}
	log.Debugf("domria: sanitizer beautified %d flats", length)
	return newFlats
}
