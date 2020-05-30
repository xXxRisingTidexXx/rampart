package domria

// TODO: easyjson from mailru
type item struct {
	BeautifulURL        string            `json:"beautiful_url"`
	MainPhoto           string            `json:"main_photo"`
	UpdatedAt           *moment           `json:"updated_at"`
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
	StateNameUK         string            `json:"state_name_uk"`
	CityNameUK          string            `json:"city_name_uk"`
	DistrictNameUK      string            `json:"district_name_uk"`
	DistrictTypeName    string            `json:"district_type_name"`
	StreetNameUK        string            `json:"street_name_uk"`
	StreetName          string            `json:"street_name"`
	BuildingNumberStr   string            `json:"building_number_str"`
}
