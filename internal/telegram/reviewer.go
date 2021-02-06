package telegram

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func NewReviewer(db *sql.DB, logger log.FieldLogger) *Reviewer {
	return &Reviewer{db, logger}
}

type Reviewer struct {
	db     *sql.DB
	logger log.FieldLogger
}

func (r *Reviewer) ReviewLookup(lookup Lookup) {
	if err := r.reviewLookup(lookup); err != nil {
		r.logger.WithField("id", lookup.ID).Error(err)
	}
}

func (r *Reviewer) reviewLookup(lookup Lookup) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("telegram: reviewer failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(`update lookups set status = 'seen' where id = $1`, lookup.ID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("telegram: reviewer failed to update a lookup, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("telegram: reviewer failed to commit a transaction, %v", err)
	}
	return nil
}
