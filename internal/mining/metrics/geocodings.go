package metrics

type Geocoding int

const (
	LocatedGeocoding Geocoding = iota
	UnlocatedGeocoding
	FailedGeocoding
	InconclusiveGeocoding
	SuccessfulGeocoding
)

var geocodingViews = map[Geocoding]string{
	LocatedGeocoding:      "located",
	UnlocatedGeocoding:    "unlocated",
	FailedGeocoding:       "failed",
	InconclusiveGeocoding: "inconclusive",
	SuccessfulGeocoding:   "successful",
}

func (geocoding Geocoding) String() string {
	if view, ok := geocodingViews[geocoding]; ok {
		return view
	}
	return "undefined"
}
