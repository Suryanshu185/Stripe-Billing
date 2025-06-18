package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ActiveResource struct {
	gorm.Model
	Provider        string    `gorm:"size:100;not null" json:"provider"`
	InstanceID      string    `gorm:"size:100;unique;not null" json:"instanceID"`
	BilledHours     uint      `gorm:"type:int;default:0" json:"billedHours"`
	StartTime       time.Time `json:"startTime"`
	LastChargedTime time.Time `json:"lastChargedTime"`
	HourlyRateCents uint      `gorm:"type:int;not null;default:0" json:"hourlyRateCents"`
	AccountID       uint      `gorm:"not null" json:"accountID"`
	Account         Account   `gorm:"foreignKey:AccountID" json:"account"`
}

func (ar *ActiveResource) BeforeSave(tx *gorm.DB) error {
	if ar.HourlyRateCents < 1 {
		return fmt.Errorf("hourly rate must be greater than 0 cents")
	}
	return nil
}
