package domain

import (
	"time"

	"github.com/google/uuid"
)

type Investment struct {
	ID         uuid.UUID
	LoanID     uuid.UUID
	InvestorID string
	Amount     int64
	CreatedAt  time.Time
}

func NewInvestment(loanID uuid.UUID, investorID string, amount int64) *Investment {
	return &Investment{
		ID:         uuid.New(),
		LoanID:     loanID,
		InvestorID: investorID,
		Amount:     amount,
		CreatedAt:  time.Now(),
	}
}
