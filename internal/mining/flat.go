package mining

import (
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Flat struct {
	URL         string
	ImageURLs   []string
	Price       float64
	TotalArea   float64
	LivingArea  float64
	KitchenArea float64
	RoomNumber  int
	Floor       int
	TotalFloor  int
	Housing     misc.Housing
	Point       orb.Point
	City        string
	Street      string
	HouseNumber string
	SSF         float64
	IZF         float64
	GZF         float64
	Miner       string
	ParsingTime time.Time
}

func (f Flat) HasLocation() bool {
	return f.Point.Lon() != 0 || f.Point.Lat() != 0
}
