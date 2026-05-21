package model

import "time"

type Customer struct {
	ID          string    `json:"id" db:"id"`
	Document    string    `json:"document" db:"document"`
	Name        string    `json:"name" db:"name"`
	Score       int       `json:"score" db:"score"`
	RiskLevel   string    `json:"risk_level" db:"risk_level"`
	IncomeRange string    `json:"income_range" db:"income_range"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCustomerRequest struct {
	Document    string `json:"document" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Score       int    `json:"score" binding:"required"`
	RiskLevel   string `json:"risk_level" binding:"required"`
	IncomeRange string `json:"income_range"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type DeleteCustomerRequest struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
