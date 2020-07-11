package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"rampart/internal/config"
	"rampart/internal/misc"
	"time"
)

// TODO: pass the registry to metrics server.
func NewGatherer(miner string, config *config.Gatherer) *Gatherer {
	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())
	return &Gatherer{
		registry,
		newCounterTracker(registry, config.RunTracker, miner, config.Statuses.Targets()),
		config.Statuses,
		newCounterTracker(registry, config.GeocodingTracker, miner, config.Categories.Targets()),
		config.Categories,
		newCounterTracker(registry, config.ValidationTracker, miner, config.Verdicts.Targets()),
		config.Verdicts,
		newCounterTracker(registry, config.StoringTracker, miner, config.Consequences.Targets()),
		config.Consequences,
		newHistogramTracker(registry, config.DurationTracker, miner, config.Processes.Targets()),
		config.Processes,
	}
}

type Gatherer struct {
	registry          *prometheus.Registry
	runTracker        *counterTracker
	statuses          *misc.Statuses
	geocodingTracker  *counterTracker
	categories        *misc.Categories
	validationTracker *counterTracker
	verdicts          *misc.Verdicts
	storingTracker    *counterTracker
	consequences      *misc.Consequences
	durationTracker   *histogramTracker
	processes         *misc.Processes
}

func (gatherer *Gatherer) GatherFailureRun() {
	gatherer.runTracker.track(gatherer.statuses.Failure)
}

func (gatherer *Gatherer) GatherSuccessRun() {
	gatherer.runTracker.track(gatherer.statuses.Success)
}

func (gatherer *Gatherer) GatherLocatedGeocoding() {
	gatherer.geocodingTracker.track(gatherer.categories.Located)
}

func (gatherer *Gatherer) GatherUnlocatedGeocoding() {
	gatherer.geocodingTracker.track(gatherer.categories.Unlocated)
}

func (gatherer *Gatherer) GatherFailedGeocoding() {
	gatherer.geocodingTracker.track(gatherer.categories.Failed)
}

func (gatherer *Gatherer) GatherInconclusiveGeocoding() {
	gatherer.geocodingTracker.track(gatherer.categories.Inconclusive)
}

func (gatherer *Gatherer) GatherSuccessfulGeocoding() {
	gatherer.geocodingTracker.track(gatherer.categories.Successful)
}

func (gatherer *Gatherer) GatherApprovedValidation() {
	gatherer.validationTracker.track(gatherer.verdicts.Approved)
}

func (gatherer *Gatherer) GatherDeniedValidation() {
	gatherer.validationTracker.track(gatherer.verdicts.Denied)
}

func (gatherer *Gatherer) GatherCreatedStoring() {
	gatherer.storingTracker.track(gatherer.consequences.Created)
}

func (gatherer *Gatherer) GatherUpdatedStoring() {
	gatherer.storingTracker.track(gatherer.consequences.Updated)
}

func (gatherer *Gatherer) GatherUnalteredStoring() {
	gatherer.storingTracker.track(gatherer.consequences.Unaltered)
}

func (gatherer *Gatherer) GatherFailedStoring() {
	gatherer.storingTracker.track(gatherer.consequences.Failed)
}

func (gatherer *Gatherer) GatherFetchingDuration(start time.Time) {
	gatherer.durationTracker.track(gatherer.processes.Fetching, time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherGeocodingDuration(start time.Time) {
	gatherer.durationTracker.track(gatherer.processes.Geocoding, time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherReadingDuration(start time.Time) {
	gatherer.durationTracker.track(gatherer.processes.Reading, time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherCreationDuration(start time.Time) {
	gatherer.durationTracker.track(gatherer.processes.Creation, time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherUpdateDuration(start time.Time) {
	gatherer.durationTracker.track(gatherer.processes.Update, time.Since(start).Seconds())
}

func (gatherer *Gatherer) GatherRunDuration(start time.Time) {
	gatherer.durationTracker.track(gatherer.processes.Run, time.Since(start).Seconds())
}

func (gatherer *Gatherer) Unregister() error {
	if err := gatherer.runTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.geocodingTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.validationTracker.unregister(); err != nil {
		return err
	}
	if err := gatherer.storingTracker.unregister(); err != nil {
		return err
	}
	return gatherer.durationTracker.unregister()
}
