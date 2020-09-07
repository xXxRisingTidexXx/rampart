package misc

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Postgres driver along with the DB opener.
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("misc: failed to instantiate the root folder")
	}
	rootDir = filepath.Dir(filepath.Dir(filepath.Dir(filePath)))
}

var rootDir = ""

func GetEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("misc: failed to find the env %s", key)
}

func ResolvePath(path string) string {
	return filepath.Join(rootDir, path)
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("misc: failed to open the db, %v", err)
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("misc: failed to ping the db, %v", err)
	}
	return db, nil
}

func CloseDB(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("misc: failed to close the db, %v", err)
	}
	return nil
}
