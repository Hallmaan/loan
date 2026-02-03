package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/agunghallmanmaliki/amartha/internal/domain"
	"github.com/agunghallmanmaliki/amartha/internal/handler/dto"
	"github.com/agunghallmanmaliki/amartha/internal/repository"
	"github.com/agunghallmanmaliki/amartha/internal/service"
	"github.com/agunghallmanmaliki/amartha/internal/storage"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LoanHandler struct {
	loanService *service.LoanService
	storage     storage.Storage
	validator   *validator.Validate
	maxFileSize int64
}

func NewLoanHandler(loanService *service.LoanService, storage storage.Storage, maxFileSize int64) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
		storage:     storage,
		validator:   validator.New(),
		maxFileSize: maxFileSize,
	}
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", formatValidationError(err))
		return
	}

	loan, err := h.loanService.CreateLoan(r.Context(), req.BorrowerID, req.PrincipalAmount, req.Rate, req.ROI)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusCreated, dto.ToLoanResponse(loan))
}

func (h *LoanHandler) GetLoan(w http.ResponseWriter, r *http.Request) {
	loanID, err := h.extractLoanID(r)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid loan ID format")
		return
	}

	loan, err := h.loanService.GetLoan(r.Context(), loanID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, dto.ToLoanResponse(loan))
}

func (h *LoanHandler) ListLoans(w http.ResponseWriter, r *http.Request) {
	filter := repository.LoanFilter{
		Limit:  10,
		Offset: 0,
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if stateStr := r.URL.Query().Get("state"); stateStr != "" {
		state := domain.LoanState(stateStr)
		filter.State = &state
	}

	loans, total, err := h.loanService.ListLoans(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	dto.WriteJSONPaginated(w, http.StatusOK, dto.ToLoanResponses(loans), total, filter.Limit, filter.Offset)
}

func (h *LoanHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	loanID, err := h.extractLoanID(r)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid loan ID format")
		return
	}

	if err := r.ParseMultipartForm(h.maxFileSize); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "INVALID_FORM", "Failed to parse multipart form")
		return
	}

	fieldValidatorID := r.FormValue("field_validator_id")
	if fieldValidatorID == "" {
		dto.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "field_validator_id is required")
		return
	}

	file, header, err := r.FormFile("picture_proof")
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "picture_proof file is required")
		return
	}
	defer file.Close()

	filename, err := h.storage.Save(r.Context(), header.Filename, file)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "STORAGE_ERROR", "Failed to save file")
		return
	}

	pictureProofURL := h.storage.GetURL(filename)

	loan, err := h.loanService.ApproveLoan(r.Context(), loanID, fieldValidatorID, pictureProofURL)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, dto.ToLoanResponse(loan))
}

func (h *LoanHandler) AddInvestment(w http.ResponseWriter, r *http.Request) {
	loanID, err := h.extractLoanID(r)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid loan ID format")
		return
	}

	var req dto.AddInvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", formatValidationError(err))
		return
	}

	loan, investment, err := h.loanService.AddInvestment(r.Context(), loanID, req.InvestorID, req.Amount)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := struct {
		Loan       *dto.LoanResponse       `json:"loan"`
		Investment *dto.InvestmentResponse `json:"investment"`
	}{
		Loan:       dto.ToLoanResponse(loan),
		Investment: dto.ToInvestmentResponse(investment),
	}

	dto.WriteJSON(w, http.StatusCreated, response)
}

func (h *LoanHandler) ListInvestments(w http.ResponseWriter, r *http.Request) {
	loanID, err := h.extractLoanID(r)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid loan ID format")
		return
	}

	investments, err := h.loanService.ListInvestments(r.Context(), loanID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, dto.ToInvestmentResponses(investments))
}

func (h *LoanHandler) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	loanID, err := h.extractLoanID(r)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid loan ID format")
		return
	}

	if err := r.ParseMultipartForm(h.maxFileSize); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "INVALID_FORM", "Failed to parse multipart form")
		return
	}

	fieldOfficerID := r.FormValue("field_officer_id")
	if fieldOfficerID == "" {
		dto.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "field_officer_id is required")
		return
	}

	file, header, err := r.FormFile("signed_agreement")
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "signed_agreement file is required")
		return
	}
	defer file.Close()

	filename, err := h.storage.Save(r.Context(), header.Filename, file)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "STORAGE_ERROR", "Failed to save file")
		return
	}

	signedAgreementURL := h.storage.GetURL(filename)

	loan, err := h.loanService.DisburseLoan(r.Context(), loanID, fieldOfficerID, signedAgreementURL)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, dto.ToLoanResponse(loan))
}

func (h *LoanHandler) extractLoanID(r *http.Request) (uuid.UUID, error) {
	// Extract from path: /api/v1/loans/{id}/...
	path := r.URL.Path
	parts := strings.Split(path, "/")

	// Find the ID after "loans"
	for i, part := range parts {
		if part == "loans" && i+1 < len(parts) {
			return uuid.Parse(parts[i+1])
		}
	}

	return uuid.Nil, errors.New("loan ID not found in path")
}

func (h *LoanHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrLoanNotFound):
		dto.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Loan not found")
	case errors.Is(err, domain.ErrInvalidStateTransition):
		dto.WriteError(w, http.StatusUnprocessableEntity, "INVALID_STATE_TRANSITION", "Invalid state transition")
	case errors.Is(err, domain.ErrInvestmentExceedsLimit):
		dto.WriteError(w, http.StatusUnprocessableEntity, "INVESTMENT_EXCEEDS_LIMIT", "Investment amount exceeds remaining principal")
	case errors.Is(err, domain.ErrLoanNotApproved):
		dto.WriteError(w, http.StatusUnprocessableEntity, "LOAN_NOT_APPROVED", "Loan must be in approved state to accept investments")
	case errors.Is(err, domain.ErrLoanNotInvested):
		dto.WriteError(w, http.StatusUnprocessableEntity, "LOAN_NOT_INVESTED", "Loan must be in invested state to disburse")
	case errors.Is(err, domain.ErrLoanAlreadyApproved):
		dto.WriteError(w, http.StatusUnprocessableEntity, "LOAN_ALREADY_APPROVED", "Loan is already approved")
	case errors.Is(err, domain.ErrLoanAlreadyDisbursed):
		dto.WriteError(w, http.StatusUnprocessableEntity, "LOAN_ALREADY_DISBURSED", "Loan is already disbursed")
	case errors.Is(err, domain.ErrInvalidAmount):
		dto.WriteError(w, http.StatusBadRequest, "INVALID_AMOUNT", "Amount must be greater than zero")
	default:
		dto.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
	}
}

func formatValidationError(err error) string {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var messages []string
		for _, e := range validationErrs {
			messages = append(messages, e.Field()+" is invalid")
		}
		return strings.Join(messages, ", ")
	}
	return "Validation failed"
}
