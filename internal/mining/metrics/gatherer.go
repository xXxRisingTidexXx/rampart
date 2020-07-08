package metrics

import (
	"rampart/internal/config"
	"time"
)

func NewGatherer(miner string, config *config.Gatherer) *Gatherer {
	return &Gatherer{

	}
}

type Gatherer struct {

}

func (gatherer *Gatherer) GatherFailedFetching() {
	gatherer.failedFetchingTracker.track(1)
}

func (gatherer *Gatherer) GatherFetchingDuration(start time.Time) {
	gatherer.fetchingDurationTracker.track(time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherFetchedFlats(fetchedFlats int) {
	gatherer.fetchedFlatsTracker.track(float64(fetchedFlats))
}

func (gatherer *Gatherer) GatherLocatedFlats() {
	gatherer.locatedFlatsTracker.track(1)
}

func (gatherer *Gatherer) GatherUnlocatedFlats() {
	gatherer.unlocatedFlatsTracker.track(1)
}

func (gatherer *Gatherer) GatherFailedGeocoding() {
	gatherer.failedGeocodingTracker.track(1)
}

func (gatherer *Gatherer) GatherEmptyGeocoding() {
	gatherer.emptyGeocodingTracker.track(1)
}

func (gatherer *Gatherer) GatherSuccessfulGeocoding() {
	gatherer.successfulGeocodingTracker.track(1)
}

func (gatherer *Gatherer) GatherGeocodingDuration(start time.Time) {
	gatherer.geocodingDurationTracker.track(time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherValidatedFlats() {
	gatherer.validatedFlatsTracker.track(1)
}

func (gatherer *Gatherer) GatherInvalidatedFlats() {
	gatherer.invalidatedFlatsTracker.track(1)
}

func (gatherer *Gatherer) GatherCreatedFlats() {
	gatherer.createdFlatsTracker.track(1)
}

func (gatherer *Gatherer) GatherUpdatedFlats() {
	gatherer.updatedFlatsTracker.track(1)
}

func (gatherer *Gatherer) GatherUnalteredFlats() {
	gatherer.unalteredFlatsTracker.track(1)
}

func (gatherer *Gatherer) GatherFailedStoring() {
	gatherer.failedStoringTracker.track(1)
}

func (gatherer *Gatherer) GatherReadingDuration(start time.Time) {
	gatherer.readingDurationTracker.track(time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherCreationDuration(start time.Time) {
	gatherer.creationDurationTracker.track(time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherUpdateDuration(start time.Time) {
	gatherer.updateDurationTracker.track(time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherRunDuration(start time.Time) {
	gatherer.runDurationTracker.track(time.Since(start).Seconds())
}

func (gatherer *Gatherer) Unregister() error {
	if err := gatherer.failedFetchingTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.fetchingDurationTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.fetchedFlatsTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.locatedFlatsTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.unlocatedFlatsTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.failedGeocodingTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.emptyGeocodingTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.successfulGeocodingTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.geocodingDurationTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.validatedFlatsTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.invalidatedFlatsTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.createdFlatsTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.updatedFlatsTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.unalteredFlatsTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.failedStoringTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.readingDurationTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.creationDurationTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.updateDurationTracker.unregister(); err != nil {
		return err
	}
	return gatherer.runDurationTracker.unregister()
}
