package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the zerolog logger
func InitLogger(debug bool) {
	// Set up the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	
	// Set the global logger
	log.Logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
	
	// Set the log level based on the debug flag
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	
	// Use a console writer for pretty output during development
	if debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}
}

// Debug logs a debug message
func Debug() *zerolog.Event {
	return log.Debug()
}

// Info logs an info message
func Info() *zerolog.Event {
	return log.Info()
}

// Warn logs a warning message
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error logs an error message
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal logs a fatal message and exits
func Fatal() *zerolog.Event {
	return log.Fatal()
}

// Panic logs a panic message and panics
func Panic() *zerolog.Event {
	return log.Panic()
}