package metrics

type Duration int

const (
	FetchingDuration Duration = iota
	GeocodingDuration
	SSFGaugingDuration
	IZFGaugingDuration
	GZFGaugingDuration
	ReadingDuration
	CreationDuration
	UpdateDuration
	TotalDuration
)

var durationViews = map[Duration]string{
	FetchingDuration:   "fetching",
	GeocodingDuration:  "geocoding",
	SSFGaugingDuration: "ssf gauging",
	IZFGaugingDuration: "izf gauging",
	GZFGaugingDuration: "gzf gauging",
	ReadingDuration:    "reading",
	CreationDuration:   "creation",
	UpdateDuration:     "update",
	TotalDuration:      "total",
}

func (duration Duration) String() string {
	if view, ok := durationViews[duration]; ok {
		return view
	}
	return "undefined"
}
