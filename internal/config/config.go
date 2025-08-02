package config

import (
	"flag"
	"path/filepath"
)


// Config represents the application configuration
type Config struct {
	// Server configuration
	Port int
	Debug bool

	// Database configuration
	DBType      string
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	SQLitePath  string
	MaxVersions int
}

// LoadConfig loads the configuration from command line flags
func LoadConfig() *Config {
	cfg := &Config{}

	// Parse command line flags
	flag.IntVar(&cfg.Port, "port", 8080, "Server port")
	flag.BoolVar(&cfg.Debug, "debug", false, "Enable debug mode")
	flag.StringVar(&cfg.DBType, "db-type", "sqlite", "Database type (postgres or sqlite)")
	flag.StringVar(&cfg.DBHost, "db-host", "localhost", "PostgreSQL host")
	flag.IntVar(&cfg.DBPort, "db-port", 5432, "PostgreSQL port")
	flag.StringVar(&cfg.DBUser, "db-user", "postgres", "PostgreSQL user")
	flag.StringVar(&cfg.DBPassword, "db-pass", "postgres", "PostgreSQL password")
	flag.StringVar(&cfg.DBName, "db-name", "heirloom", "PostgreSQL database name")
	flag.StringVar(&cfg.DBSSLMode, "db-sslmode", "disable", "PostgreSQL SSL mode")
	flag.StringVar(&cfg.SQLitePath, "sqlite-path", "heirloom.db", "SQLite database path")
	flag.IntVar(&cfg.MaxVersions, "max-versions", 10, "Maximum number of versions to keep per application/environment/region")
	
	flag.Parse()

	// Get absolute path for SQLite database
	if cfg.DBType == "sqlite" {
		absPath, err := filepath.Abs(cfg.SQLitePath)
		if err == nil {
			cfg.SQLitePath = absPath
		}
	}

	return cfg
}