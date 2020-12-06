package effects

type Effect interface {
	ApplyEffect([]byte) ([]byte, error)
}
