package domain

import (
	"time"

	"github.com/google/uuid"
)

type Disbursement struct {
	ID                 uuid.UUID
	LoanID             uuid.UUID
	FieldOfficerID     string
	SignedAgreementURL string
	DisbursedAt        time.Time
}

func NewDisbursement(loanID uuid.UUID, fieldOfficerID, signedAgreementURL string) *Disbursement {
	return &Disbursement{
		ID:                 uuid.New(),
		LoanID:             loanID,
		FieldOfficerID:     fieldOfficerID,
		SignedAgreementURL: signedAgreementURL,
		DisbursedAt:        time.Now(),
	}
}
