package domria

import (
	"encoding/json"
)

type item struct {
	Source              string
	BeautifulURL        string
	Photos              map[string]photo
	Panoramas           []panorama
	UpdatedAt           moment
	SaleDate            string
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

func (i *item) UnmarshalJSON(bytes []byte) error {
	p := publication{}
	if err := json.Unmarshal(bytes, &p); err != nil {
		return err
	}
	i.Source = string(bytes)
	i.BeautifulURL = p.BeautifulURL
	i.Photos = p.Photos
	i.Panoramas = p.Panoramas
	i.UpdatedAt = p.UpdatedAt
	i.SaleDate = p.SaleDate
	i.Inspected = p.Inspected
	i.PriceArr = p.PriceArr
	i.TotalSquareMeters = p.TotalSquareMeters
	i.LivingSquareMeters = p.LivingSquareMeters
	i.KitchenSquareMeters = p.KitchenSquareMeters
	i.RoomsCount = p.RoomsCount
	i.Floor = p.Floor
	i.FloorsCount = p.FloorsCount
	i.UserNewbuildNameUK = p.UserNewbuildNameUK
	i.Longitude = p.Longitude
	i.Latitude = p.Latitude
	i.StateNameUK = p.StateNameUK
	i.CityNameUK = p.CityNameUK
	i.DistrictNameUK = p.DistrictNameUK
	i.StreetNameUK = p.StreetNameUK
	i.StreetName = p.StreetName
	i.BuildingNumberStr = p.BuildingNumberStr
	return nil
}
