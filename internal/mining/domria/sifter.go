package domria

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"rampart/internal/config"
	"time"
)

func newSifter(db *sql.DB, config *config.Sifter) *sifter {
	return &sifter{
		db,
		config.TotalAreaBias,
		config.RoomNumberBias,
		config.FloorBias,
		config.TotalFloorBias,
		config.DistanceBias,
		time.Duration(config.UpdateTiming),
	}
}

type sifter struct {
	db             *sql.DB
	totalAreaBias  float64
	roomNumberBias int
	floorBias      int
	totalFloorBias int
	distanceBias   float64
	updateTiming   time.Duration
}

func (sifter *sifter) siftFlats(flats []*flat) ([]*flat, error) {
	length := len(flats)
	if length == 0 {
		log.Debug("domria: sifter skipped flats")
		return flats, nil
	}
	readingStmt, err := sifter.prepareReading()
	if err != nil {
		return nil, err
	}
	updateStmt, err := sifter.prepareUpdate()
	if err != nil {
		_ = readingStmt.Close()
		return nil, err
	}
	duration, newFlats := 0.0, make([]*flat, 0, length)
	for _, flat := range flats {
		start := time.Now()
		similarity, err := sifter.readFlat(readingStmt, flat)
		if err == nil {
			if similarity == nil {
				newFlats = append(newFlats, flat)
			} else if sifter.isUpdatable(similarity, flat) {
				err = sifter.updateFlat(updateStmt, similarity, flat)
			}
		}
		duration += time.Since(start).Seconds()
		if err != nil {
			log.Error(err)
		}
	}
	if err = readingStmt.Close(); err != nil {
		_ = updateStmt.Close()
		return nil, fmt.Errorf("domria: sifter failed to close the reading stmt, %v", err)
	}
	if err = updateStmt.Close(); err != nil {
		return nil, fmt.Errorf("domria: sifter failed to close the update stmt, %v", err)
	}
	log.Debugf("domria: sifter sifted %d flats (%.3fs)", len(newFlats), duration/float64(length))
	return newFlats, nil
}

func (sifter *sifter) prepareReading() (*sql.Stmt, error) {
	stmt, err := sifter.db.Prepare(
		`select id, update_time, price 
		from flats 
		where origin_url = $1 or 
		      abs(total_area - $2) <= $3 and 
		      abs(room_number - $4) <= $5 and 
		      abs(floor - $6) <= $7 and 
		      abs(total_floor - $8) <= $9 and 
		      st_distance(point, $10) <= $11`,
	)
	if err != nil {
		return nil, fmt.Errorf("domria: sifter failed to prepare the reading stmt, %v", err)
	}
	return stmt, nil
}

func (sifter *sifter) prepareUpdate() (*sql.Stmt, error) {
	stmt, err := sifter.db.Prepare(
		`update flats 
		set origin_url = $1, 
		    image_url = $2,
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
		    point = $13,
		    state = $14,
		    city = $15,
		    district = $16,
		    street = $17,
		    house_number = $18
		where id = $19`,
	)
	if err != nil {
		return nil, fmt.Errorf("domria: sifter failed to prepare the update stmt, %v", err)
	}
	return stmt, nil
}

func (sifter *sifter) readFlat(stmt *sql.Stmt, flat *flat) (*similarity, error) {
	row := stmt.QueryRow(
		flat.originURL,
		flat.totalArea,
		sifter.totalAreaBias,
		flat.roomNumber,
		sifter.roomNumberBias,
		flat.floor,
		sifter.floorBias,
		flat.totalFloor,
		sifter.totalFloorBias,
		&ewkb.Point{Point: flat.point},
		sifter.distanceBias,
	)
	var similarity similarity
	switch err := row.Scan(&similarity.id, &similarity.updateTime, &similarity.price); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &similarity, nil
	default:
		return nil, fmt.Errorf("domria: sifter failed to read a flat %s, %v", flat.originURL, err)
	}
}

func (sifter *sifter) isUpdatable(similarity *similarity, flat *flat) bool {
	isNewer := flat.updateTime.Sub(similarity.updateTime) >= sifter.updateTiming
	return isNewer || !isNewer && flat.price < similarity.price
}

func (sifter *sifter) updateFlat(stmt *sql.Stmt, similarity *similarity, flat *flat) error {
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
		similarity.id,
	)
	if err != nil {
		return fmt.Errorf("domria: sifter failed to update a flat %s, %v", flat.originURL, err)
	}
	return nil
}
