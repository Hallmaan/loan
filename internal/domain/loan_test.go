package domain

import (
	"testing"
)

func TestNewLoan(t *testing.T) {
	loan := NewLoan("borrower-123", 1000000, 0.15, 0.12)

	if loan.BorrowerID != "borrower-123" {
		t.Errorf("expected borrower_id to be 'borrower-123', got '%s'", loan.BorrowerID)
	}
	if loan.PrincipalAmount != 1000000 {
		t.Errorf("expected principal_amount to be 1000000, got %d", loan.PrincipalAmount)
	}
	if loan.State != LoanStateProposed {
		t.Errorf("expected state to be 'proposed', got '%s'", loan.State)
	}
	if loan.TotalInvested != 0 {
		t.Errorf("expected total_invested to be 0, got %d", loan.TotalInvested)
	}
}

func TestLoanStateTransitions(t *testing.T) {
	tests := []struct {
		name        string
		fromState   LoanState
		toState     LoanState
		shouldError bool
	}{
		{"proposed to approved", LoanStateProposed, LoanStateApproved, false},
		{"approved to invested", LoanStateApproved, LoanStateInvested, false},
		{"invested to disbursed", LoanStateInvested, LoanStateDisbursed, false},
		{"proposed to invested (invalid)", LoanStateProposed, LoanStateInvested, true},
		{"proposed to disbursed (invalid)", LoanStateProposed, LoanStateDisbursed, true},
		{"approved to disbursed (invalid)", LoanStateApproved, LoanStateDisbursed, true},
		{"disbursed to anything (invalid)", LoanStateDisbursed, LoanStateProposed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loan := &Loan{State: tt.fromState}
			err := loan.TransitionTo(tt.toState)

			if tt.shouldError && err == nil {
				t.Errorf("expected error for transition %s -> %s", tt.fromState, tt.toState)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error for transition %s -> %s: %v", tt.fromState, tt.toState, err)
			}
		})
	}
}

func TestLoanCanAcceptInvestment(t *testing.T) {
	tests := []struct {
		state    LoanState
		expected bool
	}{
		{LoanStateProposed, false},
		{LoanStateApproved, true},
		{LoanStateInvested, false},
		{LoanStateDisbursed, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			loan := &Loan{State: tt.state}
			if loan.CanAcceptInvestment() != tt.expected {
				t.Errorf("expected CanAcceptInvestment() to be %v for state %s", tt.expected, tt.state)
			}
		})
	}
}

func TestAddInvestment(t *testing.T) {
	loan := &Loan{
		State:           LoanStateApproved,
		PrincipalAmount: 1000000,
		TotalInvested:   0,
	}

	// Add first investment
	err := loan.AddInvestment(500000)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if loan.TotalInvested != 500000 {
		t.Errorf("expected total_invested to be 500000, got %d", loan.TotalInvested)
	}

	// Add second investment
	err = loan.AddInvestment(300000)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if loan.TotalInvested != 800000 {
		t.Errorf("expected total_invested to be 800000, got %d", loan.TotalInvested)
	}

	// Try to exceed limit
	err = loan.AddInvestment(300000)
	if err != ErrInvestmentExceedsLimit {
		t.Errorf("expected ErrInvestmentExceedsLimit, got %v", err)
	}
}

func TestAddInvestmentWrongState(t *testing.T) {
	loan := &Loan{
		State:           LoanStateProposed,
		PrincipalAmount: 1000000,
		TotalInvested:   0,
	}

	err := loan.AddInvestment(500000)
	if err != ErrLoanNotApproved {
		t.Errorf("expected ErrLoanNotApproved, got %v", err)
	}
}

func TestIsFullyInvested(t *testing.T) {
	loan := &Loan{
		PrincipalAmount: 1000000,
		TotalInvested:   500000,
	}

	if loan.IsFullyInvested() {
		t.Error("expected IsFullyInvested() to be false")
	}

	loan.TotalInvested = 1000000
	if !loan.IsFullyInvested() {
		t.Error("expected IsFullyInvested() to be true")
	}
}

func TestRemainingAmount(t *testing.T) {
	loan := &Loan{
		PrincipalAmount: 1000000,
		TotalInvested:   300000,
	}

	if loan.RemainingAmount() != 700000 {
		t.Errorf("expected remaining amount to be 700000, got %d", loan.RemainingAmount())
	}
}
