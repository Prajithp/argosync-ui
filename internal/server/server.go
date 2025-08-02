package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Prajithp/argosync/internal/api"
	"github.com/Prajithp/argosync/internal/config"
	"github.com/Prajithp/argosync/internal/logger"
	"github.com/Prajithp/argosync/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server represents the HTTP server
type Server struct {
	echo   *echo.Echo
	config *config.Config
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, handler *api.Handler) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Custom logger middleware using zerolog
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			logger.Info().
				Str("remote_ip", c.RealIP()).
				Str("host", req.Host).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("user_agent", req.UserAgent()).
				Int("status", res.Status).
				Dur("latency", time.Since(start)).
				Int64("bytes_in", req.ContentLength).
				Int64("bytes_out", res.Size).
				Msg("request")

			return err
		}
	})

	e.FileFS("/", "index.html", ui.IndexSubFS)
	e.StaticFS("/", ui.DistSubFS)

	apiGroup := e.Group("/api/v1")
	apiGroup.POST("/release", handler.Release)
	apiGroup.POST("/rollback", handler.Rollback)
	apiGroup.GET("/deployments", handler.GetActiveDeployments)
	apiGroup.GET("/history", handler.GetDeploymentHistory)
	apiGroup.GET("/all-deployments", handler.GetAllDeployments)

	return &Server{
		echo:   e,
		config: cfg,
	}
}

// Start starts the HTTP server
func (s *Server) Start() {
	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%d", s.config.Port)
		logger.Info().Msgf("Server starting on %s", addr)
		if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := s.echo.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Failed to gracefully shutdown server")
	}

	logger.Info().Msg("Server gracefully stopped")
}
