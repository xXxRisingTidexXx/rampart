package metrics

type Validation int

const (
	ApprovedValidation Validation = iota
	UninformativeValidation
	DeniedValidation
)

var validationViews = map[Validation]string{
	ApprovedValidation:      "approved",
	UninformativeValidation: "uninformative",
	DeniedValidation:        "denied",
}

func (validation Validation) String() string {
	if view, ok := validationViews[validation]; ok {
		return view
	}
	return "undefined"
}
