package mining

type Miner interface {
	Name() string
	Spec() string
	MineFlat() (Flat, error)
}
