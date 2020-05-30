package mining

type Prospector interface {
	Prospect(state, city string, housing Housing)
}
