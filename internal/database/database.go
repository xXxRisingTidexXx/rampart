package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	gourl "net/url"
)

func NewDatabase(dsn string, params map[string]string) (*sql.DB, error) {
	url, err := gourl.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("database: invalid dsn %s, %v", dsn, err)
	}
	values, err := gourl.ParseQuery(url.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("database: invalid dsn params %s, %v", dsn, err)
	}
	for key, value := range params {
		values.Set(key, value)
	}
	url.RawQuery = values.Encode()
	db, err := sql.Open("postgres", url.String())
	if err != nil {
		return nil, fmt.Errorf("database: failed to open the db, %v", err)
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database: failed to ping the db, %v", err)
	}
	return db, nil
}

//nolint:interfacer
func CloseDatabase(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("database: failed to close the db, %v", err)
	}
	return nil
}
