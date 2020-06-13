package domria

import (
	"fmt"
)

type location struct {
	Lon coordinate `json:"lon"`
	Lat coordinate `json:"lat"`
}

func (location *location) String() string {
	return fmt.Sprintf("{%.6f %.6f}", location.Lon, location.Lat)
}
