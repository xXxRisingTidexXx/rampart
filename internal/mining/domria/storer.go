package domria

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"rampart/internal/mining/metrics"
	"time"
)

func newStorer(db *sql.DB, gatherer *metrics.Gatherer) *storer {
	return &storer{db, gatherer}
}

type storer struct {
	db       *sql.DB
	gatherer *metrics.Gatherer
}

// TODO: add log with field "origin_url".
func (storer *storer) storeFlats(flats []*flat) {
	for _, flat := range flats {
		if err := storer.storeFlat(flat); err != nil {
			log.Error(err)
			storer.gatherer.GatherFailedStoring()
		}
	}
}

func (storer *storer) storeFlat(flat *flat) error {
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
	if flat.updateTime.After(origin.updateTime) {
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

func (storer *storer) readFlat(tx *sql.Tx, flat *flat) (*origin, error) {
	row := tx.QueryRow(`select update_time from flats where origin_url = $1`, flat.originURL)
	var origin origin
	switch err := row.Scan(&origin.updateTime); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &origin, nil
	default:
		return nil, fmt.Errorf("domria: storer failed to read flat %s, %v", flat.originURL, err)
	}
}

func (storer *storer) updateFlat(tx *sql.Tx, flat *flat) error {
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
		    point = $12,
		    state = $13,
		    city = $14,
		    district = $15,
		    street = $16,
		    house_number = $17
		where origin_url = $18`,
		flat.imageURL,
		flat.updateTime,
		flat.price,
		flat.totalArea,
		flat.livingArea,
		flat.kitchenArea,
		flat.roomNumber,
		flat.floor,
		flat.totalFloor,
		flat.housing.String(),
		flat.complex,
		&ewkb.Point{Point: flat.point},
		flat.state,
		flat.city,
		flat.district,
		flat.street,
		flat.houseNumber,
		flat.originURL,
	)
	if err != nil {
		return fmt.Errorf("domria: storer failed to update flat %s, %v", flat.originURL, err)
	}
	return nil
}

//nolint:funlen
func (storer *storer) createFlat(tx *sql.Tx, flat *flat) error {
	_, err := tx.Exec(
		`insert into flats(
                  origin_url, 
                  image_url, 
                  update_time, 
                  parsing_time, 
                  price, 
                  total_area, 
                  living_area, 
                  kitchen_area, 
                  room_number, 
                  floor, 
                  total_floor, 
                  housing, 
                  complex, 
                  point, 
                  state, 
                  city, 
                  district, 
                  street, 
                  house_number
                  ) values (
                            $1, 
                            $2, 
                            $3, 
                            now() at time zone 'utc', 
                            $4, 
                            $5, 
                            $6, 
                            $7, 
                            $8, 
                            $9, 
                            $10, 
                            $11, 
                            $12, 
                            $13, 
                            $14, 
                            $15, 
                            $16, 
                            $17, 
                            $18
                            )`,
		flat.originURL,
		flat.imageURL,
		flat.updateTime,
		flat.price,
		flat.totalArea,
		flat.livingArea,
		flat.kitchenArea,
		flat.roomNumber,
		flat.floor,
		flat.totalFloor,
		flat.housing.String(),
		flat.complex,
		&ewkb.Point{Point: flat.point},
		flat.state,
		flat.city,
		flat.district,
		flat.street,
		flat.houseNumber,
	)
	if err != nil {
		return fmt.Errorf("domria: storer failed to create flat %s, %v", flat.originURL, err)
	}
	return nil
}
