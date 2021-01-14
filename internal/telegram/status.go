package telegram

type Status int

const (
	CityStatus Status = iota
	PriceStatus
	RoomNumberStatus
	FloorStatus
)

var statusViews = map[Status]string{
	CityStatus:       "city",
	PriceStatus:      "price",
	RoomNumberStatus: "room-number",
	FloorStatus:      "floor",
}

func (status Status) String() string {
	if view, ok := statusViews[status]; ok {
		return view
	}
	return "undefined"
}
