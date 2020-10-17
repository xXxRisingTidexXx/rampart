package metrics

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"
)

// TODO: add the fetcher page.
func NewDrain(miner string, db *sql.DB, logger log.FieldLogger) *Drain {
	numbers := make(map[Number]int, len(numberViews))
	for number := range numberViews {
		numbers[number] = 0
	}
	durations := make(map[Duration]*bucket, len(durationViews))
	for duration := range durationViews {
		durations[duration] = &bucket{}
	}
	return &Drain{miner, numbers, durations, db, logger}
}

type Drain struct {
	miner     string
	numbers   map[Number]int
	durations map[Duration]*bucket
	db        *sql.DB
	logger    log.FieldLogger
}

func (drain *Drain) DrainNumber(number Number) {
	if _, ok := drain.numbers[number]; ok {
		drain.numbers[number]++
	} else {
		drain.logger.WithField("number", number).Errorf(
			"metrics: drain doesn't accept the number",
		)
	}
}

func (drain *Drain) DrainDuration(duration Duration, start time.Time) {
	if b, ok := drain.durations[duration]; ok {
		b.span(start)
	} else {
		drain.logger.WithField("duration", duration).Errorf(
			"metrics: drain doesn't accept the duration",
		)
	}
}

func (drain *Drain) Flush() {
	_, err := drain.db.Exec(
		`insert into minings
		(
			completion_time, miner, failed_fetching_number, state_sanitation_number,
			city_sanitation_number, district_sanitation_number, swap_sanitation_number,
			street_sanitation_number, house_number_sanitation_number, located_geocoding_number,
		 	unlocated_geocoding_number,failed_geocoding_number, inconclusive_geocoding_number,
		 	successful_geocoding_number, subwayless_ssf_gauging_number, failed_ssf_gauging_number,
		 	inconclusive_ssf_gauging_number, successful_ssf_gauging_number,
		 	failed_izf_gauging_number, inconclusive_izf_gauging_number,
		 	successful_izf_gauging_number, failed_gzf_gauging_number,
		 	inconclusive_gzf_gauging_number, successful_gzf_gauging_number,
		 	approved_validation_number, uninformative_validation_number, sold_validation_number,
		 	denied_validation_number, created_flat_storing_number, updated_flat_storing_number,
		 	unaltered_flat_storing_number, failed_flat_storing_number,
		 	created_image_storing_number, unaltered_image_storing_number,
		 	failed_image_storing_number, fetching_duration, geocoding_duration,
		 	ssf_gauging_duration, izf_gauging_duration, gzf_gauging_duration,
		 	reading_flat_storing_duration, creation_flat_storing_duration,
		 	update_flat_storing_duration, reading_image_storing_duration,
		 	creation_image_storing_duration, total_duration
		)
		values
		(
			now() at time zone 'utc', $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31,
			$32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45
		)`,
		drain.miner,
		drain.numbers[FailedFetchingNumber],
		drain.numbers[StateSanitationNumber],
		drain.numbers[CitySanitationNumber],
		drain.numbers[DistrictSanitationNumber],
		drain.numbers[SwapSanitationNumber],
		drain.numbers[StreetSanitationNumber],
		drain.numbers[HouseNumberSanitationNumber],
		drain.numbers[LocatedGeocodingNumber],
		drain.numbers[UnlocatedGeocodingNumber],
		drain.numbers[FailedGeocodingNumber],
		drain.numbers[InconclusiveGeocodingNumber],
		drain.numbers[SuccessfulGeocodingNumber],
		drain.numbers[SubwaylessSSFGaugingNumber],
		drain.numbers[FailedSSFGaugingNumber],
		drain.numbers[InconclusiveSSFGaugingNumber],
		drain.numbers[SuccessfulSSFGaugingNumber],
		drain.numbers[FailedIZFGaugingNumber],
		drain.numbers[InconclusiveIZFGaugingNumber],
		drain.numbers[SuccessfulIZFGaugingNumber],
		drain.numbers[FailedGZFGaugingNumber],
		drain.numbers[InconclusiveGZFGaugingNumber],
		drain.numbers[SuccessfulGZFGaugingNumber],
		drain.numbers[ApprovedValidationNumber],
		drain.numbers[UninformativeValidationNumber],
		drain.numbers[SoldValidationNumber],
		drain.numbers[DeniedValidationNumber],
		drain.numbers[CreatedFlatStoringNumber],
		drain.numbers[UpdatedFlatStoringNumber],
		drain.numbers[UnalteredFlatStoringNumber],
		drain.numbers[FailedFlatStoringNumber],
		drain.numbers[CreatedImageStoringNumber],
		drain.numbers[UnalteredImageStoringNumber],
		drain.numbers[FailedImageStoringNumber],
		drain.durations[FetchingDuration].avg(),
		drain.durations[GeocodingDuration].avg(),
		drain.durations[SSFGaugingDuration].avg(),
		drain.durations[IZFGaugingDuration].avg(),
		drain.durations[GZFGaugingDuration].avg(),
		drain.durations[ReadingFlatStoringDuration].avg(),
		drain.durations[CreationFlatStoringDuration].avg(),
		drain.durations[UpdateFlatStoringDuration].avg(),
		drain.durations[ReadingImageStoringDuration].avg(),
		drain.durations[CreationImageStoringDuration].avg(),
		drain.durations[TotalDuration].avg(),
	)
	if err != nil {
		drain.logger.Errorf("metrics: drain failed to flush, %v", err)
	}
	for number := range drain.numbers {
		drain.numbers[number] = 0
	}
	for _, b := range drain.durations {
		b.reset()
	}
}
