package metrics

import (
	"database/sql"
	"fmt"
	"time"
)

func NewGatherer(miner string, db *sql.DB) *Gatherer {
	return &Gatherer{miner: miner, db: db}
}

type Gatherer struct {
	miner                         string
	stateSanitizationNumber       int
	citySanitizationNumber        int
	districtSanitizationNumber    int
	swapSanitizationNumber        int
	streetSanitizationNumber      int
	houseNumberSanitizationNumber int
	locatedGeocodingNumber        int
	unlocatedGeocodingNumber      int
	failedGeocodingNumber         int
	inconclusiveGeocodingNumber   int
	successfulGeocodingNumber     int
	approvedValidationNumber      int
	deniedValidationNumber        int
	createdStoringNumber          int
	updatedStoringNumber          int
	unalteredStoringNumber        int
	failedStoringNumber           int
	fetchingDuration              float64
	geocodingDurationSum          float64
	geocodingDurationCount        float64
	readingDurationSum            float64
	readingDurationCount          float64
	creationDurationSum           float64
	creationDurationCount         float64
	updateDurationSum             float64
	updateDurationCount           float64
	totalDuration                 float64
	db                            *sql.DB
}

func (gatherer *Gatherer) GatherStateSanitization() {
	gatherer.stateSanitizationNumber++
}

func (gatherer *Gatherer) GatherCitySanitization() {
	gatherer.citySanitizationNumber++
}

func (gatherer *Gatherer) GatherDistrictSanitization() {
	gatherer.districtSanitizationNumber++
}

func (gatherer *Gatherer) GatherSwapSanitization() {
	gatherer.swapSanitizationNumber++
}

func (gatherer *Gatherer) GatherStreetSanitization() {
	gatherer.streetSanitizationNumber++
}

func (gatherer *Gatherer) GatherHouseNumberSanitization() {
	gatherer.houseNumberSanitizationNumber++
}

func (gatherer *Gatherer) GatherLocatedGeocoding() {
	gatherer.locatedGeocodingNumber++
}

func (gatherer *Gatherer) GatherUnlocatedGeocoding() {
	gatherer.unlocatedGeocodingNumber++
}

func (gatherer *Gatherer) GatherFailedGeocoding() {
	gatherer.failedGeocodingNumber++
}

func (gatherer *Gatherer) GatherInconclusiveGeocoding() {
	gatherer.inconclusiveGeocodingNumber++
}

func (gatherer *Gatherer) GatherSuccessfulGeocoding() {
	gatherer.successfulGeocodingNumber++
}

func (gatherer *Gatherer) GatherApprovedValidation() {
	gatherer.approvedValidationNumber++
}

func (gatherer *Gatherer) GatherDeniedValidation() {
	gatherer.deniedValidationNumber++
}

func (gatherer *Gatherer) GatherCreatedStoring() {
	gatherer.createdStoringNumber++
}

func (gatherer *Gatherer) GatherUpdatedStoring() {
	gatherer.updatedStoringNumber++
}

func (gatherer *Gatherer) GatherUnalteredStoring() {
	gatherer.unalteredStoringNumber++
}

func (gatherer *Gatherer) GatherFailedStoring() {
	gatherer.failedStoringNumber++
}

func (gatherer *Gatherer) GatherFetchingDuration(start time.Time) {
	gatherer.fetchingDuration = time.Since(start).Seconds()
}

func (gatherer *Gatherer) GatherGeocodingDuration(start time.Time) {
	gatherer.geocodingDurationSum += time.Since(start).Seconds()
	gatherer.geocodingDurationCount++
}

func (gatherer *Gatherer) GatherReadingDuration(start time.Time) {
	gatherer.readingDurationSum += time.Since(start).Seconds()
	gatherer.readingDurationCount++
}

func (gatherer *Gatherer) GatherCreationDuration(start time.Time) {
	gatherer.creationDurationSum += time.Since(start).Seconds()
	gatherer.creationDurationCount++
}

func (gatherer *Gatherer) GatherUpdateDuration(start time.Time) {
	gatherer.updateDurationSum += time.Since(start).Seconds()
	gatherer.updateDurationCount++
}

func (gatherer *Gatherer) GatherTotalDuration(start time.Time) {
	gatherer.totalDuration = time.Since(start).Seconds()
}

//nolint:funlen
func (gatherer *Gatherer) Flush() error {
	geocodingDuration := 0.0
	if gatherer.geocodingDurationCount != 0 {
		geocodingDuration = gatherer.geocodingDurationSum / gatherer.geocodingDurationCount
	}
	readingDuration := 0.0
	if gatherer.readingDurationCount != 0 {
		readingDuration = gatherer.readingDurationSum / gatherer.readingDurationCount
	}
	creationDuration := 0.0
	if gatherer.creationDurationCount != 0 {
		creationDuration = gatherer.creationDurationSum / gatherer.creationDurationCount
	}
	updateDuration := 0.0
	if gatherer.updateDurationCount != 0 {
		updateDuration = gatherer.updateDurationSum / gatherer.updateDurationCount
	}
	_, err := gatherer.db.Exec(
		`insert into runs
    	(
    	    completion_time, miner, state_sanitization_number, city_sanitization_number,
    	    district_sanitization_number, swap_sanitization_number, street_sanitization_number,
    	    house_number_sanitization_number, located_geocoding_number, unlocated_geocoding_number,
    	    failed_geocoding_number, inconclusive_geocoding_number, successful_geocoding_number,
    	    approved_validation_number, denied_validation_number, created_storing_number,
    	    updated_storing_number, unaltered_storing_number, failed_storing_number, fetching_duration,
    	    geocoding_duration, reading_duration, creation_duration, update_duration, total_duration
    	)
    	values
    	(
    		now() at time zone 'utc', $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
    	    $17, $18, $19, $20, $21, $22, $23, $24
    	)`,
		gatherer.miner,
		gatherer.stateSanitizationNumber,
		gatherer.citySanitizationNumber,
		gatherer.districtSanitizationNumber,
		gatherer.swapSanitizationNumber,
		gatherer.streetSanitizationNumber,
		gatherer.houseNumberSanitizationNumber,
		gatherer.locatedGeocodingNumber,
		gatherer.unlocatedGeocodingNumber,
		gatherer.failedGeocodingNumber,
		gatherer.inconclusiveGeocodingNumber,
		gatherer.successfulGeocodingNumber,
		gatherer.approvedValidationNumber,
		gatherer.deniedValidationNumber,
		gatherer.createdStoringNumber,
		gatherer.updatedStoringNumber,
		gatherer.unalteredStoringNumber,
		gatherer.failedStoringNumber,
		gatherer.fetchingDuration,
		geocodingDuration,
		readingDuration,
		creationDuration,
		updateDuration,
		gatherer.totalDuration,
	)
	gatherer.stateSanitizationNumber = 0
	gatherer.citySanitizationNumber = 0
	gatherer.districtSanitizationNumber = 0
	gatherer.swapSanitizationNumber = 0
	gatherer.streetSanitizationNumber = 0
	gatherer.houseNumberSanitizationNumber = 0
	gatherer.locatedGeocodingNumber = 0
	gatherer.unlocatedGeocodingNumber = 0
	gatherer.failedGeocodingNumber = 0
	gatherer.inconclusiveGeocodingNumber = 0
	gatherer.successfulGeocodingNumber = 0
	gatherer.approvedValidationNumber = 0
	gatherer.deniedValidationNumber = 0
	gatherer.createdStoringNumber = 0
	gatherer.updatedStoringNumber = 0
	gatherer.unalteredStoringNumber = 0
	gatherer.failedStoringNumber = 0
	gatherer.fetchingDuration = 0
	gatherer.geocodingDurationSum = 0
	gatherer.geocodingDurationCount = 0
	gatherer.readingDurationSum = 0
	gatherer.readingDurationCount = 0
	gatherer.creationDurationSum = 0
	gatherer.creationDurationCount = 0
	gatherer.updateDurationSum = 0
	gatherer.updateDurationCount = 0
	gatherer.totalDuration = 0
	if err != nil {
		return fmt.Errorf("metrics: gatherer failed to flush, %v", err)
	}
	return nil
}
