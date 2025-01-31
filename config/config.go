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

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}

	fmt.Println("Connected to SQLite database successfully at", dbPath)
	createTables()
	return DB
}

func createTables() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS agents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip_address TEXT UNIQUE NOT NULL,
		asn INTEGER,
		isp TEXT,
		last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	fmt.Println("Tables created successfully")
}
