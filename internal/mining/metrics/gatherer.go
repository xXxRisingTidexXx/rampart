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
	miner                               string
	stateSanitationNumber               int
	citySanitationNumber                int
	districtSanitationNumber            int
	swapSanitationNumber                int
	streetSanitationNumber              int
	houseNumberSanitationNumber         int
	locatedGeocodingNumber              int
	unlocatedGeocodingNumber            int
	failedGeocodingNumber               int
	inconclusiveGeocodingNumber         int
	successfulGeocodingNumber           int
	absentSubwayGaugingNumber           int
	failedSubwayGaugingNumber           int
	inconclusiveSubwayGaugingNumber     int
	successfulSubwayGaugingNumber       int
	failedIndustrialGaugingNumber       int
	inconclusiveIndustrialGaugingNumber int
	successfulIndustrialGaugingNumber   int
	failedGreenGaugingNumber            int
	inconclusiveGreenGaugingNumber      int
	successfulGreenGaugingNumber        int
	approvedValidationNumber            int
	deniedValidationNumber              int
	createdStoringNumber                int
	updatedStoringNumber                int
	unalteredStoringNumber              int
	failedStoringNumber                 int
	fetchingDuration                    float64
	geocodingDurationSum                float64
	geocodingDurationCount              float64
	subwayGaugingDurationSum            float64
	subwayGaugingDurationCount          float64
	industrialGaugingDurationSum        float64
	industrialGaugingDurationCount      float64
	greenGaugingDurationSum             float64
	greenGaugingDurationCount           float64
	readingDurationSum                  float64
	readingDurationCount                float64
	creationDurationSum                 float64
	creationDurationCount               float64
	updateDurationSum                   float64
	updateDurationCount                 float64
	totalDuration                       float64
	db                                  *sql.DB
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

func (gatherer *Gatherer) GatherAbsentSubwayGauging() {
	gatherer.absentSubwayGaugingNumber++
}

func (gatherer *Gatherer) GatherFailedSubwayGauging() {
	gatherer.failedSubwayGaugingNumber++
}

func (gatherer *Gatherer) GatherInconclusiveSubwayGauging() {
	gatherer.inconclusiveSubwayGaugingNumber++
}

func (gatherer *Gatherer) GatherSuccessfulSubwayGauging() {
	gatherer.successfulSubwayGaugingNumber++
}

func (gatherer *Gatherer) GatherFailedIndustrialGauging() {
	gatherer.failedIndustrialGaugingNumber++
}

func (gatherer *Gatherer) GatherInconclusiveIndustrialGauging() {
	gatherer.inconclusiveIndustrialGaugingNumber++
}

func (gatherer *Gatherer) GatherSuccessfulIndustrialGauging() {
	gatherer.successfulIndustrialGaugingNumber++
}

func (gatherer *Gatherer) GatherFailedGreenGauging() {
	gatherer.failedGreenGaugingNumber++
}

func (gatherer *Gatherer) GatherInconclusiveGreenGauging() {
	gatherer.inconclusiveGreenGaugingNumber++
}

func (gatherer *Gatherer) GatherSuccessfulGreenGauging() {
	gatherer.successfulGreenGaugingNumber++
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

func (gatherer *Gatherer) GatherSubwayGaugingDuration(start time.Time) {
	gatherer.subwayGaugingDurationSum += time.Since(start).Seconds()
	gatherer.subwayGaugingDurationCount++
}

func (gatherer *Gatherer) GatherIndustrialGaugingDuration(start time.Time) {
	gatherer.industrialGaugingDurationSum += time.Since(start).Seconds()
	gatherer.industrialGaugingDurationCount++
}

func (gatherer *Gatherer) GatherGreenGaugingDuration(start time.Time) {
	gatherer.greenGaugingDurationSum += time.Since(start).Seconds()
	gatherer.greenGaugingDurationCount++
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
	subwayGaugingDuration := 0.0
	if gatherer.subwayGaugingDurationCount != 0 {
		subwayGaugingDuration = gatherer.subwayGaugingDurationSum / gatherer.subwayGaugingDurationCount
	}
	industrialGaugingDuration := 0.0
	if gatherer.industrialGaugingDurationCount != 0 {
		industrialGaugingDuration = gatherer.industrialGaugingDurationSum / gatherer.industrialGaugingDurationCount
	}
	greenGaugingDuration := 0.0
	if gatherer.greenGaugingDurationCount != 0 {
		greenGaugingDuration = gatherer.greenGaugingDurationSum / gatherer.greenGaugingDurationCount
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
    	    absent_subway_gauging_number, failed_subway_gauging_number, inconclusive_subway_gauging_number,
    	    successful_subway_gauging_number, failed_industrial_gauging_number,
    	    inconclusive_industrial_gauging_number, successful_industrial_gauging_number,
    	    failed_green_gauging_number, inconclusive_green_gauging_number, successful_green_gauging_number,
    	    approved_validation_number, denied_validation_number, created_storing_number,
    	    updated_storing_number, unaltered_storing_number, failed_storing_number, fetching_duration,
    	    geocoding_duration, subway_gauging_duration, industrial_gauging_duration, green_gauging_duration,
    	    reading_duration, creation_duration, update_duration, total_duration
    	)
    	values
    	(
    		now() at time zone 'utc', $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
    	    $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35,
    	    $36, $37
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
		gatherer.absentSubwayGaugingNumber,
		gatherer.failedSubwayGaugingNumber,
		gatherer.inconclusiveSubwayGaugingNumber,
		gatherer.successfulSubwayGaugingNumber,
		gatherer.failedIndustrialGaugingNumber,
		gatherer.inconclusiveIndustrialGaugingNumber,
		gatherer.successfulIndustrialGaugingNumber,
		gatherer.failedGreenGaugingNumber,
		gatherer.inconclusiveGreenGaugingNumber,
		gatherer.successfulGreenGaugingNumber,
		gatherer.approvedValidationNumber,
		gatherer.deniedValidationNumber,
		gatherer.createdStoringNumber,
		gatherer.updatedStoringNumber,
		gatherer.unalteredStoringNumber,
		gatherer.failedStoringNumber,
		gatherer.fetchingDuration,
		geocodingDuration,
		subwayGaugingDuration,
		industrialGaugingDuration,
		greenGaugingDuration,
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
	gatherer.absentSubwayGaugingNumber = 0
	gatherer.failedSubwayGaugingNumber = 0
	gatherer.inconclusiveSubwayGaugingNumber = 0
	gatherer.successfulSubwayGaugingNumber = 0
	gatherer.failedIndustrialGaugingNumber = 0
	gatherer.inconclusiveIndustrialGaugingNumber = 0
	gatherer.successfulIndustrialGaugingNumber = 0
	gatherer.failedGreenGaugingNumber = 0
	gatherer.inconclusiveGreenGaugingNumber = 0
	gatherer.successfulGreenGaugingNumber = 0
	gatherer.approvedValidationNumber = 0
	gatherer.deniedValidationNumber = 0
	gatherer.createdStoringNumber = 0
	gatherer.updatedStoringNumber = 0
	gatherer.unalteredStoringNumber = 0
	gatherer.failedStoringNumber = 0
	gatherer.fetchingDuration = 0
	gatherer.geocodingDurationSum = 0
	gatherer.geocodingDurationCount = 0
	gatherer.subwayGaugingDurationSum = 0
	gatherer.subwayGaugingDurationCount = 0
	gatherer.industrialGaugingDurationSum = 0
	gatherer.industrialGaugingDurationCount = 0
	gatherer.greenGaugingDurationSum = 0
	gatherer.greenGaugingDurationCount = 0
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
