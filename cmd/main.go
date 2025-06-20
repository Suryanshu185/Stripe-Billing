package main

import (
	"log"

	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/api"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/config"
)

func main() {
	config.LoadConfig()
	// db.InitDB()
	server := api.NewServer()

	// scheduler.StartBillingScheduler()
	// scheduler.StartAutoTopUpScheduler()
	log.Fatal(server.Start())
}
