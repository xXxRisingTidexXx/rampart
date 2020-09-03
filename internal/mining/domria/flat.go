package domria

import (
	"fmt"
	"github.com/paulmach/orb"
	"time"
)

type Flat struct {
	OriginURL   string
	ImageURL    string
	MediaCount  int
	UpdateTime  time.Time
	IsInspected bool
	Price       float64
	TotalArea   float64
	LivingArea  float64
	KitchenArea float64
	RoomNumber  int
	Floor       int
	TotalFloor  int
	Housing     string
	Complex     string
	Point       orb.Point
	State       string
	City        string
	District    string
	Street      string
	HouseNumber string
	Source      string
}

func (flat *Flat) IsLocated() bool {
	return flat.Point.Lon() != 0 || flat.Point.Lat() != 0
}

func (flat *Flat) IsAddressable() bool {
	return flat.State != "" && flat.City != "" && flat.Street != "" && flat.HouseNumber != ""
}

func (flat *Flat) String() string {
	return fmt.Sprintf(
		"{%s %s %d %s %t %.2f %.1f %.1f %.1f %d %d %d %s %s %v %s %s %s %s %s %s}",
		flat.OriginURL,
		flat.ImageURL,
		flat.MediaCount,
		flat.UpdateTime,
		flat.IsInspected,
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
		flat.State,
		flat.City,
		flat.District,
		flat.Street,
		flat.HouseNumber,
		flat.Source,
	)
}
