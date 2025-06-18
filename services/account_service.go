package services

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	dbmodels "gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/models"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/repository"
	apimodels "gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/handlers/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type AccountService struct {
	accountRepo        repository.AccountRepository
	activeResourceRepo repository.ActiveResourceRepository
	logger             *zap.Logger
}

func NewAccountService(accountRepo repository.AccountRepository, activeResourceRepo repository.ActiveResourceRepository) *AccountService {
	logger := viper.Get("Logger").(*zap.Logger)
	return &AccountService{
		accountRepo:        accountRepo,
		activeResourceRepo: activeResourceRepo,
		logger:             logger,
	}
}

//TODO [Vikas] : do we need to mandate autoTopUp ?

func (s *AccountService) CreateAccount(userID, emailID string, topUpThresholdCents, topUpAmountCents uint) error {
	if topUpThresholdCents < 5 {
		return fmt.Errorf("TopUpThreshold can't be less than 5")
	}

	if topUpAmountCents < topUpThresholdCents {
		return fmt.Errorf("TopUpAmount can't be less than TopUpThreshold amount")
	}

	//TODO : create account in strip and get customer id
	newAccount := dbmodels.Account{
		UserID:              userID,
		StripeCustomerID:    "",
		Email:               emailID,
		BalanceCents:        0,
		AutoTopUpEnabled:    false,
		TopUpThresholdCents: topUpThresholdCents,
		TopUpAmountCents:    topUpAmountCents,
	}
	return s.accountRepo.Create(&newAccount)
}

func (s *AccountService) GetUserAccount(userID string) (*apimodels.AccountResponse, error) {
	account, err := s.accountRepo.GetByUserID(userID)

	if err != nil {
		s.logger.Error("Failed to get account", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}

	s.logger.Debug("GetUserAccount", zap.Any("account", account))

	hourlyBurnRate, err := s.getHourlyBurnRate(account.ID)

	if err != nil {
		return nil, err
	}

	s.logger.Debug("GetUserAccount", zap.Any("hourlyBurnRate", hourlyBurnRate))

	response := &apimodels.AccountResponse{
		UserID:              account.UserID,
		Email:               account.Email,
		BalanceCents:        account.BalanceCents,
		AutoTopUpEnabled:    account.AutoTopUpEnabled,
		TopUpThresholdCents: account.TopUpThresholdCents,
		TopUpAmountCents:    account.TopUpAmountCents,
		HourlyBurnRate:      hourlyBurnRate,
	}
	return response, nil
}

func (s *AccountService) UpdateAutoTopUp(userID string, TopUpThresholdCents, TopUpAmountCents uint) error {
	if TopUpThresholdCents < 500 {
		return fmt.Errorf("topUpThreshold can't be less than 5 dollers")
	}

	if TopUpAmountCents < TopUpThresholdCents {
		return fmt.Errorf("topUpAmount can't be less than TopUpThreshold amount")
	}

	account, err := s.accountRepo.GetByUserID(userID)

	if err != nil {
		return err
	}

	account.AutoTopUpEnabled = true
	account.TopUpThresholdCents = TopUpThresholdCents
	account.TopUpAmountCents = TopUpAmountCents

	return s.accountRepo.Create(account)
}

func (s *AccountService) getHourlyBurnRate(accountID uint) (uint, error) {
	activeResources, err := s.activeResourceRepo.GetByAccountID(accountID)
	if err != nil {
		return 0, err
	}
	totalCentsPerHour := uint(0)
	for _, activeResource := range activeResources {
		totalCentsPerHour += activeResource.HourlyRateCents
	}
	return totalCentsPerHour, nil
}

func (s *AccountService) CalculateUsageAndDebitBalance() error {
	s.logger.Debug("starting to calculate usage and debit balance")
	resources, err := s.activeResourceRepo.GetAll()
	if err != nil {
		s.logger.Error("failed to get all resources", zap.Error(err))
		return err
		//TODO : log error here
	}
	s.logger.Debug("got all resources", zap.Int("count", len(resources)))

	for _, resource := range resources {
		lastChargeDuration := time.Since(resource.LastChargedTime)
		hoursSinceLastCharge := uint(lastChargeDuration.Hours())

		if hoursSinceLastCharge >= 1 {
			s.logger.Debug("resource usage", zap.Uint("resource_id", resource.ID), zap.Uint("hours_used", hoursSinceLastCharge))
			// Calculate the amount to charge
			amountToChargeCents := int(hoursSinceLastCharge * resource.HourlyRateCents)

			s.logger.Info("debiting account", zap.Uint("resource_id", resource.ID), zap.Int("amount_cents", amountToChargeCents))
			// Deduct from user's balance
			account, err := s.accountRepo.GetByID(resource.AccountID)
			if err != nil {
				s.logger.Error("failed to get account by ID", zap.Uint("account_id", resource.AccountID), zap.Error(err))
				continue
			}
			account.BalanceCents -= amountToChargeCents

			err = s.accountRepo.Update(account)
			if err != nil {
				s.logger.Error("failed to update account", zap.Uint("account_id", resource.AccountID), zap.Error(err))
				continue
			}

			// Update resource's LastChargedTime
			resource.LastChargedTime = resource.LastChargedTime.Add(time.Duration(hoursSinceLastCharge) * time.Hour)

			err = s.activeResourceRepo.Update(resource)
			if err != nil {
				s.logger.Error("failed to update resource", zap.Uint("resource_id", resource.ID), zap.Error(err))
				continue
			}
		}
	}
	s.logger.Debug("Finished calculating usage and debiting balance")
	return nil
}

func (s *AccountService) CheckAndTopUpBalance() error {
	allAutoTopUpAccounts, err := s.accountRepo.GetAutoTopUpAccounts()
	s.logger.Debug("total accounts for auto_top_up_enabled true", zap.Int("count", len(allAutoTopUpAccounts)))
	if err != nil {
		s.logger.Error("failed to get accounts for auto_top_up_enabled true", zap.Error(err))
		return err
	}
	for _, account := range allAutoTopUpAccounts {
		s.logger.Debug("Checking account for top-up", zap.Uint("account_id", account.ID))
		if account.BalanceCents <= int(account.TopUpThresholdCents) {
			s.logger.Info("topping up account", zap.Uint("account_id", account.ID), zap.Uint("top_up_amount_cents", account.TopUpAmountCents))
			// Update user balance
			account.BalanceCents += int(account.TopUpAmountCents)
			err := s.accountRepo.Update(account)
			if err != nil {
				s.logger.Error("failed to update account for auto top-up", zap.Uint("account_id", account.ID), zap.Error(err))
				continue
			}
			// Record the transaction. Not implemented yet
			//transaction := models.Transaction{
			//	UserID:      user.ID,
			//	Type:        models.AutoTopUp,
			//	Amount:      user.TopUpAmount,
			//	Description: "Auto top-up",
			//}
			//utils.DB.Create(&transaction)
		}
	}
	s.logger.Debug("finished auto top-up for all the accounts")
	return nil
}
