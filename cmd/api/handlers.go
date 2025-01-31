package main

import (
	"errors"
	"github.com/Shaughny/okbio-test/internal/service"
	"github.com/Shaughny/okbio-test/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// addAgent handles the POST /agents request
// It receives an IP address from the request body, validates it, and adds the agent to the database
func (app *application) addAgent(c echo.Context) error {
	var agentRequest service.AgentRequest

	// Bind request body to agentRequest struct
	err := c.Bind(&agentRequest)
	if err != nil {
		// Log and return a bad request error if binding fails
		app.logger.Errorf("Failed to bind agent request: %v", err)
		return c.JSON(http.StatusBadRequest, utils.BadRequestResponse(err))
	}

	// Validate the input (checks required fields and IP format)
	err = c.Validate(agentRequest)
	if err != nil {
		// Log and return a validation error response
		app.logger.Errorf("Failed to validate agent request: %v", err)
		return c.JSON(http.StatusBadRequest, utils.ValidationErrorResponse(err))
	}

	// Add the agent to the database
	err = app.service.AddAgent(agentRequest.IPAddress)
	if err != nil {
		// Log and return an internal server error if insertion fails
		app.logger.Errorf("Failed to add agent: %v", err)
		return c.JSON(http.StatusInternalServerError, utils.ServerErrorResponse(err))
	}

	// Return success response
	return c.JSON(http.StatusCreated, map[string]string{"message": "agent added"})
}

// getAgents handles the GET /agents request
// It retrieves all registered agents from the database
func (app *application) getAgents(c echo.Context) error {
	// Retrieve all agents from the service
	agents, err := app.service.GetAgents()
	if err != nil {
		// Log the error and check if the error is because no agents were found
		app.logger.Errorf("Failed to retrieve agents: %v", err)
		if errors.Is(err, service.ErrAgentNotFound) {
			// Return 404 Not Found if no agents exist
			return c.JSON(http.StatusNotFound, utils.NotFoundResponse(err))
		}
		// Return 500 Internal Server Error for any other failure
		return c.JSON(http.StatusInternalServerError, utils.GenericErrorResponse())
	}

	// Return the list of agents
	return c.JSON(http.StatusOK, agents)
}

// getAgent handles the GET /agents/:id request
// It retrieves details of a specific agent based on their ID
func (app *application) getAgent(c echo.Context) error {
	id := c.Param("id")

	// Convert the ID from string to integer
	intId, err := strconv.Atoi(id)
	if err != nil {
		// Log and return a bad request error if the ID is not a valid integer
		app.logger.Errorf("Invalid agent ID: %v", err)
		return c.JSON(http.StatusBadRequest, utils.BadRequestResponse(err))
	}

	// Retrieve the agent details from the database
	agent, err := app.service.GetAgent(intId)
	if err != nil {
		// Log the error and check if the agent was not found
		app.logger.Errorf("Failed to retrieve agent with ID %d: %v", intId, err)
		if errors.Is(err, service.ErrAgentNotFound) {
			// Return 404 Not Found if no such agent exists
			return c.JSON(http.StatusNotFound, utils.NotFoundResponse(err))
		}
		// Return 500 Internal Server Error for any other failure
		return c.JSON(http.StatusInternalServerError, utils.GenericErrorResponse())
	}

	// Return the agent details
	return c.JSON(http.StatusOK, agent)
}
