package mining

type Amplifier interface {
	AmplifyFlat(Flat) (Flat, error)
}
