package metrics

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewDrain(miner string, db *sql.DB, logger log.FieldLogger) *Drain {
	numbers := make(map[Number]int, len(numberViews))
	for number := range numberViews {
		numbers[number] = 0
	}
	durations := make(map[Duration]*bucket, len(durationViews))
	for duration := range durationViews {
		durations[duration] = &bucket{}
	}
	return &Drain{miner, 0, numbers, durations, db, logger}
}

type Drain struct {
	miner     string
	page      int
	numbers   map[Number]int
	durations map[Duration]*bucket
	db        *sql.DB
	logger    log.FieldLogger
}

func (d *Drain) DrainPage(page int) {
	d.page = page
}

func (d *Drain) DrainNumber(number Number) {
	if _, ok := d.numbers[number]; ok {
		d.numbers[number]++
	} else {
		d.logger.WithField("number", number).Errorf(
			"metrics: drain doesn't accept the number",
		)
	}
}

func (d *Drain) DrainDuration(duration Duration, start time.Time) {
	if b, ok := d.durations[duration]; ok {
		b.span(start)
	} else {
		d.logger.WithField("duration", duration).Errorf(
			"metrics: drain doesn't accept the duration",
		)
	}
}

func (d *Drain) Flush() {
	_, err := d.db.Exec(
		`insert into minings
		(
			completion_time, miner, page, failed_fetching_number, state_sanitation_number,
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
			$32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46
		)`,
		d.miner,
		d.page,
		d.numbers[FailedFetchingNumber],
		d.numbers[StateSanitationNumber],
		d.numbers[CitySanitationNumber],
		d.numbers[DistrictSanitationNumber],
		d.numbers[SwapSanitationNumber],
		d.numbers[StreetSanitationNumber],
		d.numbers[HouseNumberSanitationNumber],
		d.numbers[LocatedGeocodingNumber],
		d.numbers[UnlocatedGeocodingNumber],
		d.numbers[FailedGeocodingNumber],
		d.numbers[InconclusiveGeocodingNumber],
		d.numbers[SuccessfulGeocodingNumber],
		d.numbers[SubwaylessSSFGaugingNumber],
		d.numbers[FailedSSFGaugingNumber],
		d.numbers[InconclusiveSSFGaugingNumber],
		d.numbers[SuccessfulSSFGaugingNumber],
		d.numbers[FailedIZFGaugingNumber],
		d.numbers[InconclusiveIZFGaugingNumber],
		d.numbers[SuccessfulIZFGaugingNumber],
		d.numbers[FailedGZFGaugingNumber],
		d.numbers[InconclusiveGZFGaugingNumber],
		d.numbers[SuccessfulGZFGaugingNumber],
		d.numbers[ApprovedValidationNumber],
		d.numbers[UninformativeValidationNumber],
		d.numbers[SoldValidationNumber],
		d.numbers[DeniedValidationNumber],
		d.numbers[CreatedFlatStoringNumber],
		d.numbers[UpdatedFlatStoringNumber],
		d.numbers[UnalteredFlatStoringNumber],
		d.numbers[FailedFlatStoringNumber],
		d.numbers[CreatedImageStoringNumber],
		d.numbers[UnalteredImageStoringNumber],
		d.numbers[FailedImageStoringNumber],
		d.durations[FetchingDuration].avg(),
		d.durations[GeocodingDuration].avg(),
		d.durations[SSFGaugingDuration].avg(),
		d.durations[IZFGaugingDuration].avg(),
		d.durations[GZFGaugingDuration].avg(),
		d.durations[ReadingFlatStoringDuration].avg(),
		d.durations[CreationFlatStoringDuration].avg(),
		d.durations[UpdateFlatStoringDuration].avg(),
		d.durations[ReadingImageStoringDuration].avg(),
		d.durations[CreationImageStoringDuration].avg(),
		d.durations[TotalDuration].avg(),
	)
	if err != nil {
		d.logger.Errorf("metrics: drain failed to flush, %v", err)
	}
	d.page = 0
	for number := range d.numbers {
		d.numbers[number] = 0
	}
	for _, b := range d.durations {
		b.reset()
	}
}
