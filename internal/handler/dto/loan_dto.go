package dto

import (
	"time"

	"github.com/agunghallmanmaliki/amartha/internal/domain"
)

// Request DTOs

type CreateLoanRequest struct {
	BorrowerID      string  `json:"borrower_id" validate:"required"`
	PrincipalAmount int64   `json:"principal_amount" validate:"required,gt=0"`
	Rate            float64 `json:"rate" validate:"required,gte=0"`
	ROI             float64 `json:"roi" validate:"required,gte=0"`
}

type AddInvestmentRequest struct {
	InvestorID string `json:"investor_id" validate:"required"`
	Amount     int64  `json:"amount" validate:"required,gt=0"`
}

// Response DTOs

type LoanResponse struct {
	ID                 string    `json:"id"`
	BorrowerID         string    `json:"borrower_id"`
	PrincipalAmount    int64     `json:"principal_amount"`
	Rate               float64   `json:"rate"`
	ROI                float64   `json:"roi"`
	State              string    `json:"state"`
	AgreementLetterURL *string   `json:"agreement_letter_url,omitempty"`
	TotalInvested      int64     `json:"total_invested"`
	RemainingAmount    int64     `json:"remaining_amount"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func ToLoanResponse(loan *domain.Loan) *LoanResponse {
	return &LoanResponse{
		ID:                 loan.ID.String(),
		BorrowerID:         loan.BorrowerID,
		PrincipalAmount:    loan.PrincipalAmount,
		Rate:               loan.Rate,
		ROI:                loan.ROI,
		State:              string(loan.State),
		AgreementLetterURL: loan.AgreementLetterURL,
		TotalInvested:      loan.TotalInvested,
		RemainingAmount:    loan.RemainingAmount(),
		CreatedAt:          loan.CreatedAt,
		UpdatedAt:          loan.UpdatedAt,
	}
}

func ToLoanResponses(loans []*domain.Loan) []*LoanResponse {
	responses := make([]*LoanResponse, len(loans))
	for i, loan := range loans {
		responses[i] = ToLoanResponse(loan)
	}
	return responses
}

type InvestmentResponse struct {
	ID         string    `json:"id"`
	LoanID     string    `json:"loan_id"`
	InvestorID string    `json:"investor_id"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToInvestmentResponse(inv *domain.Investment) *InvestmentResponse {
	return &InvestmentResponse{
		ID:         inv.ID.String(),
		LoanID:     inv.LoanID.String(),
		InvestorID: inv.InvestorID,
		Amount:     inv.Amount,
		CreatedAt:  inv.CreatedAt,
	}
}

func ToInvestmentResponses(investments []*domain.Investment) []*InvestmentResponse {
	responses := make([]*InvestmentResponse, len(investments))
	for i, inv := range investments {
		responses[i] = ToInvestmentResponse(inv)
	}
	return responses
}

type ApprovalResponse struct {
	ID               string    `json:"id"`
	LoanID           string    `json:"loan_id"`
	FieldValidatorID string    `json:"field_validator_id"`
	PictureProofURL  string    `json:"picture_proof_url"`
	ApprovedAt       time.Time `json:"approved_at"`
}

func ToApprovalResponse(approval *domain.Approval) *ApprovalResponse {
	return &ApprovalResponse{
		ID:               approval.ID.String(),
		LoanID:           approval.LoanID.String(),
		FieldValidatorID: approval.FieldValidatorID,
		PictureProofURL:  approval.PictureProofURL,
		ApprovedAt:       approval.ApprovedAt,
	}
}

type DisbursementResponse struct {
	ID                 string    `json:"id"`
	LoanID             string    `json:"loan_id"`
	FieldOfficerID     string    `json:"field_officer_id"`
	SignedAgreementURL string    `json:"signed_agreement_url"`
	DisbursedAt        time.Time `json:"disbursed_at"`
}

func ToDisburseResponse(disbursement *domain.Disbursement) *DisbursementResponse {
	return &DisbursementResponse{
		ID:                 disbursement.ID.String(),
		LoanID:             disbursement.LoanID.String(),
		FieldOfficerID:     disbursement.FieldOfficerID,
		SignedAgreementURL: disbursement.SignedAgreementURL,
		DisbursedAt:        disbursement.DisbursedAt,
	}
}
