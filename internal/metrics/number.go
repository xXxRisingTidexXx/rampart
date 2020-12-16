package metrics

type Number int

const (
	FailedFetchingNumber Number = iota
	StateSanitationNumber
	CitySanitationNumber
	DistrictSanitationNumber
	SwapSanitationNumber
	StreetSanitationNumber
	HouseNumberSanitationNumber
	LocatedGeocodingNumber
	UnlocatedGeocodingNumber
	FailedGeocodingNumber
	InconclusiveGeocodingNumber
	SuccessfulGeocodingNumber
	SubwaylessSSFGaugingNumber
	FailedSSFGaugingNumber
	InconclusiveSSFGaugingNumber
	SuccessfulSSFGaugingNumber
	FailedIZFGaugingNumber
	InconclusiveIZFGaugingNumber
	SuccessfulIZFGaugingNumber
	FailedGZFGaugingNumber
	InconclusiveGZFGaugingNumber
	SuccessfulGZFGaugingNumber
	ApprovedValidationNumber
	UninformativeValidationNumber
	SoldValidationNumber
	DeniedValidationNumber
	CreatedFlatStoringNumber
	UpdatedFlatStoringNumber
	UnalteredFlatStoringNumber
	FailedFlatStoringNumber
	CreatedImageStoringNumber
	UnalteredImageStoringNumber
	FailedImageStoringNumber
)

var numberViews = map[Number]string{
	FailedFetchingNumber:          "failed-fetching",
	StateSanitationNumber:         "state-sanitation",
	CitySanitationNumber:          "city-sanitation",
	DistrictSanitationNumber:      "district-sanitation",
	SwapSanitationNumber:          "swap-sanitation",
	StreetSanitationNumber:        "street-sanitation",
	HouseNumberSanitationNumber:   "house-number-sanitation",
	LocatedGeocodingNumber:        "located-geocoding",
	UnlocatedGeocodingNumber:      "unlocated-geocoding",
	FailedGeocodingNumber:         "failed-geocoding",
	InconclusiveGeocodingNumber:   "inconclusive-geocoding",
	SuccessfulGeocodingNumber:     "successful-geocoding",
	SubwaylessSSFGaugingNumber:    "subwayless-ssf-gauging",
	FailedSSFGaugingNumber:        "failed-ssf-gauging",
	InconclusiveSSFGaugingNumber:  "inconclusive-ssf-gauging",
	SuccessfulSSFGaugingNumber:    "successful-ssf-gauging",
	FailedIZFGaugingNumber:        "failed-izf-gauging",
	InconclusiveIZFGaugingNumber:  "inconclusive-izf-gauging",
	SuccessfulIZFGaugingNumber:    "successful-izf-gauging",
	FailedGZFGaugingNumber:        "failed-gzf-gauging",
	InconclusiveGZFGaugingNumber:  "inconclusive-gzf-gauging",
	SuccessfulGZFGaugingNumber:    "successful-gzf-gauging",
	ApprovedValidationNumber:      "approved-validation",
	UninformativeValidationNumber: "uninformative-validation",
	SoldValidationNumber:          "sold-validation",
	DeniedValidationNumber:        "denied-validation",
	CreatedFlatStoringNumber:      "created-flat-storing",
	UpdatedFlatStoringNumber:      "updated-flat-storing",
	UnalteredFlatStoringNumber:    "unaltered-flat-storing",
	FailedFlatStoringNumber:       "failed-flat-storing",
	CreatedImageStoringNumber:     "created-image-storing",
	UnalteredImageStoringNumber:   "unaltered-image-storing",
	FailedImageStoringNumber:      "failed-image-storing",
}

func (number Number) String() string {
	if view, ok := numberViews[number]; ok {
		return view
	}
	return "undefined"
}
