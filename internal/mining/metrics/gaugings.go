package metrics

type Gauging int

const (
	SubwaylessSSF Gauging = iota
	FailedSSF
	InconclusiveSSF
	SuccessfulSSF
	FailedIZF
	InconclusiveIZF
	SuccessfulIZF
	FailedGZF
	InconclusiveGZF
	SuccessfulGZF
)

var gaugingViews = map[Gauging]string{
	SubwaylessSSF: "subwayless ssf",
	FailedSSF: "failed ssf",
	InconclusiveSSF: "inconclusive ssf",
	SuccessfulSSF: "successful ssf",
	FailedIZF: "failed izf",
	InconclusiveIZF: "inconclusive izf",
	SuccessfulIZF: "successful izf",
	FailedGZF: "railed gzf",
	InconclusiveGZF: "inconclusive gzf",
	SuccessfulGZF: "successful gzf",
}

func (gauging Gauging) String() string {
	if view, ok := gaugingViews[gauging]; ok {
		return view
	}
	return "undefined"
}
