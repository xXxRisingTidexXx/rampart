package domria

import (
	"rampart/pkg/mining"
	"time"
)

type flat struct {
	originURL   string
	imageURL    string
	accessTime  time.Time
	price       float64
	totalArea   float64
	livingArea  float64
	kitchenArea float64
	roomNumber  int
	floor       int
	totalFloor  int
	housing     mining.Housing
	complex     string
	longitude   float64  // TODO: check golang gis geometry library integration
	latitude    float64
	state       string
	city        string
	district    string
	street      string
	houseNumber string
}
