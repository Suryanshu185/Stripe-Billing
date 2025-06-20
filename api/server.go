package api

import "github.com/spf13/viper"

type Server struct {
	listenAddr string
}

func NewServer() *Server {
	port := ":" + viper.GetString("OWN_PORT")
	return &Server{
		listenAddr: port,
	}
}

func (s *Server) Start() error {
	r := newRouter()
	return r.Run(s.listenAddr)
}
