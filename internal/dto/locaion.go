package dto

import (
	"fmt"
	"github.com/paulmach/orb/geojson"
)

type Location struct {
	ID    int           `json:"id"`
	Point geojson.Point `json:"point"`
}

func (location *Location) String() string {
	return fmt.Sprintf("{%d %v}", location.ID, location.Point)
}
