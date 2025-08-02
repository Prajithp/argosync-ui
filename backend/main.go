package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prajithp/deploy-heirloom/backend/api"
	"github.com/prajithp/deploy-heirloom/backend/database"
)

func main() {
	// Parse command line flags
	var (
		port        = flag.Int("port", 8080, "Server port")
		dbType      = flag.String("db-type", "sqlite", "Database type (postgres or sqlite)")
		dbHost      = flag.String("db-host", "localhost", "PostgreSQL host")
		dbPort      = flag.Int("db-port", 5432, "PostgreSQL port")
		dbUser      = flag.String("db-user", "postgres", "PostgreSQL user")
		dbPass      = flag.String("db-pass", "postgres", "PostgreSQL password")
		dbName      = flag.String("db-name", "heirloom", "PostgreSQL database name")
		dbSSLMode   = flag.String("db-sslmode", "disable", "PostgreSQL SSL mode")
		sqlitePath  = flag.String("sqlite-path", "heirloom.db", "SQLite database path")
		maxVersions = flag.Int("max-versions", 10, "Maximum number of versions to keep per application/environment/region")
	)
	flag.Parse()

	var db *database.DB
	var err error

	// Connect to the database based on the type
	if *dbType == "postgres" {
		// PostgreSQL configuration
		config := database.Config{
			Type:        database.PostgreSQL,
			Host:        *dbHost,
			Port:        *dbPort,
			User:        *dbUser,
			Password:    *dbPass,
			DBName:      *dbName,
			SSLMode:     *dbSSLMode,
			MaxVersions: *maxVersions,
		}

		db, err = database.New(config)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
		}

		// Migrate the database schema
		log.Println("Migrating PostgreSQL database schema...")
		err = db.MigrateSchema()
		if err != nil {
			log.Fatalf("Failed to migrate schema: %v", err)
		}

		log.Println("Database setup completed successfully")
	} else {
		// SQLite configuration
		dbPath, err := filepath.Abs(*sqlitePath)
		if err != nil {
			log.Fatalf("Failed to get absolute path to SQLite database: %v", err)
		}

		config := database.Config{
			Type:        database.SQLite,
			FilePath:    dbPath,
			MaxVersions: *maxVersions,
		}

		db, err = database.New(config)
		if err != nil {
			log.Fatalf("Failed to connect to SQLite database: %v", err)
		}

		// Check if the database file exists
		_, err = os.Stat(dbPath)
		if os.IsNotExist(err) {
			log.Println("SQLite database file does not exist, creating it...")
		}

		// Migrate the database schema
		log.Println("Migrating SQLite database schema...")
		err = db.MigrateSchema()
		if err != nil {
			log.Fatalf("Failed to migrate schema: %v", err)
		}

		log.Println("Database setup completed successfully")
	}

	defer db.Close()

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Create API handler
	h := api.NewHandler(db)

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Heirloom Deployment API")
	})

	// API routes
	e.POST("/release", h.Release)
	e.POST("/rollback", h.Rollback)
	e.GET("/deployments", h.GetActiveDeployments)
	e.GET("/history", h.GetDeploymentHistory)
	e.GET("/all-deployments", h.GetAllDeployments) // New endpoint for frontend compatibility

	// Start server
	serverAddr := fmt.Sprintf(":%d", *port)
	log.Printf("Server starting on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
