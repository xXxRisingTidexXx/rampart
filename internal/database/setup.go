package database

import (
	"database/sql"
	"fmt"
	gourl "net/url"
	"os"
)

func Setup(params map[string]string) (*sql.DB, error) {
	dsn := os.Getenv("RAMPART_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("database: dsn env not configured")
	}
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
