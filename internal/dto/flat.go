package dto

import (
	"fmt"
	"github.com/paulmach/orb/geojson"
)

type Flat struct {
	OriginURL string        `json:"originURL"`
	Point     geojson.Point `json:"point"`
	City      string        `json:"city"`
}

func (flat *Flat) String() string {
	return fmt.Sprintf("{%s %v %s}", flat.OriginURL, flat.Point, flat.City)
}
