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

// GetCodeFileMatchCounts returns the number of files that the code was found in.
func (d *Driver) GetCodeFileMatchCounts(codes []string) (map[string]promotion.CacheResult, error) {
	db, err := d.connect()
	if err != nil {
		return nil, fmt.Errorf("unable to connect when getting code validity with error: %w", err)
	}
	defer db.Close()

	results := make(map[string]promotion.CacheResult)

	// Ensure the query has the right number of placeholders.
	placeholders := make([]string, len(codes))
	args := make([]any, len(codes))

	for i, code := range codes {
		placeholders[i] = "?"
		args[i] = code
		results[code] = promotion.CacheResult{Found: false}
	}

	query := fmt.Sprintf(
		"SELECT code, matchcount FROM promocode WHERE code IN (%s)",
		strings.Join(placeholders, ","),
	)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("batch %v query error: %w", codes, err)
	}
	defer rows.Close()

	// Populate the results with data.
	for rows.Next() {
		var code string
		var matchCount int
		if err := rows.Scan(&code, &matchCount); err != nil {
			return nil, fmt.Errorf("scanning row error: %w", err)
		}
		results[code] = promotion.CacheResult{
			MatchCount: matchCount,
			Found:      true,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows error: %w", err)
	}

	return results, nil
}

// AddCodeFileMatchCounts caches the file match counts for the given codes.
func (d *Driver) AddCodeFileMatchCounts(codes map[string]int) error {
	db, err := d.connect()
	if err != nil {
		return fmt.Errorf("unable to add code validity with error: %w", err)
	}
	defer db.Close()

	// Start transaction for batch insert.
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed.

	// Prepare statement once, reuse for all inserts.
	stmt, err := tx.Prepare("INSERT INTO promocode(code, matchcount) VALUES(?, ?) ON CONFLICT(code) DO NOTHING")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert all codes.
	for code, matchCount := range codes {
		_, err := stmt.Exec(code, matchCount)
		if err != nil {
			return fmt.Errorf("failed to insert code %s: %w", code, err)
		}
	}

	// Commit transaction.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
