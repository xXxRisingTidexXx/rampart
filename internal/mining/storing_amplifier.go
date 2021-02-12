package mining

import (
	"database/sql"
	"fmt"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewStoringAmplifier(config config.StoringAmplifier, db *sql.DB) Amplifier {
	return &storingAmplifier{db, config.SRID}
}

type storingAmplifier struct {
	db   *sql.DB
	srid int
}

func (a *storingAmplifier) AmplifyFlat(flat Flat) (Flat, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return flat, fmt.Errorf("mining: amplifier failed to begin a transaction, %v", err)
	}
	if err := a.storeFlat(tx, flat); err != nil {
		_ = tx.Rollback()
		return flat, err
	}
	if err := tx.Commit(); err != nil {
		return flat, fmt.Errorf("mining: amplifier failed to close a transaction, %v", err)
	}
	return flat, nil
}

func (a *storingAmplifier) storeFlat(tx *sql.Tx, flat Flat) error {
	var count int
	err := tx.QueryRow(`select count(*) from flats where url = $1`, flat.URL).Scan(&count)
	if err != nil {
		return fmt.Errorf("mining: amplifier failed to read a flat count, %v", err)
	}
	if count > 0 {
		return nil
	}
	var flatID int
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
	if err != nil {
		return fmt.Errorf("mining: amplifier failed to create a flat, %v", err)
	}
	for _, url := range flat.ImageURLs {
		var imageID int
		err := tx.QueryRow(
			`insert into images (url)
			values ($1)
			on conflict(url) do update set url = $1
			returning id`,
			url,
		).Scan(&imageID)
		if err != nil {
			return fmt.Errorf("mining: amplifier failed to upsert an image, %v", err)
		}
		_, err = tx.Exec(`insert into visuals(flat_id, image_id) values ($1, $2)`, flatID, imageID)
		if err != nil {
			return fmt.Errorf("mining: amplifier failed to create a visual, %v", err)
		}
	}
	return nil
}
