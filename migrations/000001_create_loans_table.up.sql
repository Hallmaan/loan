CREATE TYPE loan_state AS ENUM ('proposed', 'approved', 'invested', 'disbursed');

CREATE TABLE loans (
    id UUID PRIMARY KEY,
    borrower_id VARCHAR(255) NOT NULL,
    principal_amount BIGINT NOT NULL,
    rate DECIMAL(10, 4) NOT NULL,
    roi DECIMAL(10, 4) NOT NULL,
    state loan_state NOT NULL DEFAULT 'proposed',
    agreement_letter_url TEXT,
    total_invested BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_loans_state ON loans(state);
CREATE INDEX idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX idx_loans_created_at ON loans(created_at DESC);
