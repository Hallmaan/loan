package domain

import (
	"time"

	"github.com/google/uuid"
)

type LoanState string

const (
	LoanStateProposed  LoanState = "proposed"
	LoanStateApproved  LoanState = "approved"
	LoanStateInvested  LoanState = "invested"
	LoanStateDisbursed LoanState = "disbursed"
)

var ValidTransitions = map[LoanState]LoanState{
	LoanStateProposed: LoanStateApproved,
	LoanStateApproved: LoanStateInvested,
	LoanStateInvested: LoanStateDisbursed,
}

type Loan struct {
	ID                 uuid.UUID
	BorrowerID         string
	PrincipalAmount    int64
	Rate               float64
	ROI                float64
	State              LoanState
	AgreementLetterURL *string
	TotalInvested      int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func NewLoan(borrowerID string, principalAmount int64, rate, roi float64) *Loan {
	now := time.Now()
	return &Loan{
		ID:              uuid.New(),
		BorrowerID:      borrowerID,
		PrincipalAmount: principalAmount,
		Rate:            rate,
		ROI:             roi,
		State:           LoanStateProposed,
		TotalInvested:   0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (l *Loan) CanTransitionTo(newState LoanState) bool {
	expectedNext, exists := ValidTransitions[l.State]
	return exists && expectedNext == newState
}

func (l *Loan) TransitionTo(newState LoanState) error {
	if !l.CanTransitionTo(newState) {
		return ErrInvalidStateTransition
	}
	l.State = newState
	l.UpdatedAt = time.Now()
	return nil
}

func (l *Loan) RemainingAmount() int64 {
	return l.PrincipalAmount - l.TotalInvested
}

func (l *Loan) IsFullyInvested() bool {
	return l.TotalInvested >= l.PrincipalAmount
}

func (l *Loan) CanAcceptInvestment() bool {
	return l.State == LoanStateApproved
}

func (l *Loan) AddInvestment(amount int64) error {
	if !l.CanAcceptInvestment() {
		return ErrLoanNotApproved
	}
	if amount > l.RemainingAmount() {
		return ErrInvestmentExceedsLimit
	}
	l.TotalInvested += amount
	l.UpdatedAt = time.Now()
	return nil
}
