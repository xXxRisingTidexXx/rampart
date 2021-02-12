package mining

import (
	"database/sql"
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
	return Flat{}, nil
}
