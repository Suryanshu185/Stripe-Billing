package api

import (
	"os"

	"github.com/gin-gonic/gin"
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
	r := gin.Default()

	// Add a basic health check route
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.Run(s.listenAddr)
}
