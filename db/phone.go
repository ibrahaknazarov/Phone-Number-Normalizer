// Package db handles all database operations for the phone number normalizer.
package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Phone represents a phone number record in the phone_numbers table.
type Phone struct {
	ID     int
	Number string
}

// Open establishes a connection to the database and returns a DB instance.
func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// DB wraps the underlying sql.DB connection.
type DB struct {
	db *sql.DB
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.db.Close()
}

// Seed populates the database with test phone numbers in various formats.
func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}
	return nil
}

// insertPhone adds a new phone number to the database.
func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// AllPhones retrieves all phone numbers from the database.
func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

// FindPhone looks up a phone number by value.
// Returns nil if not found, error on database failure.
func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	row := db.db.QueryRow("SELECT id, value FROM phone_numbers WHERE value=$1", number)
	err := row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

// UpdatePhone updates an existing phone number record.
func (db *DB) UpdatePhone(p *Phone) error {
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.db.Exec(statement, p.ID, p.Number)
	return err
}

// DeletePhone removes a phone number record by ID.
func (db *DB) DeletePhone(id int) error {
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.db.Exec(statement, id)
	return err
}

// Migrate creates the phone_numbers table if it doesn't exist.
func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

// createPhoneNumbersTable creates the schema for storing phone numbers.
func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
    CREATE TABLE IF NOT EXISTS phone_numbers (
      id SERIAL,
      value VARCHAR(255)
    )`
	_, err := db.Exec(statement)
	return err
}

// Reset drops and recreates a database by name.
func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

// resetDB drops an existing database and creates a new one.
func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

// createDB creates a new database.
func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

// getPhone retrieves a phone number by ID (currently unused).
func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT id, value FROM phone_numbers WHERE id=$1", id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}
