package misc

type Interior int

const (
	LuxuryInterior Interior = iota
	ComfortInterior
	JunkInterior
	ConstructionInterior
	ExcessInterior
)

var interiorViews = map[Interior]string{
	LuxuryInterior:       "luxury",
	ComfortInterior:      "comfort",
	JunkInterior:         "junk",
	ConstructionInterior: "construction",
	ExcessInterior:       "excess",
}

func (i Interior) String() string {
	if view, ok := interiorViews[i]; ok {
		return view
	}
	return "undefined"
}
