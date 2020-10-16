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
	IsInspected bool
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

func (flat Flat) IsLocated() bool {
	return flat.Point.Lon() != 0 || flat.Point.Lat() != 0
}

func (flat Flat) IsAddressable() bool {
	return flat.State != "" && flat.City != "" && flat.Street != "" && flat.HouseNumber != ""
}

func (flat Flat) MediaCount() int {
	return len(flat.Photos) + len(flat.Panoramas)
}
