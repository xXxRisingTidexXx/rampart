package misc

type Floor int

const (
	AnyFloor Floor = iota
	LowFloor
	HighFloor
)

var floorViews = map[Floor]string{AnyFloor:  "any", LowFloor:  "low", HighFloor: "high"}

func (floor Floor) String() string {
	if view, ok := floorViews[floor]; ok {
		return view
	}
	return "undefined"
}
