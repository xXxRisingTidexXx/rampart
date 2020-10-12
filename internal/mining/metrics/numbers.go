package metrics

type Number int

const (
	StateSanitationNumber Number = iota
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
	SubwaylessSSFNumber
	FailedSSFNumber
	InconclusiveSSFNumber
	SuccessfulSSFNumber
	FailedIZFNumber
	InconclusiveIZFNumber
	SuccessfulIZFNumber
	FailedGZFNumber
	InconclusiveGZFNumber
	SuccessfulGZFNumber
	ApprovedValidationNumber
	UninformativeValidationNumber
	DeniedValidationNumber
	CreatedStoringNumber
	UpdatedStoringNumber
	UnalteredStoringNumber
	FailedStoringNumber
)

var numberViews = map[Number]string{
	StateSanitationNumber:         "state sanitation",
	CitySanitationNumber:          "city sanitation",
	DistrictSanitationNumber:      "district sanitation",
	SwapSanitationNumber:          "swap sanitation",
	StreetSanitationNumber:        "street sanitation",
	HouseNumberSanitationNumber:   "house number sanitation",
	LocatedGeocodingNumber:        "located geocoding",
	UnlocatedGeocodingNumber:      "unlocated geocoding",
	FailedGeocodingNumber:         "failed geocoding",
	InconclusiveGeocodingNumber:   "inconclusive geocoding",
	SuccessfulGeocodingNumber:     "successful geocoding",
	SubwaylessSSFNumber:           "subwayless ssf",
	FailedSSFNumber:               "failed ssf",
	InconclusiveSSFNumber:         "inconclusive ssf",
	SuccessfulSSFNumber:           "successful ssf",
	FailedIZFNumber:               "failed izf",
	InconclusiveIZFNumber:         "inconclusive izf",
	SuccessfulIZFNumber:           "successful izf",
	FailedGZFNumber:               "failed gzf",
	InconclusiveGZFNumber:         "inconclusive gzf",
	SuccessfulGZFNumber:           "successful gzf",
	ApprovedValidationNumber:      "approved validation",
	UninformativeValidationNumber: "uninformative validation",
	DeniedValidationNumber:        "denied validation",
	CreatedStoringNumber:          "created storing",
	UpdatedStoringNumber:          "updated storing",
	UnalteredStoringNumber:        "unaltered storing",
	FailedStoringNumber:           "failed storing",
}

func (number Number) String() string {
	if view, ok := numberViews[number]; ok {
		return view
	}
	return "undefined"
}
