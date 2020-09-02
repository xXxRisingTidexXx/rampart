package domria

import (
	"fmt"
)

type sourcelessItem struct {
	BeautifulURL        string            `json:"beautiful_url"`
	MainPhoto           string            `json:"main_photo"`
	Photos              map[string]*photo `json:"photos"`
	Panoramas           []*panorama       `json:"panoramas"`
	UpdatedAt           moment            `json:"updated_at"`
	Inspected           int               `json:"inspected"`
	PriceArr            *prices           `json:"priceArr"`
	TotalSquareMeters   float64           `json:"total_square_meters"`
	LivingSquareMeters  float64           `json:"living_square_meters"`
	KitchenSquareMeters float64           `json:"kitchen_square_meters"`
	RoomsCount          int               `json:"rooms_count"`
	Floor               int               `json:"floor"`
	FloorsCount         int               `json:"floors_count"`
	UserNewbuildNameUK  string            `json:"user_newbuild_name_uk"`
	Longitude           coordinate        `json:"longitude"`
	Latitude            coordinate        `json:"latitude"`
	StateNameUK         string            `json:"state_name_uk"`
	CityNameUK          string            `json:"city_name_uk"`
	DistrictNameUK      string            `json:"district_name_uk"`
	StreetNameUK        string            `json:"street_name_uk"`
	StreetName          string            `json:"street_name"`
	BuildingNumberStr   string            `json:"building_number_str"`
}

func (item *sourcelessItem) String() string {
	return fmt.Sprintf(
		"{%s %s %v %v %s %d %v %.1f %.1f %.1f %d %d %d %s %.6f %.6f %s %s %s %s %s %s}",
		item.BeautifulURL,
		item.MainPhoto,
		item.Photos,
		item.Panoramas,
		item.UpdatedAt,
		item.Inspected,
		item.PriceArr,
		item.TotalSquareMeters,
		item.LivingSquareMeters,
		item.KitchenSquareMeters,
		item.RoomsCount,
		item.Floor,
		item.FloorsCount,
		item.UserNewbuildNameUK,
		item.Longitude,
		item.Latitude,
		item.StateNameUK,
		item.CityNameUK,
		item.DistrictNameUK,
		item.StreetNameUK,
		item.StreetName,
		item.BuildingNumberStr,
	)
}
