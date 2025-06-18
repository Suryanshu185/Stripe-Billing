package models

import (
	"gorm.io/gorm"
	"time"
)

type BillingTransaction struct {
	gorm.Model
	AccountID  uint
	ResourceID uint
	Amount     float64
	ChargedAt  time.Time
}
