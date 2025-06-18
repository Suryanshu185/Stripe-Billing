package models

import (
	"fmt"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID              string `gorm:"size:100;unique_index;not null" json:"userID"`
	StripeCustomerID    string `gorm:"size:100;unique_index;not null" json:"stripeCustomerID"`
	Email               string `gorm:"size:100;unique_index;not null" json:"email"`
	BalanceCents        int    `gorm:"type:int;default:0" json:"balanceCents"`
	AutoTopUpEnabled    bool   `gorm:"default:false" json:"autoTopUpEnabled"`
	TopUpThresholdCents uint   `gorm:"type:int;default:0" json:"topUpThresholdCents"`
	TopUpAmountCents    uint   `gorm:"type:int;default:0" json:"topUpAmountCents"`
}

func (a *Account) BeforeSave(tx *gorm.DB) error {
	if a.TopUpAmountCents < a.TopUpThresholdCents {
		return fmt.Errorf("topUpAmount can't be less than topUpThreshold amount")
	}
	return nil
}
