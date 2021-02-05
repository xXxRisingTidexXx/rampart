package domria

import (
	"database/sql"
	"fmt"
	"github.com/paulmach/orb/encoding/wkb"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"time"
)

// TODO: shorten house number column (but research the actual width before).
// TODO: should we use flats' event sourcing? This means we don't mutate rows,
//  but add updated ones to avoid inconsistent lookups. We achieve this thing
//  using flat statuses.
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

func (s *Storer) StoreFlats(flats []Flat) {
	for _, flat := range flats {
		entry := s.logger.WithField("url", flat.URL)
		if id, err := s.storeFlat(flat); err != nil {
			entry.Error(err)
			s.drain.DrainNumber(metrics.FailedFlatStoringNumber)
		} else {
			s.storeImages(id, flat.Photos, PhotoKind, entry)
			s.storeImages(id, flat.Panoramas, PanoramaKind, entry)
		}
	}
}

func (s *Storer) storeFlat(flat Flat) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("domria: storer failed to begin a transaction, %v", err)
	}
	start := time.Now()
	o, err := s.readFlat(tx, flat)
	s.drain.DrainDuration(metrics.ReadingFlatStoringDuration, start)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	id, number := o.id, metrics.UnalteredFlatStoringNumber
	if !o.isFound {
		start := time.Now()
		id, err = s.createFlat(tx, flat)
		s.drain.DrainDuration(metrics.CreationFlatStoringDuration, start)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
		number = metrics.CreatedFlatStoringNumber
	} else if flat.UpdateTime.After(o.upsertTime) {
		start := time.Now()
		err := s.updateFlat(tx, flat)
		s.drain.DrainDuration(metrics.UpdateFlatStoringDuration, start)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
		number = metrics.UpdatedFlatStoringNumber
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("domria: storer failed to commit a transaction, %v", err)
	}
	s.drain.DrainNumber(number)
	return id, nil
}

func (s *Storer) readFlat(tx *sql.Tx, flat Flat) (origin, error) {
	var o origin
	row := tx.QueryRow(`select id, upsert_time from flats where url = $1`, flat.URL)
	switch err := row.Scan(&o.id, &o.upsertTime); err {
	case sql.ErrNoRows:
		return o, nil
	case nil:
		o.isFound = true
		return o, nil
	default:
		return o, fmt.Errorf("domria: storer failed to read a flat, %v", err)
	}
}

func (s *Storer) updateFlat(tx *sql.Tx, flat Flat) error {
	_, err := tx.Exec(
		`update flats 
		set upsert_time = now() at time zone 'utc',
		    price = $2,
		    total_area = $3,
		    living_area = $4,
		    kitchen_area = $5,
		    room_number = $6,
		    floor = $7,
		    total_floor = $8,
		    housing = $9,
		    complex = $10,
		    point = st_geomfromwkb($11, $12),
		    state = $13,
		    city = $14,
		    district = $15,
		    street = $16,
		    house_number = $17,
		    ssf = $18,
		    izf = $19,
		    gzf = $20
		where url = $1`,
		flat.URL,
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
		s.srid,
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
		return fmt.Errorf("domria: storer failed to update a flat, %v", err)
	}
	return nil
}

func (s *Storer) createFlat(tx *sql.Tx, flat Flat) (int, error) {
	var id int
	err := tx.QueryRow(
		`insert into flats
		(
		 	url, upsert_time, price, total_area, living_area, kitchen_area, room_number, floor,
		    total_floor, housing, complex, point, state, city, district, street, house_number, ssf,
		    izf, gzf
		)
		values 
		(
		    $1, now() at time zone 'utc', $2, $3, $4, $5, $6, $7, $8, $9, $10, 
			st_geomfromwkb($11, $12), $13, $14, $15, $16, $17, $18, $19, $20
		)
		returning id`,
		flat.URL,
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
		s.srid,
		flat.State,
		flat.City,
		flat.District,
		flat.Street,
		flat.HouseNumber,
		flat.SSF,
		flat.IZF,
		flat.GZF,
	).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("domria: storer failed to create a flat, %v", err)
	}
	return id, nil
}

func (s *Storer) storeImages(
	flatID int,
	images []string,
	kind Kind,
	logger log.FieldLogger,
) {
	for _, url := range images {
		if err := s.storeImage(image{flatID, url, kind}); err != nil {
			s.drain.DrainNumber(metrics.FailedImageStoringNumber)
			logger.WithFields(log.Fields{"flat_id": flatID, "url": url}).Error(err)
		}
	}
}

func (s *Storer) storeImage(i image) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("domria: storer failed to begin a transaction, %v", err)
	}
	start := time.Now()
	isFound, err := s.readImage(tx, i)
	s.drain.DrainDuration(metrics.ReadingImageStoringDuration, start)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	number := metrics.UnalteredImageStoringNumber
	if !isFound {
		start := time.Now()
		err := s.createImage(tx, i)
		s.drain.DrainDuration(metrics.CreationImageStoringDuration, start)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		number = metrics.CreatedImageStoringNumber
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("domria: storer failed to commit a transaction, %v", err)
	}
	s.drain.DrainNumber(number)
	return nil
}

func (s *Storer) readImage(tx *sql.Tx, i image) (bool, error) {
	var count int
	row := tx.QueryRow(
		`select count(*) from images where flat_id = $1 and url = $2`,
		i.flatID,
		i.url,
	)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("domria: storer failed to read an image, %v", err)
	}
	return count > 0, nil
}

func (s *Storer) createImage(tx *sql.Tx, i image) error {
	_, err := tx.Exec(
		`insert into images (flat_id, url, kind) values ($1, $2, $3)`,
		i.flatID,
		i.url,
		i.kind.String(),
	)
	if err != nil {
		return fmt.Errorf("domria: storer failed to create an image, %v", err)
	}
	return nil
}
