package metrics

import (
	"rampart/internal/config"
)

func NewGatherer(miner string, config *config.Gatherer) *Gatherer {
	return &Gatherer{
		newCounterTracker(miner, config.FailedFetchingTracker),
		newHistogramTracker(miner, config.FetchingDurationTracker),
		newCounterTracker(miner, config.FetchedFlatsTracker),
		newCounterTracker(miner, config.LocatedFlatsTracker),
		newCounterTracker(miner, config.UnlocatedFlatsTracker),
		newCounterTracker(miner, config.FailedGeocodingTracker),
		newCounterTracker(miner, config.EmptyGeocodingTracker),
		newCounterTracker(miner, config.SuccessfulGeocodingTracker),
		newHistogramTracker(miner, config.GeocodingDurationTracker),
		newCounterTracker(miner, config.ValidatedFlatsTracker),
		newCounterTracker(miner, config.InvalidatedFlatsTracker),
		newCounterTracker(miner, config.CreatedFlatsTracker),
		newCounterTracker(miner, config.UpdatedFlatsTracker),
		newCounterTracker(miner, config.UnalteredFlatsTracker),
		newCounterTracker(miner, config.FailedStoringTracker),
		newHistogramTracker(miner, config.ReadingDurationTracker),
		newHistogramTracker(miner, config.CreationDurationTracker),
		newHistogramTracker(miner, config.UpdateDurationTracker),
		newHistogramTracker(miner, config.RunDurationTracker),
	}
}

type Gatherer struct {
	failedFetchingTracker      *counterTracker
	fetchingDurationTracker    *histogramTracker
	fetchedFlatsTracker        *counterTracker
	locatedFlatsTracker        *counterTracker
	unlocatedFlatsTracker      *counterTracker
	failedGeocodingTracker     *counterTracker
	emptyGeocodingTracker      *counterTracker
	successfulGeocodingTracker *counterTracker
	geocodingDurationTracker   *histogramTracker
	validatedFlatsTracker      *counterTracker
	invalidatedFlatsTracker    *counterTracker
	createdFlatsTracker        *counterTracker
	updatedFlatsTracker        *counterTracker
	unalteredFlatsTracker      *counterTracker
	failedStoringTracker       *counterTracker
	readingDurationTracker     *histogramTracker
	creationDurationTracker    *histogramTracker
	updateDurationTracker      *histogramTracker
	runDurationTracker         *histogramTracker
}

func (gatherer *Gatherer) GatherFailedFetching() {
	gatherer.failedFetchingTracker.track(1)
}

func (gatherer *Gatherer) GatherFetchingDuration(fetchingDuration float64) {
	gatherer.fetchingDurationTracker.track(fetchingDuration)
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

func (gatherer *Gatherer) GatherGeocodingDuration(geocodingDuration float64) {
	gatherer.geocodingDurationTracker.track(geocodingDuration)
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

func (gatherer *Gatherer) GatherReadingDuration(readingDuration float64) {
	gatherer.readingDurationTracker.track(readingDuration)
}

func (gatherer *Gatherer) GatherCreationDuration(creationDuration float64) {
	gatherer.creationDurationTracker.track(creationDuration)
}

func (gatherer *Gatherer) GatherUpdateDuration(updateDuration float64) {
	gatherer.updateDurationTracker.track(updateDuration)
}

func (gatherer *Gatherer) GatherRunDuration(runDuration float64) {
	gatherer.runDurationTracker.track(runDuration)
}

func (gatherer *Gatherer) Unregister() error {
	return gatherer.runDurationTracker.unregister()
}
