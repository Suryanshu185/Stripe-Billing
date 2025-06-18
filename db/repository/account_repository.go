package repository

import (
	"fmt"
	dbmodels "gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/models"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *dbmodels.Account) error
	Update(account *dbmodels.Account) error
	GetByID(id uint) (*dbmodels.Account, error)
	GetByStripeCustomerID(stripeCustomerID string) (*dbmodels.Account, error)
	GetByUserID(userID string) (*dbmodels.Account, error)
	GetAutoTopUpAccounts() ([]*dbmodels.Account, error)
}

type GormAccountRepo struct {
	DB *gorm.DB
}

func NewGormAccountRepo(db *gorm.DB) AccountRepository {
	return &GormAccountRepo{DB: db}
}

func (r *GormAccountRepo) Create(account *dbmodels.Account) error {
	if err := r.DB.Create(account).Error; err != nil {
		return fmt.Errorf("failed to create account : %w", err)
	}
	return nil
}

func (r *GormAccountRepo) Update(account *dbmodels.Account) error {
	if err := r.DB.Save(account).Error; err != nil {
		return fmt.Errorf("failed to save account : %w", err)
	}
	return nil
}

func (r *GormAccountRepo) GetByID(id uint) (*dbmodels.Account, error) {
	var account dbmodels.Account
	if err := r.DB.First(&account, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get account by ID: %w", err)
	}
	return &account, nil
}

func (r *GormAccountRepo) GetByStripeCustomerID(stripeCustomerID string) (*dbmodels.Account, error) {
	var account dbmodels.Account
	if err := r.DB.Where("stripe_customer_id = ?", stripeCustomerID).First(&account).Error; err != nil {
		return nil, fmt.Errorf("failed to get account by Stripe customer ID: %w", err)
	}
	return &account, nil
}

func (r *GormAccountRepo) GetByUserID(userID string) (*dbmodels.Account, error) {
	var account dbmodels.Account
	if err := r.DB.Where("user_id = ?", userID).First(&account).Error; err != nil {
		return nil, GetDBError(err)
	}
	return &account, nil

}

// GetAutoTopUpAccounts returns the accounts with auto_top_up_enabled set to true
func (r *GormAccountRepo) GetAutoTopUpAccounts() ([]*dbmodels.Account, error) {
	var accounts []*dbmodels.Account
	if err := r.DB.Where("auto_top_up_enabled = ?", true).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get accounts for auto_top_up_enabled true")
	}
	return accounts, nil
}
