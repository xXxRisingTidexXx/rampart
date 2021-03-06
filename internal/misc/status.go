package misc

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

func (s Status) String() string {
	if view, ok := statusViews[s]; ok {
		return view
	}
	return "undefined"
}

var viewStatuses = map[string]Status{
	statusViews[CityStatus]:       CityStatus,
	statusViews[PriceStatus]:      PriceStatus,
	statusViews[RoomNumberStatus]: RoomNumberStatus,
	statusViews[FloorStatus]:      FloorStatus,
}

func ToStatus(view string) (Status, bool) {
	status, ok := viewStatuses[view]
	return status, ok
}
