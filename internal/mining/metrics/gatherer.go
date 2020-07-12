package metrics

import (
	"database/sql"
	"fmt"
	"time"
)

func NewGatherer(miner string, db *sql.DB) *Gatherer {
	return &Gatherer{
		miner,
		false,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		db,
	}
}

type Gatherer struct {
	miner                       string
	isSuccessful                bool
	locatedGeocodingNumber      int
	unlocatedGeocodingNumber    int
	failedGeocodingNumber       int
	inconclusiveGeocodingNumber int
	successfulGeocodingNumber   int
	approvedValidationNumber    int
	deniedValidationNumber      int
	createdStoringNumber        int
	updatedStoringNumber        int
	unalteredStoringNumber      int
	failedStoringNumber         int
	fetchingDuration            float64
	geocodingDuration           float64
	readingDuration             float64
	creationDuration            float64
	updateDuration              float64
	totalDuration               float64
	db                          *sql.DB
}

func (gatherer *Gatherer) GatherSuccess() {
	gatherer.isSuccessful = true
}

func (gatherer *Gatherer) GatherLocatedGeocoding() {}

func (gatherer *Gatherer) GatherUnlocatedGeocoding() {}

func (gatherer *Gatherer) GatherFailedGeocoding() {}

func (gatherer *Gatherer) GatherInconclusiveGeocoding() {}

func (gatherer *Gatherer) GatherSuccessfulGeocoding() {}

func (gatherer *Gatherer) GatherApprovedValidation() {}

func (gatherer *Gatherer) GatherDeniedValidation() {}

func (gatherer *Gatherer) GatherCreatedStoring() {}

func (gatherer *Gatherer) GatherUpdatedStoring() {}

func (gatherer *Gatherer) GatherUnalteredStoring() {}

func (gatherer *Gatherer) GatherFailedStoring() {}

func (gatherer *Gatherer) GatherFetchingDuration(start time.Time) {}

func (gatherer *Gatherer) GatherGeocodingDuration(start time.Time) {}

func (gatherer *Gatherer) GatherReadingDuration(start time.Time) {}

func (gatherer *Gatherer) GatherCreationDuration(start time.Time) {}

func (gatherer *Gatherer) GatherUpdateDuration(start time.Time) {}

func (gatherer *Gatherer) GatherTotalDuration(start time.Time) {}

func (gatherer *Gatherer) Flush() error {
	// TODO: write the query.
	_, err := gatherer.db.Exec(``)
	// TODO: nullify fields.
	if err != nil {
		return fmt.Errorf("metrics: gatherer failed to flush, %v", err)
	}
	return nil
}
