package scheduler

import (
	"github.com/spf13/viper"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/repository"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/services"
	"go.uber.org/zap"
	"time"
)

func StartBillingScheduler() {
	logger := viper.Get("Logger").(*zap.Logger)
	dbInstance := db.DB
	accountRepo := repository.NewGormAccountRepo(dbInstance)
	activeResourceRepo := repository.NewGormActiveResourceRepo(dbInstance)
	accountService := services.NewAccountService(accountRepo, activeResourceRepo)

	ticker := time.NewTicker(5 * time.Second)
	logger.Info("starting billing scheduler")
	go func() {
		for {
			select {
			case <-ticker.C:
				err := accountService.CalculateUsageAndDebitBalance()
				if err != nil {
					logger.Error("failed to calculate usage and debit balance", zap.Error(err))
				}
			}
		}
	}()
}

func StartAutoTopUpScheduler() {
	logger := viper.Get("Logger").(*zap.Logger)
	dbInstance := db.DB
	accountRepo := repository.NewGormAccountRepo(dbInstance)
	activeResourceRepo := repository.NewGormActiveResourceRepo(dbInstance)
	accountService := services.NewAccountService(accountRepo, activeResourceRepo)

	//ticker := time.NewTicker(1 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	logger.Info("starting auto top-up scheduler")
	go func() {
		for {
			select {
			case <-ticker.C:
				err := accountService.CheckAndTopUpBalance()
				if err != nil {
					logger.Error("Failed to check and top-up balance", zap.Error(err))
				}
			}
		}
	}()
}
