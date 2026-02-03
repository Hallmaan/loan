package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/agunghallmanmaliki/amartha/internal/domain"
	"github.com/agunghallmanmaliki/amartha/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LoanRepository struct {
	db *DB
}

func NewLoanRepository(db *DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	conn := r.db.GetConn(ctx)
	query := `
		INSERT INTO loans (id, borrower_id, principal_amount, rate, roi, state, agreement_letter_url, total_invested, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := conn.Exec(ctx, query,
		loan.ID,
		loan.BorrowerID,
		loan.PrincipalAmount,
		loan.Rate,
		loan.ROI,
		loan.State,
		loan.AgreementLetterURL,
		loan.TotalInvested,
		loan.CreatedAt,
		loan.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create loan: %w", err)
	}
	return nil
}

func (r *LoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	conn := r.db.GetConn(ctx)
	query := `
		SELECT id, borrower_id, principal_amount, rate, roi, state, agreement_letter_url, total_invested, created_at, updated_at
		FROM loans
		WHERE id = $1
	`
	return r.scanLoan(conn.QueryRow(ctx, query, id))
}

func (r *LoanRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	conn := r.db.GetConn(ctx)
	query := `
		SELECT id, borrower_id, principal_amount, rate, roi, state, agreement_letter_url, total_invested, created_at, updated_at
		FROM loans
		WHERE id = $1
		FOR UPDATE
	`
	return r.scanLoan(conn.QueryRow(ctx, query, id))
}

func (r *LoanRepository) scanLoan(row pgx.Row) (*domain.Loan, error) {
	var loan domain.Loan
	err := row.Scan(
		&loan.ID,
		&loan.BorrowerID,
		&loan.PrincipalAmount,
		&loan.Rate,
		&loan.ROI,
		&loan.State,
		&loan.AgreementLetterURL,
		&loan.TotalInvested,
		&loan.CreatedAt,
		&loan.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLoanNotFound
		}
		return nil, fmt.Errorf("failed to scan loan: %w", err)
	}
	return &loan, nil
}

func (r *LoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	conn := r.db.GetConn(ctx)
	query := `
		UPDATE loans
		SET borrower_id = $2, principal_amount = $3, rate = $4, roi = $5, state = $6,
		    agreement_letter_url = $7, total_invested = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := conn.Exec(ctx, query,
		loan.ID,
		loan.BorrowerID,
		loan.PrincipalAmount,
		loan.Rate,
		loan.ROI,
		loan.State,
		loan.AgreementLetterURL,
		loan.TotalInvested,
		loan.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}
	return nil
}

func (r *LoanRepository) List(ctx context.Context, filter repository.LoanFilter) ([]*domain.Loan, int64, error) {
	conn := r.db.GetConn(ctx)

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.State != nil {
		conditions = append(conditions, fmt.Sprintf("state = $%d", argIndex))
		args = append(args, *filter.State)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM loans %s", whereClause)
	var total int64
	err := conn.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loans: %w", err)
	}

	// List query
	listQuery := fmt.Sprintf(`
		SELECT id, borrower_id, principal_amount, rate, roi, state, agreement_letter_url, total_invested, created_at, updated_at
		FROM loans
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := conn.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list loans: %w", err)
	}
	defer rows.Close()

	var loans []*domain.Loan
	for rows.Next() {
		var loan domain.Loan
		err := rows.Scan(
			&loan.ID,
			&loan.BorrowerID,
			&loan.PrincipalAmount,
			&loan.Rate,
			&loan.ROI,
			&loan.State,
			&loan.AgreementLetterURL,
			&loan.TotalInvested,
			&loan.CreatedAt,
			&loan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loan: %w", err)
		}
		loans = append(loans, &loan)
	}

	return loans, total, nil
}

// ApprovalRepository

type ApprovalRepository struct {
	db *DB
}

func NewApprovalRepository(db *DB) *ApprovalRepository {
	return &ApprovalRepository{db: db}
}

func (r *ApprovalRepository) Create(ctx context.Context, approval *domain.Approval) error {
	conn := r.db.GetConn(ctx)
	query := `
		INSERT INTO approvals (id, loan_id, field_validator_id, picture_proof_url, approved_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := conn.Exec(ctx, query,
		approval.ID,
		approval.LoanID,
		approval.FieldValidatorID,
		approval.PictureProofURL,
		approval.ApprovedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create approval: %w", err)
	}
	return nil
}

func (r *ApprovalRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Approval, error) {
	conn := r.db.GetConn(ctx)
	query := `
		SELECT id, loan_id, field_validator_id, picture_proof_url, approved_at
		FROM approvals
		WHERE loan_id = $1
	`
	var approval domain.Approval
	err := conn.QueryRow(ctx, query, loanID).Scan(
		&approval.ID,
		&approval.LoanID,
		&approval.FieldValidatorID,
		&approval.PictureProofURL,
		&approval.ApprovedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrApprovalNotFound
		}
		return nil, fmt.Errorf("failed to get approval: %w", err)
	}
	return &approval, nil
}

// InvestmentRepository

type InvestmentRepository struct {
	db *DB
}

func NewInvestmentRepository(db *DB) *InvestmentRepository {
	return &InvestmentRepository{db: db}
}

func (r *InvestmentRepository) Create(ctx context.Context, investment *domain.Investment) error {
	conn := r.db.GetConn(ctx)
	query := `
		INSERT INTO investments (id, loan_id, investor_id, amount, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := conn.Exec(ctx, query,
		investment.ID,
		investment.LoanID,
		investment.InvestorID,
		investment.Amount,
		investment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create investment: %w", err)
	}
	return nil
}

func (r *InvestmentRepository) ListByLoanID(ctx context.Context, loanID uuid.UUID) ([]*domain.Investment, error) {
	conn := r.db.GetConn(ctx)
	query := `
		SELECT id, loan_id, investor_id, amount, created_at
		FROM investments
		WHERE loan_id = $1
		ORDER BY created_at ASC
	`
	rows, err := conn.Query(ctx, query, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to list investments: %w", err)
	}
	defer rows.Close()

	var investments []*domain.Investment
	for rows.Next() {
		var inv domain.Investment
		err := rows.Scan(
			&inv.ID,
			&inv.LoanID,
			&inv.InvestorID,
			&inv.Amount,
			&inv.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan investment: %w", err)
		}
		investments = append(investments, &inv)
	}

	return investments, nil
}

func (r *InvestmentRepository) GetInvestorsByLoanID(ctx context.Context, loanID uuid.UUID) ([]string, error) {
	conn := r.db.GetConn(ctx)
	query := `
		SELECT DISTINCT investor_id
		FROM investments
		WHERE loan_id = $1
	`
	rows, err := conn.Query(ctx, query, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get investors: %w", err)
	}
	defer rows.Close()

	var investors []string
	for rows.Next() {
		var investorID string
		if err := rows.Scan(&investorID); err != nil {
			return nil, fmt.Errorf("failed to scan investor: %w", err)
		}
		investors = append(investors, investorID)
	}

	return investors, nil
}

// DisbursementRepository

type DisbursementRepository struct {
	db *DB
}

func NewDisbursementRepository(db *DB) *DisbursementRepository {
	return &DisbursementRepository{db: db}
}

func (r *DisbursementRepository) Create(ctx context.Context, disbursement *domain.Disbursement) error {
	conn := r.db.GetConn(ctx)
	query := `
		INSERT INTO disbursements (id, loan_id, field_officer_id, signed_agreement_url, disbursed_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := conn.Exec(ctx, query,
		disbursement.ID,
		disbursement.LoanID,
		disbursement.FieldOfficerID,
		disbursement.SignedAgreementURL,
		disbursement.DisbursedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create disbursement: %w", err)
	}
	return nil
}

func (r *DisbursementRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Disbursement, error) {
	conn := r.db.GetConn(ctx)
	query := `
		SELECT id, loan_id, field_officer_id, signed_agreement_url, disbursed_at
		FROM disbursements
		WHERE loan_id = $1
	`
	var disbursement domain.Disbursement
	err := conn.QueryRow(ctx, query, loanID).Scan(
		&disbursement.ID,
		&disbursement.LoanID,
		&disbursement.FieldOfficerID,
		&disbursement.SignedAgreementURL,
		&disbursement.DisbursedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDisbursementNotFound
		}
		return nil, fmt.Errorf("failed to get disbursement: %w", err)
	}
	return &disbursement, nil
}
