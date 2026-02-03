CREATE TABLE approvals (
    id UUID PRIMARY KEY,
    loan_id UUID NOT NULL UNIQUE REFERENCES loans(id) ON DELETE CASCADE,
    field_validator_id VARCHAR(255) NOT NULL,
    picture_proof_url TEXT NOT NULL,
    approved_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_approvals_loan_id ON approvals(loan_id);
