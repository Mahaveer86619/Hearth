package web

import (
	"context"
	"fmt"

	"github.com/Mahaveer86619/Hearth/pkg/config"
)

// Server defines the common interface for all server types in Hearth
type Server interface {
	Start() error
	Stop(ctx context.Context) error
	Addr() string
	Name() string
}

func NewServers() []Server {
	return []Server{
		NewHTTPServer(fmt.Sprintf(":%d", config.AppConfig.HTTPPort)),
		NewTCPServer(fmt.Sprintf(":%d", config.AppConfig.TCPPort)),
	}
}
