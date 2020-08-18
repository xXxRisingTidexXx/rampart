package dto

import (
	"fmt"
	"github.com/paulmach/orb/geojson"
)

type Location struct {
	OriginURL string        `json:"originURL"`
	Point     geojson.Point `json:"point"`
	City      string        `json:"city"`
}

func (location *Location) String() string {
	return fmt.Sprintf("{%s %v %s}", location.OriginURL, location.Point, location.City)
}
