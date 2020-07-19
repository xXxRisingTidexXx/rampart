package domria

import (
	"rampart/internal/config"
	"strings"
)

func NewSanitizer(config *config.Sanitizer) *Sanitizer {
	return &Sanitizer{
		config.OriginURLPrefix,
		config.ImageURLPrefix,
		config.StateEnding,
		config.StateSuffix,
		config.DistrictEnding,
		config.DistrictSuffix,
	}
}

type Sanitizer struct {
	originURLPrefix string
	imageURLPrefix  string
	stateEnding     string
	stateSuffix     string
	districtEnding  string
	districtSuffix  string
}

func (sanitizer *Sanitizer) SanitizeFlats(flats []*Flat) []*Flat {
	newFlats := make([]*Flat, len(flats))
	for i, flat := range flats {
		newFlats[i] = sanitizer.sanitizeFlat(flat)
	}
	return newFlats
}

func (sanitizer *Sanitizer) sanitizeFlat(f *Flat) *Flat {
	originURL := f.OriginURL
	if originURL != "" {
		originURL = sanitizer.originURLPrefix + originURL
	}
	imageURL := f.ImageURL
	if imageURL != "" {
		imageURL = sanitizer.imageURLPrefix + imageURL
	}
	state := f.State
	if strings.HasSuffix(state, sanitizer.stateEnding) {
		state += sanitizer.stateSuffix
	}
	district := f.District
	if strings.HasSuffix(district, sanitizer.districtEnding) {
		district += sanitizer.districtSuffix
	}
	return &Flat{
		originURL,
		imageURL,
		f.UpdateTime,
		f.Price,
		f.TotalArea,
		f.LivingArea,
		f.KitchenArea,
		f.RoomNumber,
		f.Floor,
		f.TotalFloor,
		f.Housing,
		f.Complex,
		f.Point,
		state,
		f.City,
		district,
		f.Street,
		f.HouseNumber,
	}
}
