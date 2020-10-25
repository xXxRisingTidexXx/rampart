package domria

type Kind int

const (
	PhotoKind Kind = iota
	PanoramaKind
)

var kindViews = map[Kind]string{PhotoKind: "photo", PanoramaKind: "panorama"}

func (kind Kind) String() string {
	if view, ok := kindViews[kind]; ok {
		return view
	}
	return "undefined"
}
