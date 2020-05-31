package mining

type Housing int

const (
	Primary Housing = iota
	Secondary
)

func (housing Housing) String() string {
	switch housing {
	case Primary:
		return "primary"
	case Secondary:
		return "secondary"
	default:
		return "undefined"
	}
}
