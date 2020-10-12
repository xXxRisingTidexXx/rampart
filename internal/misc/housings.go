package misc

const (
	HousingPrimary   = "primary"
	HousingSecondary = "secondary"
)

type Housing int

const (
	PrimaryHousing Housing = iota
	SecondaryHousing
)

var housingViews = map[Housing]string{
	PrimaryHousing: "primary",
	SecondaryHousing: "secondary",
}

func (housing Housing) String() string {
	if view, ok := housingViews[housing]; ok {
		return view
	}
	return "undefined"
}
