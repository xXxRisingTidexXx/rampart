package telegram

// TODO: add UUID generation: https://github.com/gofrs/uuid .
type subscription struct {
	ID         int
	City       string
	Price      string
	RoomNumber string
	Floor      string
}
