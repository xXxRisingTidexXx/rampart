package config

import (
	"fmt"
)

type Gatherer struct {
	FailedFetchingTracker      *CounterTracker   `yaml:"failedFetchingTracker"`
	FetchingDurationTracker    *HistogramTracker `yaml:"fetchingDurationTracker"`
	FetchedFlatsTracker        *CounterTracker   `yaml:"fetchedFlatsTracker"`
	LocatedFlatsTracker        *CounterTracker   `yaml:"locatedFlatsTracker"`
	UnlocatedFlatsTracker      *CounterTracker   `yaml:"unlocatedFlatsTracker"`
	FailedGeocodingTracker     *CounterTracker   `yaml:"failedGeocodingTracker"`
	EmptyGeocodingTracker      *CounterTracker   `yaml:"emptyGeocodingTracker"`
	SuccessfulGeocodingTracker *CounterTracker   `yaml:"successfulGeocodingTracker"`
	GeocodingDurationTracker   *HistogramTracker `yaml:"geocodingDurationTracker"`
	ValidatedFlatsTracker      *CounterTracker   `yaml:"validatedFlatsTracker"`
	InvalidatedFlatsTracker    *CounterTracker   `yaml:"invalidatedFlatsTracker"`
	CreatedFlatsTracker        *CounterTracker   `yaml:"createdFlatsTracker"`
	UpdatedFlatsTracker        *CounterTracker   `yaml:"updatedFlatsTracker"`
	UnalteredFlatsTracker      *CounterTracker   `yaml:"unalteredFlatsTracker"`
	FailedStoringTracker       *CounterTracker   `yaml:"failedStoringTracker"`
	ReadingDurationTracker     *HistogramTracker `yaml:"readingDurationTracker"`
	CreationDurationTracker    *HistogramTracker `yaml:"creationDurationTracker"`
	UpdateDurationTracker      *HistogramTracker `yaml:"updateDurationTracker"`
	RunDurationTracker         *HistogramTracker `yaml:"runDurationTracker"`
}

func (gatherer *Gatherer) String() string {
	return fmt.Sprintf(
		"{%v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v}",
		gatherer.FailedFetchingTracker,
		gatherer.FetchingDurationTracker,
		gatherer.FetchedFlatsTracker,
		gatherer.LocatedFlatsTracker,
		gatherer.UnlocatedFlatsTracker,
		gatherer.FailedGeocodingTracker,
		gatherer.EmptyGeocodingTracker,
		gatherer.SuccessfulGeocodingTracker,
		gatherer.GeocodingDurationTracker,
		gatherer.ValidatedFlatsTracker,
		gatherer.InvalidatedFlatsTracker,
		gatherer.CreatedFlatsTracker,
		gatherer.UpdatedFlatsTracker,
		gatherer.UnalteredFlatsTracker,
		gatherer.FailedStoringTracker,
		gatherer.ReadingDurationTracker,
		gatherer.CreationDurationTracker,
		gatherer.UpdateDurationTracker,
		gatherer.RunDurationTracker,
	)
}
