package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/agunghallmanmaliki/amartha/internal/handler/middleware"
)

type Router struct {
	mux     *http.ServeMux
	handler *LoanHandler
	logger  *slog.Logger
}

func NewRouter(handler *LoanHandler, logger *slog.Logger) *Router {
	return &Router{
		mux:     http.NewServeMux(),
		handler: handler,
		logger:  logger,
	}
}

func (r *Router) Setup() http.Handler {
	// Register routes
	r.mux.HandleFunc("/api/v1/loans", r.loansHandler)
	r.mux.HandleFunc("/api/v1/loans/", r.loanDetailHandler)

	// Health check
	r.mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Serve static files for uploads
	r.mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Apply middleware
	var handler http.Handler = r.mux
	handler = middleware.Logger(r.logger)(handler)
	handler = middleware.Recovery(r.logger)(handler)

	return handler
}

func (r *Router) loansHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		r.handler.CreateLoan(w, req)
	case http.MethodGet:
		r.handler.ListLoans(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) loanDetailHandler(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/loans/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// /api/v1/loans/{id}
	if len(parts) == 1 {
		switch req.Method {
		case http.MethodGet:
			r.handler.GetLoan(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// /api/v1/loans/{id}/{action}
	if len(parts) == 2 {
		action := parts[1]
		switch action {
		case "approve":
			if req.Method == http.MethodPost {
				r.handler.ApproveLoan(w, req)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case "investments":
			switch req.Method {
			case http.MethodPost:
				r.handler.AddInvestment(w, req)
			case http.MethodGet:
				r.handler.ListInvestments(w, req)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case "disburse":
			if req.Method == http.MethodPost {
				r.handler.DisburseLoan(w, req)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}
