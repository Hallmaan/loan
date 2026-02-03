package domain

import "errors"

var (
	ErrLoanNotFound          = errors.New("loan not found")
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrInvestmentExceedsLimit = errors.New("investment amount exceeds remaining principal")
	ErrLoanNotApproved       = errors.New("loan must be in approved state to accept investments")
	ErrLoanNotInvested       = errors.New("loan must be in invested state to disburse")
	ErrLoanAlreadyApproved   = errors.New("loan is already approved")
	ErrLoanAlreadyDisbursed  = errors.New("loan is already disbursed")
	ErrInvalidAmount         = errors.New("invalid amount")
	ErrApprovalNotFound      = errors.New("approval not found")
	ErrDisbursementNotFound  = errors.New("disbursement not found")
)
