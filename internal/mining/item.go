package mining

type item struct {
	BeautifulURL        string           `json:"beautiful_url"`
	Photos              map[string]photo `json:"photos"`
	SaleDate            string           `json:"sale_date"`
	PriceArr            prices           `json:"priceArr"`
	TotalSquareMeters   float64          `json:"total_square_meters"`
	LivingSquareMeters  float64          `json:"living_square_meters"`
	KitchenSquareMeters float64          `json:"kitchen_square_meters"`
	RoomsCount          int              `json:"rooms_count"`
	Floor               int              `json:"floor"`
	FloorsCount         int              `json:"floors_count"`
	RealtySaleType      int              `json:"realty_sale_type"`
	Longitude           coordinate       `json:"longitude"`
	Latitude            coordinate       `json:"latitude"`
	CityNameUK          string           `json:"city_name_uk"`
	DistrictNameUK      string           `json:"district_name_uk"`
	StreetNameUK        string           `json:"street_name_uk"`
	StreetName          string           `json:"street_name"`
	BuildingNumberStr   number           `json:"building_number_str"`
}
