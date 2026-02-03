CREATE TABLE disbursements (
    id UUID PRIMARY KEY,
    loan_id UUID NOT NULL UNIQUE REFERENCES loans(id) ON DELETE CASCADE,
    field_officer_id VARCHAR(255) NOT NULL,
    signed_agreement_url TEXT NOT NULL,
    disbursed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_disbursements_loan_id ON disbursements(loan_id);
