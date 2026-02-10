package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/Mahaveer86619/Hearth/pkg/handlers"
	"github.com/Mahaveer86619/Hearth/pkg/services"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	router *gin.Engine
	srv    *http.Server
	addr   string
}

func NewHTTPServer(addr string) *HTTPServer {
	router := gin.Default()

	// Initialize services
	healthService := services.NewHealthService()

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(healthService)

	s := &HTTPServer{
		router: router,
		addr:   addr,
	}

	s.setupRoutes(healthHandler)

	return s
}

func (s *HTTPServer) setupRoutes(hh *handlers.HealthHandler) {
	s.router.GET("/health", hh.Check)
	s.router.GET("/ping", hh.Ping)

	// WebSocket endpoint
	s.router.GET("/ws", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "websocket endpoint placeholder"})
	})
}

// Start starts the Gin HTTP server
func (s *HTTPServer) Start() error {
	s.srv = &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop gracefully shuts down the HTTP server
func (s *HTTPServer) Stop(ctx context.Context) error {
	if s.srv != nil {
		return s.srv.Shutdown(ctx)
	}
	return nil
}

func (s *HTTPServer) Addr() string {
	return s.addr
}

func (s *HTTPServer) Name() string {
	return "HTTP Server"
}
