package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite", "file:"+dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &Database{conn: conn}

	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("📦 SQLite database connected and initialized")
	return db, nil
}

func (db *Database) createTables() error {
	createFocusareasTable := `
	CREATE TABLE IF NOT EXISTS focusareas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		deadline TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createSessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		focusarea_id INTEGER NOT NULL,
		start_time TEXT NOT NULL,
		end_time TEXT,
		duration_minutes INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (focusarea_id) REFERENCES focusareas (id)
	);`

	if _, err := db.conn.Exec(createFocusareasTable); err != nil {
		return fmt.Errorf("failed to create focusareas table: %w", err)
	}

	if _, err := db.conn.Exec(createSessionsTable); err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}

	return nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

func (db *Database) GetConnection() *sql.DB {
	return db.conn
}
