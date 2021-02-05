package domria

type Kind int

const (
	PhotoKind Kind = iota
	PanoramaKind
)

var kindViews = map[Kind]string{PhotoKind: "photo", PanoramaKind: "panorama"}

func (k Kind) String() string {
	if view, ok := kindViews[k]; ok {
		return view
	}
	return "undefined"
}
