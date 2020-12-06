package imaging

type Effect interface {
	Apply([]byte) ([]byte, error)
	Name() string
}
