package domria

import (
	"fmt"
	"github.com/paulmach/orb"
	"time"
)

type Flat struct {
	OriginURL             string
	ImageURL              string
	UpdateTime            time.Time
	Price                 float64
	TotalArea             float64
	LivingArea            float64
	KitchenArea           float64
	RoomNumber            int
	Floor                 int
	TotalFloor            int
	Housing               string
	Complex               string
	Point                 orb.Point
	SubwayStationDistance float64
	State                 string
	City                  string
	District              string
	Street                string
	HouseNumber           string
	Source                string
}

func (flat *Flat) String() string {
	return fmt.Sprintf(
		"{%s %s %s %.2f %.1f %.1f %.1f %d %d %d %s %s %v %.5f %s %s %s %s %s %s}",
		flat.OriginURL,
		flat.ImageURL,
		flat.UpdateTime,
		flat.Price,
		flat.TotalArea,
		flat.LivingArea,
		flat.KitchenArea,
		flat.RoomNumber,
		flat.Floor,
		flat.TotalFloor,
		flat.Housing,
		flat.Complex,
		flat.Point,
		flat.SubwayStationDistance,
		flat.State,
		flat.City,
		flat.District,
		flat.Street,
		flat.HouseNumber,
		flat.Source,
	)
}
