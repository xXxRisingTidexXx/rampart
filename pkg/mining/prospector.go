package mining

type Prospector interface {
	Prospect(housing Housing) error
}
