package domain

import (
	"time"

	"github.com/google/uuid"
)

type Approval struct {
	ID               uuid.UUID
	LoanID           uuid.UUID
	FieldValidatorID string
	PictureProofURL  string
	ApprovedAt       time.Time
}

func NewApproval(loanID uuid.UUID, fieldValidatorID, pictureProofURL string) *Approval {
	return &Approval{
		ID:               uuid.New(),
		LoanID:           loanID,
		FieldValidatorID: fieldValidatorID,
		PictureProofURL:  pictureProofURL,
		ApprovedAt:       time.Now(),
	}
}
