package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agunghallmanmaliki/amartha/internal/config"
	"github.com/agunghallmanmaliki/amartha/internal/handler"
	"github.com/agunghallmanmaliki/amartha/internal/repository/postgres"
	"github.com/agunghallmanmaliki/amartha/internal/service"
	"github.com/agunghallmanmaliki/amartha/internal/storage/local"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load configuration
	cfg := config.Load()

	// Initialize database
	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize storage
	storage, err := local.NewLocalStorage(cfg.StoragePath, cfg.ServerHost)
	if err != nil {
		logger.Error("failed to initialize storage", "error", err)
		os.Exit(1)
	}

	// Initialize repositories
	loanRepo := postgres.NewLoanRepository(db)
	approvalRepo := postgres.NewApprovalRepository(db)
	investmentRepo := postgres.NewInvestmentRepository(db)
	disbursementRepo := postgres.NewDisbursementRepository(db)

	// Initialize services
	emailService := service.NewMockEmailService(logger)
	loanService := service.NewLoanService(
		loanRepo,
		approvalRepo,
		investmentRepo,
		disbursementRepo,
		db,
		emailService,
		logger,
	)

	// Initialize handlers
	loanHandler := handler.NewLoanHandler(loanService, storage, cfg.MaxFileSize)

	// Setup router
	router := handler.NewRouter(loanHandler, logger)
	httpHandler := router.Setup()

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("starting server", "port", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server stopped")
}
