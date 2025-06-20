// package api

// import (
// 	"github.com/gin-gonic/gin"
// 	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db"
// 	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/repository"
// 	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/handlers"
// 	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/services"
// )

// func newRouter() *gin.Engine {
// 	r := gin.Default()
// 	r.GET("health-check", func(c *gin.Context) {
// 		c.JSON(200, gin.H{"message": "OK"})
// 	})
// 	registerAccountRoutes(r)
// 	registerResourceRoutes(r)
// 	return r
// }

// func registerAccountRoutes(c *gin.Engine) {
// 	dbInstance := db.DB
// 	accountRepo := repository.NewGormAccountRepo(dbInstance)
// 	activeResourceRepo := repository.NewGormActiveResourceRepo(dbInstance)
// 	accountService := services.NewAccountService(accountRepo, activeResourceRepo)
// 	accountHandler := handlers.NewAccountHandler(accountService)
// 	c.GET("/account", accountHandler.GetAccount)
// 	c.POST("/account", accountHandler.CreateAccount)
// }

//	func registerResourceRoutes(c *gin.Engine) {
//		dbInstance := db.DB
//		accountRepo := repository.NewGormAccountRepo(dbInstance)
//		activeResourceRepo := repository.NewGormActiveResourceRepo(dbInstance)
//		accountService := services.NewActiveResourceService(accountRepo, activeResourceRepo)
//		resourceHandler := handlers.NewResourceHandler(accountService)
//		c.GET("/active-resources", resourceHandler.GetUserResources)
//		c.POST("/active-resources", resourceHandler.AddResource)
//	}
package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func newRouter() *gin.Engine {
	r := gin.Default()

	r.StaticFile("/", "./cmd/index.html")

	// Health check
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	// Stripe PaymentIntent demo
	r.POST("/demo-payment", demoPaymentIntent)

	// Mocked endpoints
	registerMockAccountRoutes(r)
	registerMockResourceRoutes(r)

	return r
}

func demoPaymentIntent(c *gin.Context) {
	clientSecret := "pi_" + uuid.NewString() + "_secret_" + uuid.NewString()

	response := gin.H{
		"id":            "pi_" + uuid.NewString(),
		"object":        "payment_intent",
		"amount":        1000,
		"currency":      "usd",
		"status":        "requires_payment_method",
		"client_secret": clientSecret,
		"created":       time.Now().Unix(),
		"livemode":      false,
	}

	c.JSON(200, response)
}

func registerMockAccountRoutes(r *gin.Engine) {
	// GET /account
	r.GET("/account", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id":      1,
			"name":    "Demo User",
			"email":   "demo@example.com",
			"status":  "active",
			"credits": 100.0,
		})
	})

	// POST /account
	r.POST("/account", func(c *gin.Context) {
		c.JSON(201, gin.H{
			"message": "Mock account created successfully",
			"id":      2,
		})
	})
}

func registerMockResourceRoutes(r *gin.Engine) {
	// GET /active-resources
	r.GET("/active-resources", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"user_id": 1,
			"resources": []gin.H{
				{"type": "GPU", "usage_hours": 12, "status": "running"},
				{"type": "CPU", "usage_hours": 5, "status": "stopped"},
			},
		})
	})

	// POST /active-resources
	r.POST("/active-resources", func(c *gin.Context) {
		c.JSON(201, gin.H{
			"message":     "Mock resource added successfully",
			"resource_id": "resource_abc123",
		})
	})
}
