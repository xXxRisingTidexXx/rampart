package domria

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"rampart/internal/mining"
	"time"
)

type flat struct {
	originURL   string
	imageURL    string
	updateTime  *time.Time
	price       float64
	totalArea   float64
	livingArea  float64
	kitchenArea float64
	roomNumber  int
	floor       int
	totalFloor  int
	housing     mining.Housing
	complex     string
	point       *geom.Point
	state       string
	city        string
	district    string
	street      string
	houseNumber string
}

func (flat *flat) String() string {
	return fmt.Sprintf(
		"{%s %s %v %.2f %.1f %.1f %.1f %d %d %d %s %s %v %s %s %s %s %s}",
		flat.originURL,
		flat.imageURL,
		flat.updateTime,
		flat.price,
		flat.totalArea,
		flat.livingArea,
		flat.kitchenArea,
		flat.roomNumber,
		flat.floor,
		flat.totalFloor,
		flat.housing,
		flat.complex,
		flat.point,
		flat.state,
		flat.city,
		flat.district,
		flat.street,
		flat.houseNumber,
	)
}
