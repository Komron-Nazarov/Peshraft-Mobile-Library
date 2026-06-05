package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB(connString string) (*DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	return &DB{conn: db}, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) GetConn() *sql.DB {
	return d.conn
}

// В файле internal/repositories/db.go добавь:
//
//	func (d *DB) Exec(query string, args ...interface{}) error {
//		_, err := d.conn.Exec(query, args...)
//		return err
//	}
//
// В файле internal/repositories/db.go
func (d *DB) Exec(query string, args ...interface{}) error {
	// Используем внутреннее соединение sql.DB
	_, err := d.conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}
	return nil
}
