package service

import "database/sql"

// ServiceI defines the interface for the service layer
type ServiceI interface {
	// AddAgent adds a new agent or updates existing agent details.
	AddAgent(ipAddress string) error

	// GetAgents retrieves a list of all registered agents.
	GetAgents() ([]Agent, error)

	// GetAgent retrieves the detailed information of a specific agent by ID.
	GetAgent(id int) (DetailedAgentResponse, error)
}

// Service is the concrete implementation of the ServiceI interface.
// It interacts with the SQLite database.
type Service struct {
	DB *sql.DB
}
