package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Agent represents a minimal agent model with just ID and IP address.
type Agent struct {
	ID        int    `json:"id"`
	IPAddress string `json:"ip_address"`
}

// AgentRequest defines the structure for incoming agent registration requests.
type AgentRequest struct {
	IPAddress string `json:"ip_address" validate:"required,ip"` // Ensures IP format validation
}

// DetailedAgentResponse provides full details about an agent, including ASN and ISP.
type DetailedAgentResponse struct {
	ID        int    `json:"id"`
	IPAddress string `json:"ip_address"`
	ASN       string `json:"asn"`
	ISP       string `json:"isp"`
}

// DetailedAgentRequest is used internally when fetching IP details from an external API.
type DetailedAgentRequest struct {
	IPAddress string `json:"ip_address" validate:"required"`
	ASN       string `json:"as" validate:"required"`
	ISP       string `json:"isp" validate:"required"`
}

// ErrAgentNotFound is returned when an agent is not found in the database.
var ErrAgentNotFound = errors.New("agent not found")

// AddAgent inserts an agent into the database or updates its details if it already exists.
func (s *Service) AddAgent(ipAddress string) error {
	// Fetch IP details from an external API
	agent, err := getIPInformation(ipAddress)
	if err != nil {
		return fmt.Errorf("error getting IP information: %w", err)
	}

	// SQL query to insert or update the agent record in the database
	query := `
	INSERT INTO agents (ip_address, asn, isp, last_updated) 
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
	ON CONFLICT(ip_address) DO UPDATE 
	SET asn = EXCLUDED.asn,
	    isp = EXCLUDED.isp,
	    last_updated = CURRENT_TIMESTAMP;
	`

	// Execute the query
	_, err = s.DB.Exec(query, agent.IPAddress, agent.ASN, agent.ISP)
	if err != nil {
		return fmt.Errorf("error inserting agent: %w", err)
	}

	return nil
}

// GetAgents retrieves a list of all registered agents from the database.
func (s *Service) GetAgents() ([]Agent, error) {
	query := "SELECT id, ip_address FROM agents"
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Iterate over the result set
	agents := []Agent{}
	for rows.Next() {
		var agent Agent
		err = rows.Scan(&agent.ID, &agent.IPAddress)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		agents = append(agents, agent)
	}

	// Check if there were any errors during iteration
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return agents, nil
}

// GetAgent retrieves an agent's detailed information from the database based on the given ID.
func (s *Service) GetAgent(ID int) (DetailedAgentResponse, error) {
	query := "SELECT id, ip_address, asn, isp FROM agents WHERE id = $1"
	var agent DetailedAgentResponse

	// Execute the query and scan the result into the agent struct
	err := s.DB.QueryRow(query, ID).Scan(&agent.ID, &agent.IPAddress, &agent.ASN, &agent.ISP)
	if err != nil {
		// Return a custom error if no rows were found
		if errors.Is(err, sql.ErrNoRows) {
			return agent, ErrAgentNotFound
		}
		return agent, fmt.Errorf("error fetching agent from database: %w", err)
	}

	return agent, nil
}

// getIPInformation fetches ASN and ISP details for a given IP address from an external API.
func getIPInformation(ip string) (DetailedAgentRequest, error) {
	var agent DetailedAgentRequest
	url := "http://ip-api.com/json/"

	// Call the external API to get IP details
	resp, err := http.Get(url + ip)
	if err != nil {
		return agent, fmt.Errorf("error fetching IP information: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return agent, fmt.Errorf("error reading response body: %w", err)
	}

	// Unmarshal JSON response into the agent struct
	err = json.Unmarshal(body, &agent)
	if err != nil {
		return agent, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	// Ensure the API returned valid data
	if agent.ASN == "" || agent.ISP == "" {
		return agent, fmt.Errorf("IP API information incomplete for IP: %s", ip)
	}

	// Set the IP address and normalize ASN field
	agent.IPAddress = ip
	agent.ASN = strings.Split(agent.ASN, " ")[0]

	return agent, nil
}
