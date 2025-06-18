package services

import (
	"github.com/spf13/viper"
	dbModels "gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/models"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/repository"
	"go.uber.org/zap"
	"math"
	"time"
)

type ActiveResourceService struct {
	activeResourceRepository repository.ActiveResourceRepository
	accountRepository        repository.AccountRepository
	logger                   *zap.Logger
}

func NewActiveResourceService(accountRepo repository.AccountRepository, activeResourceRepo repository.ActiveResourceRepository) *ActiveResourceService {
	logger := viper.Get("Logger").(*zap.Logger)
	return &ActiveResourceService{
		activeResourceRepository: activeResourceRepo,
		accountRepository:        accountRepo,
		logger:                   logger,
	}
}

func (s *ActiveResourceService) AddResource(userID, instanceID, provider string, hourlyRateCents uint) error {
	account, err := s.accountRepository.GetByUserID(userID)
	if err != nil {
		return err
	}

	resource := dbModels.ActiveResource{
		AccountID:       account.ID,
		Provider:        provider,
		InstanceID:      instanceID,
		HourlyRateCents: hourlyRateCents,
		StartTime:       time.Now(),
		LastChargedTime: time.Now(),
	}
	if err := s.activeResourceRepository.Create(&resource); err != nil {
		return err
	}

	// Immediately charge for the first hour
	account.BalanceCents -= int(hourlyRateCents) // Charge for one hour

	if err := s.accountRepository.Update(account); err != nil {
		return err
	}

	// TODO : Create a billing transaction
	//transaction := BillingTransaction{
	//	UserID:     user.ID,
	//	ResourceID: resource.ID,
	//	Amount:     hourlyRate,
	//	ChargedAt:  time.Now(),
	//}
	//db.Create(&transaction)

	return nil
}

func (s *ActiveResourceService) TerminateResource(resourceID uint) error {
	resource, err := s.activeResourceRepository.GetByID(resourceID)
	if err != nil {
		return err
	}
	// Calculate total hours used
	totalDuration := time.Since(resource.StartTime)
	totalHoursUsed := uint(math.Ceil(totalDuration.Hours()))

	// Calculate hours charged so far
	chargedDuration := resource.LastChargedTime.Sub(resource.StartTime)
	hoursCharged := uint(chargedDuration.Hours())

	// If total hours used > hours charged, charge the difference
	if totalHoursUsed > hoursCharged {
		hoursToCharge := totalHoursUsed - hoursCharged
		amountToCharge := hoursToCharge * resource.HourlyRateCents

		account, err := s.accountRepository.GetByID(resource.AccountID)
		if err != nil {
			return err
		}
		account.BalanceCents -= int(amountToCharge)
		err = s.accountRepository.Update(account)
		if err != nil {
			return err
		}

		// Update resource's LastChargedTime
		resource.LastChargedTime = resource.StartTime.Add(time.Duration(totalHoursUsed) * time.Hour)
		err = s.activeResourceRepository.Update(resource)
		if err != nil {
			return err
		}

		//TODO do not delete below code
		// Create a billing transaction
		//transaction := BillingTransaction{
		//	UserID:     user.ID,
		//	ResourceID: resource.ID,
		//	Amount:     amountToCharge,
		//	ChargedAt:  time.Now(),
		//}
		//db.Create(&transaction)
	}

	// Delete or deactivate the resource
	err = s.activeResourceRepository.Delete(resource)
	if err != nil {
		return err
	}
	return nil
}

func (s *ActiveResourceService) GetUserResource(userID string) ([]*dbModels.ActiveResource, error) {
	account, err := s.accountRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	activeResources, err := s.activeResourceRepository.GetByAccountID(account.ID)
	if err != nil {
		return nil, err
	}
	return activeResources, nil
}
