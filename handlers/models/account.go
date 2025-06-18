package models

type CreateAccountRequest struct {
	Email               string `json:"email" binding:"required"`
	AutoTopUpEnabled    bool   `json:"autoTopUpEnabled"`
	TopUpThresholdCents uint   `json:"topUpThresholdCents"`
	TopUpAmountCents    uint   `json:"topUpAmountCents"`
}

type AccountResponse struct {
	UserID              string `json:"userID"`
	Email               string `json:"email"`
	BalanceCents        int    `json:"balanceCents"`
	AutoTopUpEnabled    bool   `json:"autoTopUpEnabled"`
	TopUpThresholdCents uint   `json:"topUpThresholdCents"`
	TopUpAmountCents    uint   `json:"topUpAmountCents"`
	HourlyBurnRate      uint   `json:"hourlyBurnRate"`
}
