package dto

import (
	"fmt"
	"github.com/paulmach/orb/geojson"
)

type Location struct {
	OriginURL string        `json:"originURL"`
	Point     geojson.Point `json:"point"`
}

func (location *Location) String() string {
	return fmt.Sprintf("{%s %v}", location.OriginURL, location.Point)
}
