package migrator

import (
	"github.com/golang-migrate/migrate/v4"
)

// Migration function
func MigrateDB(dbURL string) error {
	m, err := migrate.New("file://./migrations", dbURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}
