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
	config config.Storer,
	db *sql.DB,
	drain *metrics.Drain,
	logger log.FieldLogger,
) *Storer {
	return &Storer{config.SRID, db, drain, logger}
}

type Storer struct {
	srid   int
	db     *sql.DB
	drain  *metrics.Drain
	logger log.FieldLogger
}

func (storer *Storer) StoreFlats(flats []Flat) {
	for _, flat := range flats {
		if err := storer.storeFlat(flat); err != nil {
			storer.logger.WithFields(
				log.Fields{"source": flat.Source, "url": flat.OriginURL},
			).Error(err)
			storer.drain.DrainNumber(metrics.FailedStoringNumber)
		}
	}
}

func (storer *Storer) storeFlat(flat Flat) error {
	tx, err := storer.db.Begin()
	if err != nil {
		return fmt.Errorf("domria: storer failed to begin a transaction, %v", err)
	}
	start := time.Now()
	o, err := storer.readFlat(tx, flat)
	storer.drain.DrainDuration(metrics.ReadingDuration, start)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	message := "domria: storer failed to commit a transaction, %v"
	if !o.isFound {
		start := time.Now()
		err = storer.createFlat(tx, flat)
		storer.drain.DrainDuration(metrics.CreationDuration, start)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return fmt.Errorf(message, err)
		}
		storer.drain.DrainNumber(metrics.CreatedStoringNumber)
		return nil
	}
	if flat.UpdateTime.After(o.updateTime) {
		start := time.Now()
		err = storer.updateFlat(tx, flat)
		storer.drain.DrainDuration(metrics.UpdateDuration, start)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return fmt.Errorf(message, err)
		}
		storer.drain.DrainNumber(metrics.UpdatedStoringNumber)
		return nil
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf(message, err)
	}
	storer.drain.DrainNumber(metrics.UnalteredStoringNumber)
	return nil
}

func (storer *Storer) readFlat(tx *sql.Tx, flat Flat) (origin, error) {
	row := tx.QueryRow(`select update_time from flats where origin_url = $1`, flat.OriginURL)
	o := origin{}
	switch err := row.Scan(&o.updateTime); err {
	case sql.ErrNoRows:
		return o, nil
	case nil:
		o.isFound = true
		return o, nil
	default:
		return o, fmt.Errorf("domria: storer failed to read the flat, %v", err)
	}
}

func (storer *Storer) updateFlat(tx *sql.Tx, flat Flat) error {
	_, err := tx.Exec(
		`update flats 
		set image_url = $2,
		    update_time = $3,
		    parsing_time = now() at time zone 'utc',
		    price = $4,
		    total_area = $5,
		    living_area = $6,
		    kitchen_area = $7,
		    room_number = $8,
		    floor = $9,
		    total_floor = $10,
		    housing = $11,
		    complex = $12,
		    point = st_geomfromwkb($13, $14),
		    state = $15,
		    city = $16,
		    district = $17,
		    street = $18,
		    house_number = $19,
		    ssf = $20,
		    izf = $21,
		    gzf = $22
		where origin_url = $1`,
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
		flat.Housing,
		flat.Complex,
		wkb.Value(flat.Point),
		storer.srid,
		flat.State,
		flat.City,
		flat.District,
		flat.Street,
		flat.HouseNumber,
		flat.SSF,
		flat.IZF,
		flat.GZF,
	)
	if err != nil {
		return fmt.Errorf("domria: storer failed to update the flat, %v", err)
	}
	return nil
}

func (storer *Storer) createFlat(tx *sql.Tx, flat Flat) error {
	_, err := tx.Exec(
		`insert into flats
        (
         	origin_url, image_url, update_time, parsing_time, price, total_area, living_area,
            kitchen_area, room_number, floor, total_floor, housing, complex, point, state, city,
            district, street, house_number, ssf, izf, gzf
        )
        values 
		(
		    $1, $2, $3, now() at time zone 'utc', $4, $5, $6, $7, $8, $9, $10, $11, $12, 
		    st_geomfromwkb($13, $14), $15, $16, $17, $18, $19, $20, $21, $22
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
		flat.Housing,
		flat.Complex,
		wkb.Value(flat.Point),
		storer.srid,
		flat.State,
		flat.City,
		flat.District,
		flat.Street,
		flat.HouseNumber,
		flat.SSF,
		flat.IZF,
		flat.GZF,
	)
	if err != nil {
		return fmt.Errorf("domria: storer failed to create the flat, %v", err)
	}
	return nil
}
