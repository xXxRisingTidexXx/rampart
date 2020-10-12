package metrics

import (
	"database/sql"
	"fmt"
	"time"
)

func NewDrain(miner string, db *sql.DB) *Drain {
	return &Drain{miner: miner, db: db}
}

type Drain struct {
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

func (drain *Drain) GatherStateSanitation() {
	drain.stateSanitationNumber++
}

func (drain *Drain) GatherCitySanitation() {
	drain.citySanitationNumber++
}

func (drain *Drain) GatherDistrictSanitation() {
	drain.districtSanitationNumber++
}

func (drain *Drain) GatherSwapSanitation() {
	drain.swapSanitationNumber++
}

func (drain *Drain) GatherStreetSanitation() {
	drain.streetSanitationNumber++
}

func (drain *Drain) GatherHouseNumberSanitation() {
	drain.houseNumberSanitationNumber++
}

func (drain *Drain) GatherLocatedGeocoding() {
	drain.locatedGeocodingNumber++
}

func (drain *Drain) GatherUnlocatedGeocoding() {
	drain.unlocatedGeocodingNumber++
}

func (drain *Drain) GatherFailedGeocoding() {
	drain.failedGeocodingNumber++
}

func (drain *Drain) GatherInconclusiveGeocoding() {
	drain.inconclusiveGeocodingNumber++
}

func (drain *Drain) GatherSuccessfulGeocoding() {
	drain.successfulGeocodingNumber++
}

func (drain *Drain) GatherSubwaylessSSFGauging() {
	drain.subwaylessSSFGaugingNumber++
}

func (drain *Drain) GatherFailedSSFGauging() {
	drain.failedSSFGaugingNumber++
}

func (drain *Drain) GatherInconclusiveSSFGauging() {
	drain.inconclusiveSSFGaugingNumber++
}

func (drain *Drain) GatherSuccessfulSSFGauging() {
	drain.successfulSSFGaugingNumber++
}

func (drain *Drain) GatherFailedIZFGauging() {
	drain.failedIZFGaugingNumber++
}

func (drain *Drain) GatherInconclusiveIZFGauging() {
	drain.inconclusiveIZFGaugingNumber++
}

func (drain *Drain) GatherSuccessfulIZFGauging() {
	drain.successfulIZFGaugingNumber++
}

func (drain *Drain) GatherFailedGZFGauging() {
	drain.failedGZFGaugingNumber++
}

func (drain *Drain) GatherInconclusiveGZFGauging() {
	drain.inconclusiveGZFGaugingNumber++
}

func (drain *Drain) GatherSuccessfulGZFGauging() {
	drain.successfulGZFGaugingNumber++
}

func (drain *Drain) GatherApprovedValidation() {
	drain.approvedValidationNumber++
}

func (drain *Drain) GatherUninformativeValidation() {
	drain.uninformativeValidationNumber++
}

func (drain *Drain) GatherDeniedValidation() {
	drain.deniedValidationNumber++
}

func (drain *Drain) GatherCreatedStoring() {
	drain.createdStoringNumber++
}

func (drain *Drain) GatherUpdatedStoring() {
	drain.updatedStoringNumber++
}

func (drain *Drain) GatherUnalteredStoring() {
	drain.unalteredStoringNumber++
}

func (drain *Drain) GatherFailedStoring() {
	drain.failedStoringNumber++
}

func (drain *Drain) GatherFetchingDuration(start time.Time) {
	drain.fetchingDuration = time.Since(start).Seconds()
}

func (drain *Drain) GatherGeocodingDuration(start time.Time) {
	drain.geocodingDurationSum += time.Since(start).Seconds()
	drain.geocodingDurationCount++
}

func (drain *Drain) GatherSSFGaugingDuration(start time.Time) {
	drain.ssfGaugingDurationSum += time.Since(start).Seconds()
	drain.ssfGaugingDurationCount++
}

func (drain *Drain) GatherIZFGaugingDuration(start time.Time) {
	drain.izfGaugingDurationSum += time.Since(start).Seconds()
	drain.izfGaugingDurationCount++
}

func (drain *Drain) GatherGZFGaugingDuration(start time.Time) {
	drain.gzfGaugingDurationSum += time.Since(start).Seconds()
	drain.gzfGaugingDurationCount++
}

func (drain *Drain) GatherReadingDuration(start time.Time) {
	drain.readingDurationSum += time.Since(start).Seconds()
	drain.readingDurationCount++
}

func (drain *Drain) GatherCreationDuration(start time.Time) {
	drain.creationDurationSum += time.Since(start).Seconds()
	drain.creationDurationCount++
}

func (drain *Drain) GatherUpdateDuration(start time.Time) {
	drain.updateDurationSum += time.Since(start).Seconds()
	drain.updateDurationCount++
}

func (drain *Drain) GatherTotalDuration(start time.Time) {
	drain.totalDuration = time.Since(start).Seconds()
}

func (drain *Drain) Flush() error {
	geocodingDuration := 0.0
	if drain.geocodingDurationCount != 0 {
		geocodingDuration = drain.geocodingDurationSum / drain.geocodingDurationCount
	}
	ssfGaugingDuration := 0.0
	if drain.ssfGaugingDurationCount != 0 {
		ssfGaugingDuration = drain.ssfGaugingDurationSum / drain.ssfGaugingDurationCount
	}
	izfGaugingDuration := 0.0
	if drain.izfGaugingDurationCount != 0 {
		izfGaugingDuration = drain.izfGaugingDurationSum / drain.izfGaugingDurationCount
	}
	gzfGaugingDuration := 0.0
	if drain.gzfGaugingDurationCount != 0 {
		gzfGaugingDuration = drain.gzfGaugingDurationSum / drain.gzfGaugingDurationCount
	}
	readingDuration := 0.0
	if drain.readingDurationCount != 0 {
		readingDuration = drain.readingDurationSum / drain.readingDurationCount
	}
	creationDuration := 0.0
	if drain.creationDurationCount != 0 {
		creationDuration = drain.creationDurationSum / drain.creationDurationCount
	}
	updateDuration := 0.0
	if drain.updateDurationCount != 0 {
		updateDuration = drain.updateDurationSum / drain.updateDurationCount
	}
	_, err := drain.db.Exec(
		`insert into runs
    	(
    	    completion_time, miner, state_sanitation_number, city_sanitation_number,
    	    district_sanitation_number, swap_sanitation_number, street_sanitation_number,
    	    house_number_sanitation_number, located_geocoding_number, unlocated_geocoding_number,
    	    failed_geocoding_number, inconclusive_geocoding_number, successful_geocoding_number,
    	    subwayless_ssf_gauging_number, failed_ssf_gauging_number, inconclusive_ssf_gauging_number,
    	    successful_ssf_gauging_number, failed_izf_gauging_number, inconclusive_izf_gauging_number,
    	    successful_izf_gauging_number, failed_gzf_gauging_number, inconclusive_gzf_gauging_number,
    	    successful_gzf_gauging_number, approved_validation_number, uninformative_validation_number,
    	    denied_validation_number, created_storing_number, updated_storing_number,
    	    unaltered_storing_number, failed_storing_number, fetching_duration, geocoding_duration,
    	    ssf_gauging_duration, izf_gauging_duration, gzf_gauging_duration, reading_duration,
    	    creation_duration, update_duration, total_duration
    	)
    	values
    	(
    		now() at time zone 'utc', $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
    	    $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31,
    	    $32, $33, $34, $35, $36, $37, $38
    	)`,
		drain.miner,
		drain.stateSanitationNumber,
		drain.citySanitationNumber,
		drain.districtSanitationNumber,
		drain.swapSanitationNumber,
		drain.streetSanitationNumber,
		drain.houseNumberSanitationNumber,
		drain.locatedGeocodingNumber,
		drain.unlocatedGeocodingNumber,
		drain.failedGeocodingNumber,
		drain.inconclusiveGeocodingNumber,
		drain.successfulGeocodingNumber,
		drain.subwaylessSSFGaugingNumber,
		drain.failedSSFGaugingNumber,
		drain.inconclusiveSSFGaugingNumber,
		drain.successfulSSFGaugingNumber,
		drain.failedIZFGaugingNumber,
		drain.inconclusiveIZFGaugingNumber,
		drain.successfulIZFGaugingNumber,
		drain.failedGZFGaugingNumber,
		drain.inconclusiveGZFGaugingNumber,
		drain.successfulGZFGaugingNumber,
		drain.approvedValidationNumber,
		drain.uninformativeValidationNumber,
		drain.deniedValidationNumber,
		drain.createdStoringNumber,
		drain.updatedStoringNumber,
		drain.unalteredStoringNumber,
		drain.failedStoringNumber,
		drain.fetchingDuration,
		geocodingDuration,
		ssfGaugingDuration,
		izfGaugingDuration,
		gzfGaugingDuration,
		readingDuration,
		creationDuration,
		updateDuration,
		drain.totalDuration,
	)
	drain.stateSanitationNumber = 0
	drain.citySanitationNumber = 0
	drain.districtSanitationNumber = 0
	drain.swapSanitationNumber = 0
	drain.streetSanitationNumber = 0
	drain.houseNumberSanitationNumber = 0
	drain.locatedGeocodingNumber = 0
	drain.unlocatedGeocodingNumber = 0
	drain.failedGeocodingNumber = 0
	drain.inconclusiveGeocodingNumber = 0
	drain.successfulGeocodingNumber = 0
	drain.subwaylessSSFGaugingNumber = 0
	drain.failedSSFGaugingNumber = 0
	drain.inconclusiveSSFGaugingNumber = 0
	drain.successfulSSFGaugingNumber = 0
	drain.failedIZFGaugingNumber = 0
	drain.inconclusiveIZFGaugingNumber = 0
	drain.successfulIZFGaugingNumber = 0
	drain.failedGZFGaugingNumber = 0
	drain.inconclusiveGZFGaugingNumber = 0
	drain.successfulGZFGaugingNumber = 0
	drain.approvedValidationNumber = 0
	drain.uninformativeValidationNumber = 0
	drain.deniedValidationNumber = 0
	drain.createdStoringNumber = 0
	drain.updatedStoringNumber = 0
	drain.unalteredStoringNumber = 0
	drain.failedStoringNumber = 0
	drain.fetchingDuration = 0
	drain.geocodingDurationSum = 0
	drain.geocodingDurationCount = 0
	drain.ssfGaugingDurationSum = 0
	drain.ssfGaugingDurationCount = 0
	drain.izfGaugingDurationSum = 0
	drain.izfGaugingDurationCount = 0
	drain.gzfGaugingDurationSum = 0
	drain.gzfGaugingDurationCount = 0
	drain.readingDurationSum = 0
	drain.readingDurationCount = 0
	drain.creationDurationSum = 0
	drain.creationDurationCount = 0
	drain.updateDurationSum = 0
	drain.updateDurationCount = 0
	drain.totalDuration = 0
	if err != nil {
		return fmt.Errorf("metrics: drain failed to flush, %v", err)
	}
	return nil
}
