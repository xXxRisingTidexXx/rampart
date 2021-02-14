package mining

import (
	"database/sql"
	"fmt"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"time"
)

func NewStoringAmplifier(config config.StoringAmplifier, db *sql.DB) Amplifier {
	return &storingAmplifier{db, config.SRID}
}

type storingAmplifier struct {
	db   *sql.DB
	srid int
}

func (a *storingAmplifier) AmplifyFlat(flat Flat) (Flat, error) {
	ok, err := a.transactFlat(flat)
	status := "skipped"
	if err != nil {
		status = "failure"
	} else if ok {
		status = "success"
	}
	metrics.MessisStorings.WithLabelValues("flat", status).Inc()
	metrics.MessisStorings.WithLabelValues("visual", status).Add(float64(len(flat.ImageURLs)))
	return flat, err
}

func (a *storingAmplifier) transactFlat(flat Flat) (bool, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return false, fmt.Errorf("mining: amplifier failed to begin a transaction, %v", err)
	}
	ok, err := a.storeFlat(tx, flat)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("mining: amplifier failed to close a transaction, %v", err)
	}
	return ok, nil
}

func (a *storingAmplifier) storeFlat(tx *sql.Tx, flat Flat) (bool, error) {
	var count int
	now := time.Now()
	err := tx.QueryRow(`select count(*) from flats where url = $1`, flat.URL).Scan(&count)
	metrics.MessisStoringDuration.WithLabelValues("flat", "select").Observe(
		time.Since(now).Seconds(),
	)
	if err != nil {
		return false, fmt.Errorf("mining: amplifier failed to read a flat count, %v", err)
	}
	if count > 0 {
		return false, nil
	}
	var flatID int
	now = time.Now()
	err = tx.QueryRow(
		`insert into flats
		(
		 	url, price, total_area, living_area, kitchen_area, room_number, floor, total_floor,
		 	housing, point, city, street, house_number, ssf, izf, gzf
		)
		values 
		(
		    $1, $2, $3, $4, $5, $6, $7, $8, $9, st_geomfromwkb($10, $11), $12, $13, $14, $15, $16,
		 	$17
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
		wkb.Value(flat.Point),
		a.srid,
		flat.City,
		flat.Street,
		flat.HouseNumber,
		flat.SSF,
		flat.IZF,
		flat.GZF,
	).Scan(&flatID)
	metrics.MessisStoringDuration.WithLabelValues("flat", "insert").Observe(
		time.Since(now).Seconds(),
	)
	if err != nil {
		return false, fmt.Errorf("mining: amplifier failed to create a flat, %v", err)
	}
	for _, url := range flat.ImageURLs {
		var imageID int
		now := time.Now()
		err := tx.QueryRow(
			`insert into images (url)
			values ($1)
			on conflict(url) do update set url = $1
			returning id`,
			url,
		).Scan(&imageID)
		metrics.MessisStoringDuration.WithLabelValues("image", "upsert").Observe(
			time.Since(now).Seconds(),
		)
		if err != nil {
			return false, fmt.Errorf("mining: amplifier failed to upsert an image, %v", err)
		}
		now = time.Now()
		_, err = tx.Exec(`insert into visuals(flat_id, image_id) values ($1, $2)`, flatID, imageID)
		metrics.MessisStoringDuration.WithLabelValues("visual", "insert").Observe(
			time.Since(now).Seconds(),
		)
		if err != nil {
			return false, fmt.Errorf("mining: amplifier failed to create a visual, %v", err)
		}
	}
	return true, nil
}
