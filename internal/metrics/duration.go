package metrics

type Duration int

const (
	FetchingDuration Duration = iota
	GeocodingDuration
	SSFGaugingDuration
	IZFGaugingDuration
	GZFGaugingDuration
	ReadingFlatStoringDuration
	CreationFlatStoringDuration
	UpdateFlatStoringDuration
	ReadingImageStoringDuration
	CreationImageStoringDuration
	TotalDuration
)

var durationViews = map[Duration]string{
	FetchingDuration:             "fetching",
	GeocodingDuration:            "geocoding",
	SSFGaugingDuration:           "ssf-gauging",
	IZFGaugingDuration:           "izf-gauging",
	GZFGaugingDuration:           "gzf-gauging",
	ReadingFlatStoringDuration:   "reading-flat-storing",
	CreationFlatStoringDuration:  "creation-flat-storing",
	UpdateFlatStoringDuration:    "update-flat-storing",
	ReadingImageStoringDuration:  "reading-image-storing",
	CreationImageStoringDuration: "creation-image-storing",
	TotalDuration:                "total",
}

func (duration Duration) String() string {
	if view, ok := durationViews[duration]; ok {
		return view
	}
	return "undefined"
}
