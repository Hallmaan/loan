# Loan Service API - Technical Specification

## Overview

A RESTful Loan Service API in Golang implementing a forward-only state machine for loan lifecycle management.

## Loan State Machine

```
proposed → approved → invested → disbursed
```

- **proposed**: Initial state when loan is created
- **approved**: Approved by staff (requires: picture proof, employee ID, date)
- **invested**: Fully funded by investors (auto-transitions when total = principal)
- **disbursed**: Loan given to borrower (requires: signed agreement, employee ID, date)

## REST API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/loans` | Create loan (proposed state) |
| GET | `/api/v1/loans` | List loans with pagination/filters |
| GET | `/api/v1/loans/{id}` | Get loan details |
| POST | `/api/v1/loans/{id}/approve` | Approve loan (multipart: picture proof) |
| POST | `/api/v1/loans/{id}/investments` | Add investment |
| GET | `/api/v1/loans/{id}/investments` | List investments |
| POST | `/api/v1/loans/{id}/disburse` | Disburse loan (multipart: signed agreement) |

## API Request/Response Examples

### Create Loan

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": "borrower-123",
    "principal_amount": 1000000,
    "rate": 0.15,
    "roi": 0.12
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "borrower_id": "borrower-123",
    "principal_amount": 1000000,
    "rate": 0.15,
    "roi": 0.12,
    "state": "proposed",
    "total_invested": 0,
    "remaining_amount": 1000000,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Approve Loan

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/loans/{id}/approve \
  -F "field_validator_id=validator-456" \
  -F "picture_proof=@proof.jpg"
```

### Add Investment

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/loans/{id}/investments \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": "investor-789",
    "amount": 500000
  }'
```

### Disburse Loan

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/loans/{id}/disburse \
  -F "field_officer_id=officer-123" \
  -F "signed_agreement=@agreement.pdf"
```

## Database Schema

### loans
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| borrower_id | VARCHAR(255) | Borrower identifier |
| principal_amount | BIGINT | Loan amount in cents |
| rate | DECIMAL(10,4) | Interest rate |
| roi | DECIMAL(10,4) | Return on investment |
| state | ENUM | proposed, approved, invested, disbursed |
| agreement_letter_url | TEXT | URL to signed agreement |
| total_invested | BIGINT | Total invested amount |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |

### approvals
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| loan_id | UUID | Foreign key to loans |
| field_validator_id | VARCHAR(255) | Validator employee ID |
| picture_proof_url | TEXT | URL to proof picture |
| approved_at | TIMESTAMP | Approval timestamp |

### investments
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| loan_id | UUID | Foreign key to loans |
| investor_id | VARCHAR(255) | Investor identifier |
| amount | BIGINT | Investment amount |
| created_at | TIMESTAMP | Investment timestamp |

### disbursements
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| loan_id | UUID | Foreign key to loans |
| field_officer_id | VARCHAR(255) | Officer employee ID |
| signed_agreement_url | TEXT | URL to signed agreement |
| disbursed_at | TIMESTAMP | Disbursement timestamp |

## Error Responses

| Status | Code | Description |
|--------|------|-------------|
| 400 | BAD_REQUEST | Invalid request format |
| 400 | VALIDATION_ERROR | Validation failed |
| 404 | NOT_FOUND | Resource not found |
| 422 | INVALID_STATE_TRANSITION | Invalid state transition |
| 422 | LOAN_NOT_APPROVED | Loan must be approved for investments |
| 422 | LOAN_NOT_INVESTED | Loan must be invested for disbursement |
| 422 | INVESTMENT_EXCEEDS_LIMIT | Investment exceeds remaining principal |
| 500 | INTERNAL_ERROR | Internal server error |

## Technology Stack

- **Language**: Go 1.21+
- **HTTP**: Standard library (net/http)
- **Database**: PostgreSQL 16
- **Database Driver**: pgx/v5
- **Migrations**: golang-migrate
- **Validation**: go-playground/validator
- **UUID**: google/uuid
- **Logging**: log/slog
- **File Storage**: Local filesystem

## Getting Started

1. Start PostgreSQL:
   ```bash
   docker-compose up -d
   ```

2. Install migrate tool:
   ```bash
   make install-migrate
   ```

3. Run migrations:
   ```bash
   make migrate
   ```

4. Start the server:
   ```bash
   make run
   ```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| SERVER_PORT | 8080 | Server port |
| SERVER_HOST | http://localhost:8080 | Server host URL |
| DATABASE_URL | postgres://postgres:postgres@localhost:5432/amartha?sslmode=disable | PostgreSQL connection string |
| STORAGE_PATH | ./uploads | File storage directory |
| MAX_FILE_SIZE | 10485760 | Max upload size (10MB) |
