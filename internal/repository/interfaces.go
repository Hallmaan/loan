package repository

import (
	"context"

	"github.com/agunghallmanmaliki/amartha/internal/domain"
	"github.com/google/uuid"
)

type LoanRepository interface {
	Create(ctx context.Context, loan *domain.Loan) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Loan, error)
	Update(ctx context.Context, loan *domain.Loan) error
	List(ctx context.Context, filter LoanFilter) ([]*domain.Loan, int64, error)
}

type LoanFilter struct {
	State    *domain.LoanState
	Limit    int
	Offset   int
}

type ApprovalRepository interface {
	Create(ctx context.Context, approval *domain.Approval) error
	GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Approval, error)
}

type InvestmentRepository interface {
	Create(ctx context.Context, investment *domain.Investment) error
	ListByLoanID(ctx context.Context, loanID uuid.UUID) ([]*domain.Investment, error)
	GetInvestorsByLoanID(ctx context.Context, loanID uuid.UUID) ([]string, error)
}

type DisbursementRepository interface {
	Create(ctx context.Context, disbursement *domain.Disbursement) error
	GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Disbursement, error)
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
