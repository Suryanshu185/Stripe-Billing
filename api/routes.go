package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/repository"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/handlers"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/services"
)

func newRouter() *gin.Engine {
	r := gin.Default()
	r.GET("health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})
	registerAccountRoutes(r)
	registerResourceRoutes(r)
	return r
}

func registerAccountRoutes(c *gin.Engine) {
	dbInstance := db.DB
	accountRepo := repository.NewGormAccountRepo(dbInstance)
	activeResourceRepo := repository.NewGormActiveResourceRepo(dbInstance)
	accountService := services.NewAccountService(accountRepo, activeResourceRepo)
	accountHandler := handlers.NewAccountHandler(accountService)
	c.GET("/account", accountHandler.GetAccount)
	c.POST("/account", accountHandler.CreateAccount)
}

func registerResourceRoutes(c *gin.Engine) {
	dbInstance := db.DB
	accountRepo := repository.NewGormAccountRepo(dbInstance)
	activeResourceRepo := repository.NewGormActiveResourceRepo(dbInstance)
	accountService := services.NewActiveResourceService(accountRepo, activeResourceRepo)
	resourceHandler := handlers.NewResourceHandler(accountService)
	c.GET("/active-resources", resourceHandler.GetUserResources)
	c.POST("/active-resources", resourceHandler.AddResource)
}
