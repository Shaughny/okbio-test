package main

import (
	"github.com/Shaughny/okbio-test/config"
	"github.com/Shaughny/okbio-test/internal/service"
	"github.com/Shaughny/okbio-test/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"os"
	"strconv"
)

// Application struct contains logger and service layer.
type application struct {
	logger  echo.Logger
	service service.ServiceI
}

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Initialize database connection
	DB := config.InitializeDB()

	// Initialize application dependencies
	app := &application{
		logger: e.Logger, //
		service: &service.Service{
			DB: DB,
		},
	}

	// Middleware setup
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	rateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		rateLimit = 10
	}
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(rateLimit))))
	// Register validator
	e.Validator = utils.NewValidator()

	// Define routes
	e.GET("/agents", app.getAgents)
	e.POST("/agents", app.addAgent)
	e.GET("/agents/:id", app.getAgent)

	// Start the server
	app.logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
