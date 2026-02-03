package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/agunghallmanmaliki/amartha/internal/domain"
	"github.com/agunghallmanmaliki/amartha/internal/repository"
	"github.com/google/uuid"
)

type LoanService struct {
	loanRepo         repository.LoanRepository
	approvalRepo     repository.ApprovalRepository
	investmentRepo   repository.InvestmentRepository
	disbursementRepo repository.DisbursementRepository
	txManager        repository.TransactionManager
	emailService     EmailService
	logger           *slog.Logger
}

func NewLoanService(
	loanRepo repository.LoanRepository,
	approvalRepo repository.ApprovalRepository,
	investmentRepo repository.InvestmentRepository,
	disbursementRepo repository.DisbursementRepository,
	txManager repository.TransactionManager,
	emailService EmailService,
	logger *slog.Logger,
) *LoanService {
	return &LoanService{
		loanRepo:         loanRepo,
		approvalRepo:     approvalRepo,
		investmentRepo:   investmentRepo,
		disbursementRepo: disbursementRepo,
		txManager:        txManager,
		emailService:     emailService,
		logger:           logger,
	}
}

func (s *LoanService) CreateLoan(ctx context.Context, borrowerID string, principalAmount int64, rate, roi float64) (*domain.Loan, error) {
	if principalAmount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	loan := domain.NewLoan(borrowerID, principalAmount, rate, roi)

	if err := s.loanRepo.Create(ctx, loan); err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	s.logger.Info("loan created",
		"loan_id", loan.ID,
		"borrower_id", borrowerID,
		"principal", principalAmount,
	)

	return loan, nil
}

func (s *LoanService) GetLoan(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	return s.loanRepo.GetByID(ctx, id)
}

func (s *LoanService) ListLoans(ctx context.Context, filter repository.LoanFilter) ([]*domain.Loan, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.loanRepo.List(ctx, filter)
}

func (s *LoanService) ApproveLoan(ctx context.Context, loanID uuid.UUID, fieldValidatorID, pictureProofURL string) (*domain.Loan, error) {
	var loan *domain.Loan

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		loan, err = s.loanRepo.GetByIDForUpdate(txCtx, loanID)
		if err != nil {
			return err
		}

		if loan.State != domain.LoanStateProposed {
			return domain.ErrLoanAlreadyApproved
		}

		if err := loan.TransitionTo(domain.LoanStateApproved); err != nil {
			return err
		}

		if err := s.loanRepo.Update(txCtx, loan); err != nil {
			return err
		}

		approval := domain.NewApproval(loanID, fieldValidatorID, pictureProofURL)
		if err := s.approvalRepo.Create(txCtx, approval); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.logger.Info("loan approved",
		"loan_id", loanID,
		"field_validator_id", fieldValidatorID,
	)

	return loan, nil
}

func (s *LoanService) AddInvestment(ctx context.Context, loanID uuid.UUID, investorID string, amount int64) (*domain.Loan, *domain.Investment, error) {
	if amount <= 0 {
		return nil, nil, domain.ErrInvalidAmount
	}

	var loan *domain.Loan
	var investment *domain.Investment

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		loan, err = s.loanRepo.GetByIDForUpdate(txCtx, loanID)
		if err != nil {
			return err
		}

		if err := loan.AddInvestment(amount); err != nil {
			return err
		}

		investment = domain.NewInvestment(loanID, investorID, amount)
		if err := s.investmentRepo.Create(txCtx, investment); err != nil {
			return err
		}

		// Auto-transition to invested if fully funded
		if loan.IsFullyInvested() {
			if err := loan.TransitionTo(domain.LoanStateInvested); err != nil {
				return err
			}
		}

		if err := s.loanRepo.Update(txCtx, loan); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	s.logger.Info("investment added",
		"loan_id", loanID,
		"investor_id", investorID,
		"amount", amount,
		"total_invested", loan.TotalInvested,
	)

	// Send emails to all investors if fully invested (async)
	if loan.IsFullyInvested() && loan.AgreementLetterURL != nil {
		go s.notifyInvestors(context.Background(), loanID, *loan.AgreementLetterURL)
	}

	return loan, investment, nil
}

func (s *LoanService) notifyInvestors(ctx context.Context, loanID uuid.UUID, agreementURL string) {
	investors, err := s.investmentRepo.GetInvestorsByLoanID(ctx, loanID)
	if err != nil {
		s.logger.Error("failed to get investors for notification",
			"loan_id", loanID,
			"error", err,
		)
		return
	}

	for _, investorID := range investors {
		if err := s.emailService.SendAgreementEmail(ctx, investorID, loanID.String(), agreementURL); err != nil {
			s.logger.Error("failed to send email to investor",
				"investor_id", investorID,
				"loan_id", loanID,
				"error", err,
			)
		}
	}
}

func (s *LoanService) ListInvestments(ctx context.Context, loanID uuid.UUID) ([]*domain.Investment, error) {
	// Verify loan exists
	_, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	return s.investmentRepo.ListByLoanID(ctx, loanID)
}

func (s *LoanService) DisburseLoan(ctx context.Context, loanID uuid.UUID, fieldOfficerID, signedAgreementURL string) (*domain.Loan, error) {
	var loan *domain.Loan

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		loan, err = s.loanRepo.GetByIDForUpdate(txCtx, loanID)
		if err != nil {
			return err
		}

		if loan.State != domain.LoanStateInvested {
			if loan.State == domain.LoanStateDisbursed {
				return domain.ErrLoanAlreadyDisbursed
			}
			return domain.ErrLoanNotInvested
		}

		if err := loan.TransitionTo(domain.LoanStateDisbursed); err != nil {
			return err
		}

		loan.AgreementLetterURL = &signedAgreementURL

		if err := s.loanRepo.Update(txCtx, loan); err != nil {
			return err
		}

		disbursement := domain.NewDisbursement(loanID, fieldOfficerID, signedAgreementURL)
		if err := s.disbursementRepo.Create(txCtx, disbursement); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.logger.Info("loan disbursed",
		"loan_id", loanID,
		"field_officer_id", fieldOfficerID,
	)

	// Notify all investors about disbursement with the signed agreement
	go s.notifyInvestors(context.Background(), loanID, signedAgreementURL)

	return loan, nil
}

func (s *LoanService) GetApproval(ctx context.Context, loanID uuid.UUID) (*domain.Approval, error) {
	return s.approvalRepo.GetByLoanID(ctx, loanID)
}

func (s *LoanService) GetDisbursement(ctx context.Context, loanID uuid.UUID) (*domain.Disbursement, error) {
	return s.disbursementRepo.GetByLoanID(ctx, loanID)
}
