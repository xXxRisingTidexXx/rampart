package domria

import (
	"fmt"
)

type position struct {
	Lon coordinate `json:"lon"`
	Lat coordinate `json:"lat"`
}

func (position *position) String() string {
	return fmt.Sprintf("{%.6f %.6f}", position.Lon, position.Lat)
}
