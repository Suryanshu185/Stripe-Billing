package repository

import (
	"fmt"
	dbmodels "gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/models"
	"gorm.io/gorm"
)

type ActiveResourceRepository interface {
	Create(account *dbmodels.ActiveResource) error
	Delete(account *dbmodels.ActiveResource) error
	Update(account *dbmodels.ActiveResource) error
	GetByID(accountID uint) (*dbmodels.ActiveResource, error)
	GetAll() ([]*dbmodels.ActiveResource, error)
	GetByAccountID(accountID uint) ([]*dbmodels.ActiveResource, error)
}

type GormActiveResourceRepo struct {
	DB *gorm.DB
}

func NewGormActiveResourceRepo(db *gorm.DB) ActiveResourceRepository {
	return &GormActiveResourceRepo{DB: db}
}

func (r *GormActiveResourceRepo) Create(account *dbmodels.ActiveResource) error {
	if err := r.DB.Create(account).Error; err != nil {
		return GetDBError(err)
	}
	return nil
}

func (r *GormActiveResourceRepo) Delete(account *dbmodels.ActiveResource) error {
	if err := r.DB.Delete(account).Error; err != nil {
		return fmt.Errorf("failed to delete active resource entry : %w", err)
	}
	return nil
}

func (r *GormActiveResourceRepo) Update(account *dbmodels.ActiveResource) error {
	if err := r.DB.Save(account).Error; err != nil {
		return fmt.Errorf("failed to update active resource entry : %w", err)
	}
	return nil
}

func (r *GormActiveResourceRepo) GetByID(id uint) (*dbmodels.ActiveResource, error) {
	var account dbmodels.ActiveResource
	if err := r.DB.First(&account, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get active resource by id : %w", err)
	}
	return &account, nil
}

func (r *GormActiveResourceRepo) GetAll() ([]*dbmodels.ActiveResource, error) {
	var accounts []*dbmodels.ActiveResource
	if err := r.DB.Find(&accounts).Error; err != nil {
		return nil, GetDBError(err)
	}
	return accounts, nil
}

func (r *GormActiveResourceRepo) GetByAccountID(accountID uint) ([]*dbmodels.ActiveResource, error) {
	var accounts []*dbmodels.ActiveResource
	if err := r.DB.Where("account_id = ?", accountID).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get active resource by acount ID: %w", err)
	}
	return accounts, nil
}
