package testutils

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Estriper0/EventHub/internal/config"
	_ "github.com/lib/pq"
)

func GetDb(t *testing.T, config *config.Database) (*sql.DB, func()) {
	t.Helper()
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.DbHost,
			config.DbPort,
			config.DbUser,
			config.DbPassword,
			config.DbName,
			config.SSLMode,
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	teardown := func() {
		db.Exec("TRUNCATE event CASCADE")
		db.Close()
	}
	return db, teardown
}
