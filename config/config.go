package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var DB *sql.DB

// InitializeDB initializes SQLite connection and creates tables if needed
func InitializeDB() *sql.DB {
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "database.sqlite" // ✅ Default for local use
	}
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}

	fmt.Println("Connected to SQLite database successfully at", dbPath)
	createTables()
	return DB
}

// createTables creates the necessary tables in the database
func createTables() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS agents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip_address varchar(15) UNIQUE NOT NULL,
		asn TEXT,
		isp TEXT,
		last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	fmt.Println("Tables created successfully")
}
