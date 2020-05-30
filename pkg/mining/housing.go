package mining

type Housing int

const (
	Primary Housing = iota
	Secondary
	All
)

func (housing Housing) String() string {
	switch housing {
	case Primary:
		return "primary"
	case Secondary:
		return "secondary"
	case All:
		return "all"
	default:
		return "undefined"
	}
}
