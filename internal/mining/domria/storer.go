package domria

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"time"
)

func newStorer(db *sql.DB) *storer {
	return &storer{db}
}

type storer struct {
	db *sql.DB
}

func (storer *storer) storeFlats(flats []*flat) error {
	length := len(flats)
	if length == 0 {
		log.Debug("domria: storer skipped flats")
		return nil
	}
	stmt, err := storer.prepare()
	if err != nil {
		return err
	}
	storedNumber, duration := 0.0, 0.0
	for _, flat := range flats {
		start := time.Now()
		if err = storer.createFlat(stmt, flat); err != nil {
			log.Error(err)
		} else {
			storedNumber++
		}
		duration += time.Since(start).Seconds()
	}
	if err = stmt.Close(); err != nil {
		return fmt.Errorf("domria: storer failed to close the stmt, %v", err)
	}
	log.Debugf("domria: storer stored %.0f flats (%.3fs)", storedNumber, duration/float64(length))
	return nil
}

func (storer *storer) prepare() (*sql.Stmt, error) {
	stmt, err := storer.db.Prepare(
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
	)
	if err != nil {
		return nil, fmt.Errorf("domria: storer failed to prepare the statement, %v", err)
	}
	return stmt, nil
}

func (storer *storer) createFlat(stmt *sql.Stmt, flat *flat) error {
	_, err := stmt.Exec(
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
		return fmt.Errorf("domria: storer failed to create a flat %s, %v", flat.originURL, err)
	}
	return nil
}
