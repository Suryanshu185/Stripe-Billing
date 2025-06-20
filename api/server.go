package api

import (
	"os"

	"github.com/spf13/viper"
)

type Server struct {
	listenAddr string
}

// NewServer creates a new server instance and resolves the port to listen on.
// - On Render: it reads PORT from environment (os.Getenv("PORT"))
// - Locally: it falls back to OWN_PORT from your .env via Viper
// - Final fallback: port 8080
func NewServer() *Server {
	port := os.Getenv("PORT") // Render provides this
	if port == "" {
		port = viper.GetString("OWN_PORT") // fallback to local config
		if port == "" {
			port = "8080" // default if nothing is set
		}
	}

	return &Server{
		listenAddr: ":" + port,
	}
}

func (s *Server) Start() error {
	r := newRouter()
	return r.Run(s.listenAddr)
}
