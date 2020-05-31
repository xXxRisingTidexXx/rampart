package mining

type Housing string

const (
	Primary Housing = "primary"
	Secondary Housing = "secondary"
)

func (housing Housing) String() string {
	return string(housing)
}
