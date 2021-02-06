package domria

import (
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Flat struct {
	Source      string
	URL         string
	Photos      []string
	Panoramas   []string
	UpdateTime  time.Time
	IsSold      bool
	Price       float64
	TotalArea   float64
	LivingArea  float64
	KitchenArea float64
	RoomNumber  int
	Floor       int
	TotalFloor  int
	Housing     misc.Housing
	Complex     string
	Point       orb.Point
	State       string
	City        string
	District    string
	Street      string
	HouseNumber string
	SSF         float64
	IZF         float64
	GZF         float64
}

func (f Flat) IsLocated() bool {
	return f.Point.Lon() != 0 || f.Point.Lat() != 0
}

func (f Flat) IsAddressable() bool {
	return f.State != "" && f.City != "" && f.Street != "" && f.HouseNumber != ""
}

func (f Flat) ImageCount() int {
	return len(f.Photos) + len(f.Panoramas)
}
