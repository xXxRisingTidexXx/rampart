package domria

import (
	"encoding/json"
)

type item struct {
	Source              string
	BeautifulURL        string
	MainPhoto           string
	Photos              map[string]photo
	Panoramas           []panorama
	UpdatedAt           moment
	Inspected           int
	PriceArr            prices
	TotalSquareMeters   float64
	LivingSquareMeters  float64
	KitchenSquareMeters float64
	RoomsCount          int
	Floor               int
	FloorsCount         int
	UserNewbuildNameUK  string
	Longitude           coordinate
	Latitude            coordinate
	StateNameUK         string
	CityNameUK          string
	DistrictNameUK      string
	StreetNameUK        string
	StreetName          string
	BuildingNumberStr   string
}

func (item *item) UnmarshalJSON(bytes []byte) error {
	sourceless := publication{}
	if err := json.Unmarshal(bytes, &sourceless); err != nil {
		return err
	}
	item.Source = string(bytes)
	item.BeautifulURL = sourceless.BeautifulURL
	item.MainPhoto = sourceless.MainPhoto
	item.Photos = sourceless.Photos
	item.Panoramas = sourceless.Panoramas
	item.UpdatedAt = sourceless.UpdatedAt
	item.Inspected = sourceless.Inspected
	item.PriceArr = sourceless.PriceArr
	item.TotalSquareMeters = sourceless.TotalSquareMeters
	item.LivingSquareMeters = sourceless.LivingSquareMeters
	item.KitchenSquareMeters = sourceless.KitchenSquareMeters
	item.RoomsCount = sourceless.RoomsCount
	item.Floor = sourceless.Floor
	item.FloorsCount = sourceless.FloorsCount
	item.UserNewbuildNameUK = sourceless.UserNewbuildNameUK
	item.Longitude = sourceless.Longitude
	item.Latitude = sourceless.Latitude
	item.StateNameUK = sourceless.StateNameUK
	item.CityNameUK = sourceless.CityNameUK
	item.DistrictNameUK = sourceless.DistrictNameUK
	item.StreetNameUK = sourceless.StreetNameUK
	item.StreetName = sourceless.StreetName
	item.BuildingNumberStr = sourceless.BuildingNumberStr
	return nil
}
