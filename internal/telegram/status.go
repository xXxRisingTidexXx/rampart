package telegram

import "fmt"

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

var viewStatuses = map[string]Status{
	statusViews[CityStatus]:       CityStatus,
	statusViews[PriceStatus]:      PriceStatus,
	statusViews[RoomNumberStatus]: RoomNumberStatus,
	statusViews[FloorStatus]:      FloorStatus,
}

func ToStatus(view string) (Status, error) {
	if status, ok := viewStatuses[view]; ok {
		return status, nil
	}
	return 0, fmt.Errorf("telegram: status not found, %s", view)
}
