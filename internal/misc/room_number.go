package misc

type RoomNumber int

const (
	AnyRoomNumber RoomNumber = iota
	OneRoomNumber
	TwoRoomNumber
	ThreeRoomNumber
	ManyRoomNumber
)

var roomNumberViews = map[RoomNumber]string{
	AnyRoomNumber:   "any",
	OneRoomNumber:   "one",
	TwoRoomNumber:   "two",
	ThreeRoomNumber: "three",
	ManyRoomNumber:  "many",
}

func (roomNumber RoomNumber) String() string {
	if view, ok := roomNumberViews[roomNumber]; ok {
		return view
	}
	return "undefined"
}
