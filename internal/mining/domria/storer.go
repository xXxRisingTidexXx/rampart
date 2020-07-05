package domria

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"rampart/internal/config"
	"time"
)

func newStorer(db *sql.DB, config *config.Storer) *storer {
	return &storer{db, time.Duration(config.UpdateTiming)}
}

type storer struct {
	db           *sql.DB
	updateTiming time.Duration
}

func (storer *storer) storeFlats(flats []*flat) {
	for _, flat := range flats {
		if err := storer.storeFlat(flat); err != nil {
			log.Error(err)  // TODO: add log with field "origin_url"
		}
	}
}

func (storer *storer) storeFlat(flat *flat) error {
	tx, err := storer.db.Begin()
	if err != nil {
		return fmt.Errorf("domria: storer failed to begin a transaction, %v", err)
	}
	origin, err := storer.readFlat(tx, flat)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if origin != nil {
		isNewer := flat.updateTime.Sub(origin.updateTime) >= storer.updateTiming
		if isNewer || !isNewer && flat.price < origin.price {
			if err = storer.updateFlat(tx, flat); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	} else if err = storer.createFlat(tx, flat); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("domria: storer failed to commit a transaction, %v", err)
	}
	return nil
}

func (storer *storer) readFlat(tx *sql.Tx, flat *flat) (*origin, error) {
	row := tx.QueryRow(`select update_time, price from flats where origin_url = $1`, flat.originURL)
	var origin origin
	switch err := row.Scan(&origin.updateTime, &origin.price); err {
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
