package main

import (
	"log"
	"os"

	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/api"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/config"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/scheduler"
)

func main() {
	// Load config
	config.LoadConfig()

	// DB connection
	db.InitDB()

	// Log which port we are binding to (for Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port :%s", port)

	// Create and start server
	server := api.NewServer()
	go scheduler.StartBillingScheduler()
	go scheduler.StartAutoTopUpScheduler()

	// Start server
	log.Fatal(server.Start())
}
