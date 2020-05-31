package domria

import (
	"fmt"
	"rampart/pkg/mining"
	"time"
)

// TODO: check golang gis geometry library integration.
type flat struct {
	originURL   string
	imageURL    string
	updatedAt   *time.Time
	price       float64
	totalArea   float64
	livingArea  float64
	kitchenArea float64
	roomNumber  int
	floor       int
	totalFloor  int
	housing     mining.Housing
	complex     string
	longitude   float64
	latitude    float64
	state       string
	city        string
	district    string
	street      string
	houseNumber string
}

func (flat *flat) String() string {
	return fmt.Sprintf(
		"{%s %s %v %.2f %.1f %.1f %.1f %d %d %d %s %s %.6f %.6f %s %s %s %s %s}",
		flat.originURL,
		flat.imageURL,
		flat.updatedAt,
		flat.price,
		flat.totalArea,
		flat.livingArea,
		flat.kitchenArea,
		flat.roomNumber,
		flat.floor,
		flat.totalFloor,
		flat.housing,
		flat.complex,
		flat.longitude,
		flat.latitude,
		flat.state,
		flat.city,
		flat.district,
		flat.street,
		flat.houseNumber,
	)
}
