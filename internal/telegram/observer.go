package telegram

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func NewObserver(db *sql.DB, logger log.FieldLogger) *Observer {
	return &Observer{db, logger}
}

type Observer struct {
	db     *sql.DB
	logger log.FieldLogger
}

func (o *Observer) ObserveLookups() []Lookup {
	tx, err := o.db.Begin()
	if err != nil {
		o.logger.Errorf("telegram: observer failed to begin a transaction, %v", err)
		return make([]Lookup, 0)
	}
	lookups, err := o.observeLookups(tx)
	if err != nil {
		_ = tx.Rollback()
		o.logger.Error(err)
		return make([]Lookup, 0)
	}
	if err := tx.Commit(); err != nil {
		o.logger.Errorf("telegram: observer failed to commit a transaction, %v", err)
		return make([]Lookup, 0)
	}
	return lookups
}

func (o *Observer) observeLookups(tx *sql.Tx) ([]Lookup, error) {
	rows, err := tx.Query(
		`select cte.id, cte.chat_id, cte.url, cte.uuid
		from (
			select lookups.id,
				chat_id,
				flats.url,
				uuid,
				row_number() over (partition by uuid) as count
			from lookups
				join subscriptions on subscriptions.id = lookups.subscription_id
				join flats on flats.id = lookups.flat_id
			where lookups.status = 'unseen'
		) as cte
		where count = 1`,
	)
	if err != nil {
		return nil, fmt.Errorf("telegram: observer failed to read lookups, %v", err)
	}
	lookups := make([]Lookup, 0)
	for rows.Next() {
		var lookup Lookup
		err := rows.Scan(&lookup.ID, &lookup.ChatID, &lookup.URL, &lookup.UUID)
		if err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("telegram: observer failed to scan a row, %v", err)
		}
		lookups = append(lookups, lookup)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("telegram: observer failed to finish iteration, %v", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("telegram: observer failed to close rows, %v", err)
	}
	return lookups, nil
}
