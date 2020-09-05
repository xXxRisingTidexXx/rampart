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
	stateSanitationNumber         int
	citySanitationNumber          int
	districtSanitationNumber      int
	swapSanitationNumber          int
	streetSanitationNumber        int
	houseNumberSanitationNumber   int
	locatedGeocodingNumber        int
	unlocatedGeocodingNumber      int
	failedGeocodingNumber         int
	inconclusiveGeocodingNumber   int
	successfulGeocodingNumber     int
	subwaylessSSFGaugingNumber    int
	failedSSFGaugingNumber        int
	inconclusiveSSFGaugingNumber  int
	successfulSSFGaugingNumber    int
	failedIZFGaugingNumber        int
	inconclusiveIZFGaugingNumber  int
	successfulIZFGaugingNumber    int
	failedGZFGaugingNumber        int
	inconclusiveGZFGaugingNumber  int
	successfulGZFGaugingNumber    int
	approvedValidationNumber      int
	uninformativeValidationNumber int
	deniedValidationNumber        int
	createdStoringNumber          int
	updatedStoringNumber          int
	unalteredStoringNumber        int
	failedStoringNumber           int
	fetchingDuration              float64
	geocodingDurationSum          float64
	geocodingDurationCount        float64
	ssfGaugingDurationSum         float64
	ssfGaugingDurationCount       float64
	izfGaugingDurationSum         float64
	izfGaugingDurationCount       float64
	gzfGaugingDurationSum         float64
	gzfGaugingDurationCount       float64
	readingDurationSum            float64
	readingDurationCount          float64
	creationDurationSum           float64
	creationDurationCount         float64
	updateDurationSum             float64
	updateDurationCount           float64
	totalDuration                 float64
	db                            *sql.DB
}

func (gatherer *Gatherer) GatherStateSanitation() {
	gatherer.stateSanitationNumber++
}

func (gatherer *Gatherer) GatherCitySanitation() {
	gatherer.citySanitationNumber++
}

func (gatherer *Gatherer) GatherDistrictSanitation() {
	gatherer.districtSanitationNumber++
}

func (gatherer *Gatherer) GatherSwapSanitation() {
	gatherer.swapSanitationNumber++
}

func (gatherer *Gatherer) GatherStreetSanitation() {
	gatherer.streetSanitationNumber++
}

func (gatherer *Gatherer) GatherHouseNumberSanitation() {
	gatherer.houseNumberSanitationNumber++
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

func (gatherer *Gatherer) GatherSubwaylessSSFGauging() {
	gatherer.subwaylessSSFGaugingNumber++
}

func (gatherer *Gatherer) GatherFailedSSFGauging() {
	gatherer.failedSSFGaugingNumber++
}

func (gatherer *Gatherer) GatherInconclusiveSSFGauging() {
	gatherer.inconclusiveSSFGaugingNumber++
}

func (gatherer *Gatherer) GatherSuccessfulSSFGauging() {
	gatherer.successfulSSFGaugingNumber++
}

func (gatherer *Gatherer) GatherFailedIZFGauging() {
	gatherer.failedIZFGaugingNumber++
}

func (gatherer *Gatherer) GatherInconclusiveIZFGauging() {
	gatherer.inconclusiveIZFGaugingNumber++
}

func (gatherer *Gatherer) GatherSuccessfulIZFGauging() {
	gatherer.successfulIZFGaugingNumber++
}

func (gatherer *Gatherer) GatherFailedGZFGauging() {
	gatherer.failedGZFGaugingNumber++
}

func (gatherer *Gatherer) GatherInconclusiveGZFGauging() {
	gatherer.inconclusiveGZFGaugingNumber++
}

func (gatherer *Gatherer) GatherSuccessfulGZFGauging() {
	gatherer.successfulGZFGaugingNumber++
}

func (gatherer *Gatherer) GatherApprovedValidation() {
	gatherer.approvedValidationNumber++
}

func (gatherer *Gatherer) GatherUninformativeValidation() {
	gatherer.uninformativeValidationNumber++
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
    	    completion_time, miner, state_sanitation_number, city_sanitation_number,
    	    district_sanitation_number, swap_sanitation_number, street_sanitation_number,
    	    house_number_sanitation_number, located_geocoding_number, unlocated_geocoding_number,
    	    failed_geocoding_number, inconclusive_geocoding_number, successful_geocoding_number,
    	    approved_validation_number, uninformative_validation_number, denied_validation_number,
    	    created_storing_number, updated_storing_number, unaltered_storing_number,
    	    failed_storing_number, fetching_duration, geocoding_duration, reading_duration,
    	    creation_duration, update_duration, total_duration
    	)
    	values
    	(
    		now() at time zone 'utc', $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
    	    $17, $18, $19, $20, $21, $22, $23, $24, $25
    	)`,
		gatherer.miner,
		gatherer.stateSanitationNumber,
		gatherer.citySanitationNumber,
		gatherer.districtSanitationNumber,
		gatherer.swapSanitationNumber,
		gatherer.streetSanitationNumber,
		gatherer.houseNumberSanitationNumber,
		gatherer.locatedGeocodingNumber,
		gatherer.unlocatedGeocodingNumber,
		gatherer.failedGeocodingNumber,
		gatherer.inconclusiveGeocodingNumber,
		gatherer.successfulGeocodingNumber,
		gatherer.approvedValidationNumber,
		gatherer.uninformativeValidationNumber,
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
	gatherer.stateSanitationNumber = 0
	gatherer.citySanitationNumber = 0
	gatherer.districtSanitationNumber = 0
	gatherer.swapSanitationNumber = 0
	gatherer.streetSanitationNumber = 0
	gatherer.houseNumberSanitationNumber = 0
	gatherer.locatedGeocodingNumber = 0
	gatherer.unlocatedGeocodingNumber = 0
	gatherer.failedGeocodingNumber = 0
	gatherer.inconclusiveGeocodingNumber = 0
	gatherer.successfulGeocodingNumber = 0
	gatherer.approvedValidationNumber = 0
	gatherer.uninformativeValidationNumber = 0
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
