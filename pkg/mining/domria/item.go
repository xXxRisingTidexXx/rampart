package domria

type item struct {
	BeautifulURL        string            `json:"beautiful_url"`
	MainPhoto           string            `json:"main_photo"`
	PriceArr            map[string]string `json:"priceArr"`
	TotalSquareMeters   float64           `json:"total_square_meters"`
	LivingSquareMeters  float64           `json:"living_square_meters"`
	KitchenSquareMeters float64           `json:"kitchen_square_meters"`
	RoomsCount          int               `json:"rooms_count"`
	Floor               int               `json:"floor"`
	FloorsCount         int               `json:"floors_count"`
	UserNewbuildNameUK  string            `json:"user_newbuild_name_uk"`
	Longitude           coordinate        `json:"longitude"`
	Latitude            coordinate        `json:"latitude"`
	DistrictNameUK      string            `json:"district_name_uk"`
	StreetNameUK        string            `json:"street_name_uk"`
	StreetName          string            `json:"street_name"`
	BuildingNumberStr   string            `json:"building_number_str"`
}
