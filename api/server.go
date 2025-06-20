package api

import (
	"os"

	"github.com/spf13/viper"
)

type Server struct {
	listenAddr string
}

func NewServer() *Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("OWN_PORT")
		if port == "" {
			port = "8080"
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
