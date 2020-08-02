package domria

import (
	"encoding/json"
	"fmt"
)

type item struct {
	BeautifulURL        string     `json:"beautiful_url"`
	MainPhoto           string     `json:"main_photo"`
	UpdatedAt           moment     `json:"updated_at"`
	PriceArr            *prices    `json:"priceArr"`
	TotalSquareMeters   float64    `json:"total_square_meters"`
	LivingSquareMeters  float64    `json:"living_square_meters"`
	KitchenSquareMeters float64    `json:"kitchen_square_meters"`
	RoomsCount          int        `json:"rooms_count"`
	Floor               int        `json:"floor"`
	FloorsCount         int        `json:"floors_count"`
	UserNewbuildNameUK  string     `json:"user_newbuild_name_uk"`
	Longitude           coordinate `json:"longitude"`
	Latitude            coordinate `json:"latitude"`
	StateNameUK         string     `json:"state_name_uk"`
	CityNameUK          string     `json:"city_name_uk"`
	DistrictNameUK      string     `json:"district_name_uk"`
	StreetNameUK        string     `json:"street_name_uk"`
	StreetName          string     `json:"street_name"`
	BuildingNumberStr   string     `json:"building_number_str"`
	Source              string     `json:"-"`
}

func (i *item) UnmarshalJSON(bytes []byte) error {
	raw := item{}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}
	*i = raw
	i.Source = string(bytes)
	return nil
}

func (i *item) String() string {
	return fmt.Sprintf(
		"{%s %s %s %v %.1f %.1f %.1f %d %d %d %s %.6f %.6f %s %s %s %s %s %s %s}",
		i.BeautifulURL,
		i.MainPhoto,
		i.UpdatedAt,
		i.PriceArr,
		i.TotalSquareMeters,
		i.LivingSquareMeters,
		i.KitchenSquareMeters,
		i.RoomsCount,
		i.Floor,
		i.FloorsCount,
		i.UserNewbuildNameUK,
		i.Longitude,
		i.Latitude,
		i.StateNameUK,
		i.CityNameUK,
		i.DistrictNameUK,
		i.StreetNameUK,
		i.StreetName,
		i.BuildingNumberStr,
		i.Source,
	)
}
