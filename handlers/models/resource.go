package models

import "time"

type AddResourceRequest struct {
	Provider        string    `json:"provider" binding:"required"`
	InstanceID      string    `json:"instanceID" binding:"required"`
	StartTime       time.Time `json:"startTime"`
	LastChargedTime time.Time `json:"lastChargedTime"`
	HourlyRateCents uint      `json:"hourlyRateCents" binding:"required,gt=0"`
}

type ResourcesResponse struct {
	Provider        string    `json:"provider"`
	InstanceID      string    `json:"instanceID"`
	StartTime       time.Time `json:"startTime"`
	LastChargedTime time.Time `json:"lastChargedTime"`
	HourlyRateCents uint      `json:"hourlyRateCents"`
}
