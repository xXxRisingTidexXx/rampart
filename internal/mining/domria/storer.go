package domria

import (
	"database/sql"
	"fmt"
	"github.com/paulmach/orb/encoding/wkb"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"time"
)

func NewStorer(
	config *config.Storer,
	db *sql.DB,
	gatherer *metrics.Gatherer,
	logger log.FieldLogger,
) *Storer {
	return &Storer{config.SRID, db, gatherer, logger}
}

type Storer struct {
	srid     int
	db       *sql.DB
	gatherer *metrics.Gatherer
	logger   log.FieldLogger
}

func (storer *Storer) StoreFlats(flats []*Flat) {
	for _, flat := range flats {
		if err := storer.storeFlat(flat); err != nil {
			storer.logger.WithField("origin_url", flat.OriginURL).Error(err)
			storer.gatherer.GatherFailedStoring()
		}
	}
}

func (storer *Storer) storeFlat(flat *Flat) error {
	tx, err := storer.db.Begin()
	if err != nil {
		return fmt.Errorf("domria: storer failed to begin a transaction, %v", err)
	}
	start := time.Now()
	origin, err := storer.readFlat(tx, flat)
	storer.gatherer.GatherReadingDuration(start)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	message := "domria: storer failed to commit a transaction, %v"
	if origin == nil {
		start := time.Now()
		err = storer.createFlat(tx, flat)
		storer.gatherer.GatherCreationDuration(start)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return fmt.Errorf(message, err)
		}
		storer.gatherer.GatherCreatedStoring()
		return nil
	}
	if flat.UpdateTime.After(origin.updateTime) {
		start := time.Now()
		err = storer.updateFlat(tx, flat)
		storer.gatherer.GatherUpdateDuration(start)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return fmt.Errorf(message, err)
		}
		storer.gatherer.GatherUpdatedStoring()
		return nil
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf(message, err)
	}
	storer.gatherer.GatherUnalteredStoring()
	return nil
}

func (storer *Storer) readFlat(tx *sql.Tx, flat *Flat) (*origin, error) {
	row := tx.QueryRow(`select update_time from flats where origin_url = $1`, flat.OriginURL)
	var origin origin
	switch err := row.Scan(&origin.updateTime); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &origin, nil
	default:
		return nil, fmt.Errorf("domria: storer failed to read flat %s, %v", flat.OriginURL, err)
	}
}

func (storer *Storer) updateFlat(tx *sql.Tx, flat *Flat) error {
	_, err := tx.Exec(
		`update flats 
		set image_url = $1,
		    update_time = $2,
		    parsing_time = now() at time zone 'utc',
		    price = $3,
		    total_area = $4,
		    living_area = $5,
		    kitchen_area = $6,
		    room_number = $7,
		    floor = $8,
		    total_floor = $9,
		    housing = $10,
		    complex = $11,
		    point = st_geomfromwkb($12, $13),
		    state = $14,
		    city = $15,
		    district = $16,
		    street = $17,
		    house_number = $18
		where origin_url = $19`,
		flat.ImageURL,
		flat.UpdateTime,
		flat.Price,
		flat.TotalArea,
		flat.LivingArea,
		flat.KitchenArea,
		flat.RoomNumber,
		flat.Floor,
		flat.TotalFloor,
		flat.Housing.String(),
		flat.Complex,
		wkb.Value(flat.Point),
		storer.srid,
		flat.State,
		flat.City,
		flat.District,
		flat.Street,
		flat.HouseNumber,
		flat.OriginURL,
	)
	if err != nil {
		return fmt.Errorf("domria: storer failed to update flat %s, %v", flat.OriginURL, err)
	}
	return nil
}

func (storer *Storer) createFlat(tx *sql.Tx, flat *Flat) error {
	_, err := tx.Exec(
		`insert into flats
        (
         	origin_url, image_url, update_time, parsing_time, price, total_area, living_area, kitchen_area,
            room_number, floor, total_floor, housing, complex, point, state, city, district, street,
            house_number
        )
        values 
		(
		    $1, $2, $3, now() at time zone 'utc', $4, $5, $6, $7, $8, $9, $10, $11, $12, 
		    st_geomfromwkb($13, $14), $15, $16, $17, $18, $19
		)`,
		flat.OriginURL,
		flat.ImageURL,
		flat.UpdateTime,
		flat.Price,
		flat.TotalArea,
		flat.LivingArea,
		flat.KitchenArea,
		flat.RoomNumber,
		flat.Floor,
		flat.TotalFloor,
		flat.Housing.String(),
		flat.Complex,
		wkb.Value(flat.Point),
		storer.srid,
		flat.State,
		flat.City,
		flat.District,
		flat.Street,
		flat.HouseNumber,
	)
	if err != nil {
		return fmt.Errorf("domria: storer failed to create flat %s, %v", flat.OriginURL, err)
	}
	return nil
}
