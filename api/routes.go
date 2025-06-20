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
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

func newRouter() *gin.Engine {
	r := gin.Default()

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
	stripe.Key = viper.GetString("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		c.JSON(500, gin.H{"error": "STRIPE_SECRET_KEY not set"})
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(1000),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"client_secret": pi.ClientSecret})
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
