package main

import (
	"bytes"
	"errors"
	"github.com/Shaughny/obkio-test/internal/service"
	"github.com/Shaughny/obkio-test/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock Service
type MockService struct {
	agents    []service.Agent
	agentResp service.DetailedAgentResponse
	err       error
}

func (m *MockService) AddAgent(ipAddress string) error {
	return m.err
}

func (m *MockService) GetAgents() ([]service.Agent, error) {
	return m.agents, m.err
}

func (m *MockService) GetAgent(id int) (service.DetailedAgentResponse, error) {
	return m.agentResp, m.err
}

func getEchoInstance() *echo.Echo {
	e := echo.New()
	e.Validator = utils.NewValidator()
	e.HideBanner = true
	return e
}
func TestAddAgentHandler(t *testing.T) {
	e := getEchoInstance()
	mockService := &MockService{}
	app := &application{logger: e.Logger, service: mockService}

	t.Run("Successfully Add Agent", func(t *testing.T) {
		mockService.err = nil // No error
		requestBody := `{"ip_address": "8.8.8.8"}`

		req := httptest.NewRequest(http.MethodPost, "/agents", bytes.NewReader([]byte(requestBody)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := app.addAgent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"message": "agent added"}`, rec.Body.String())
	})

	t.Run("Invalid IP Address", func(t *testing.T) {
		requestBody := `{"ip_address": "999.999.999.999"}`

		req := httptest.NewRequest(http.MethodPost, "/agents", bytes.NewReader([]byte(requestBody)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := app.addAgent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockService.err = errors.New("DB error")
		requestBody := `{"ip_address": "8.8.8.8"}`

		req := httptest.NewRequest(http.MethodPost, "/agents", bytes.NewReader([]byte(requestBody)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := app.addAgent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetAgentsHandler(t *testing.T) {
	e := getEchoInstance()
	mockService := &MockService{
		agents: []service.Agent{
			{ID: 1, IPAddress: "8.8.8.8"},
			{ID: 2, IPAddress: "1.1.1.1"},
		},
	}
	app := &application{logger: e.Logger, service: mockService}

	t.Run("Successfully Retrieve Agents", func(t *testing.T) {
		mockService.err = nil

		req := httptest.NewRequest(http.MethodGet, "/agents", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := app.getAgents(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := `[{"id":1, "ip_address":"8.8.8.8"}, {"id":2, "ip_address":"1.1.1.1"}]`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	})

	t.Run("Database Error", func(t *testing.T) {
		mockService.err = errors.New("DB error")

		req := httptest.NewRequest(http.MethodGet, "/agents", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := app.getAgents(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetAgentHandler(t *testing.T) {
	e := getEchoInstance()
	mockService := &MockService{
		agentResp: service.DetailedAgentResponse{
			ID: 1, IPAddress: "8.8.8.8", ASN: "15169", ISP: "Google LLC",
		},
	}
	app := &application{logger: e.Logger, service: mockService}

	t.Run("Successfully Retrieve Agent", func(t *testing.T) {
		mockService.err = nil

		req := httptest.NewRequest(http.MethodGet, "/agents/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := app.getAgent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := `{"id":1, "ip_address":"8.8.8.8", "asn":"15169", "isp":"Google LLC"}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	})

	t.Run("Agent Not Found", func(t *testing.T) {
		mockService.err = service.ErrAgentNotFound

		req := httptest.NewRequest(http.MethodGet, "/agents/99", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("99")

		err := app.getAgent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/agents/abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("abc")

		err := app.getAgent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
