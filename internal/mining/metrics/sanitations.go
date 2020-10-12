package metrics

type Sanitation int

const (
	StateSanitation Sanitation = iota
	CitySanitation
	DistrictSanitation
	SwapSanitation
	StreetSanitation
	HouseNumberSanitation
)

var views = map[Sanitation]string{
	StateSanitation:       "state",
	CitySanitation:        "city",
	DistrictSanitation:    "district",
	SwapSanitation:        "swap",
	StreetSanitation:      "street",
	HouseNumberSanitation: "house number",
}

func (sanitation Sanitation) String() string {
	if view, ok := views[sanitation]; ok {
		return view
	}
	return "undefined"
}
