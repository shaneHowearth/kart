package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shanehowearth/kart/promotion"
)

type Driver struct{}

var _ promotion.Store = (*Driver)(nil)

func (*Driver) connect() (*sql.DB, error) {
	dbName := "promotion_data.db"

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Lock the connection options down.
	db.SetMaxOpenConns(1) // SQLite works best with fewer open connections
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping the database to verify the connection and the driver are working
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting/pinging database: %w", err)
	}

	return db, nil
}

// InitialiseDataStore creates the table needed for the queries to work.
func (d *Driver) InitialiseDataStore() error {
	db, err := d.connect()
	if err != nil {
		return fmt.Errorf("unable to connect when initialising datastore with error: %w", err)
	}
	defer db.Close()

	query := `
	CREATE TABLE IF NOT EXISTS promocode (
		code TEXT PRIMARY KEY NOT NULL,
		matchcount INTEGER NOT NULL
	);`
	_, err = db.Exec(query)

	return err
}

// GetCodeFileMatchCount returns the number of files that the code was found in.
func (d *Driver) GetCodeFileMatchCount(code string) (int, error) {
	var matchCount int

	db, err := d.connect()
	if err != nil {
		return 0, fmt.Errorf("unable to connect when getting code validity with error: %w", err)
	}
	defer db.Close()

	query := "SELECT matchcount FROM promocode WHERE code = ?"

	row := db.QueryRow(query, code)
	if err := row.Scan(&matchCount); err != nil {
		if err == sql.ErrNoRows {
			return 0, err // Let the caller know that there were no rows found.
		}

		return 0, fmt.Errorf("getting code validity for %s database query error: %w", code, err)
	}

	return matchCount, nil
}

// AddCodeFileMatchCount caches the file match count for the given code.
func (d *Driver) AddCodeFileMatchCount(code string, matchCount int) error {
	db, err := d.connect()
	if err != nil {
		return fmt.Errorf("unable to add code validity with error: %w", err)
	}
	defer db.Close()

	const insertSQL = "INSERT INTO promocode(code, matchcount) VALUES(?, ?)"

	_, err = db.Exec(insertSQL, code, matchCount)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil // Not an error, just skip it
		}
		return err // Real error
	}

	return nil
}
