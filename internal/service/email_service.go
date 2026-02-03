package service

import (
	"context"
	"log/slog"
)

type EmailService interface {
	SendAgreementEmail(ctx context.Context, investorID string, loanID string, agreementURL string) error
}

type MockEmailService struct {
	logger *slog.Logger
}

func NewMockEmailService(logger *slog.Logger) *MockEmailService {
	return &MockEmailService{logger: logger}
}

func (s *MockEmailService) SendAgreementEmail(ctx context.Context, investorID string, loanID string, agreementURL string) error {
	s.logger.Info("sending agreement email",
		"investor_id", investorID,
		"loan_id", loanID,
		"agreement_url", agreementURL,
	)
	return nil
}
