package service

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var db *sql.DB

func TestMain(m *testing.M) {

	var err error

	testDB, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("error opening test database: %v\n", err)

	}
	db = testDB
	setupTestDB(db)
	code := m.Run()
	db.Close()
	os.Remove("test.db")
	os.Exit(code)
}
func TestAddAgent(t *testing.T) {
	svc := &Service{DB: db}

	t.Run("Successfully adds agent", func(t *testing.T) {
		err := svc.AddAgent("8.8.8.8")

		assert.NoError(t, err)

		// Verify the agent was added
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM agents WHERE ip_address = ?", "8.8.8.8").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("updates existing agent", func(t *testing.T) {
		err := svc.AddAgent("8.8.8.8") // Same IP
		assert.NoError(t, err)
	})

}

func TestGetAgents(t *testing.T) {
	svc := &Service{DB: db}
	t.Run("Retrieves agents", func(t *testing.T) {
		db.Exec("INSERT INTO agents (ip_address, asn, isp) VALUES ('1.1.1.1', '15169', 'Cloudflare')")

		agents, err := svc.GetAgents()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(agents), 1) // At least one agent should exist
	})

	t.Run("Returns empty list if no agents", func(t *testing.T) {
		_, err := db.Exec("DELETE FROM agents") // Clear table
		assert.NoError(t, err)

		agents, err := svc.GetAgents()
		assert.NoError(t, err)
		assert.Len(t, agents, 0)
	})

}
func TestGetAgent(t *testing.T) {
	svc := &Service{DB: db}
	t.Run("Successfully retrieves agent", func(t *testing.T) {
		db.Exec("INSERT INTO agents (ip_address, asn, isp) VALUES ('9.9.9.9', '12345', 'Quad9')")

		var id int
		db.QueryRow("SELECT id FROM agents WHERE ip_address = ?", "9.9.9.9").Scan(&id)

		agent, err := svc.GetAgent(id)
		assert.NoError(t, err)
		assert.Equal(t, "9.9.9.9", agent.IPAddress)
	})

	t.Run("Returns error if agent not found", func(t *testing.T) {
		_, err := svc.GetAgent(999) // Non-existent ID
		assert.ErrorIs(t, err, ErrAgentNotFound)
	})
}

func setupTestDB(db *sql.DB) {
	schema := `
	CREATE TABLE agents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip_address TEXT UNIQUE NOT NULL,
		asn TEXT,
		isp TEXT,
		last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to create test schema: %v", err)
	}
}
