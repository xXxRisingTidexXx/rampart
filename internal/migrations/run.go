package migrations

import (
	"database/sql"
)

// TODO: shorten the house number column varchar length.
func Run(db *sql.DB) error {
	versions, err := listVersions()
	if err != nil {
		return err
	}
	migrator, err := newMigrator(db)
	if err != nil {
		return err
	}
	latest, err := migrator.ensureVersions()
	if err != nil {
		return err
	}
	for _, version := range versions {
		if version.id > latest {
			if err = migrator.applyVersion(version); err != nil {
				return err
			}
		}
	}
	if err = migrator.commit(); err != nil {
		return err
	}
	return nil
}
