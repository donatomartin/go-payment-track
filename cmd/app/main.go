package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/invoice"
	"app/internal/payment"
	"app/internal/platform/config"
	"app/internal/platform/database"
	"app/internal/platform/router"
)

func main() {

	// Load cfg
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Open or create the log file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Setup logger
	logger := log.New(io.MultiWriter(os.Stdout, logFile), "app", log.LstdFlags)

	// Setup DB
	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Do Development tasks
	if cfg.Env == "dev" {

		// Apply schema if in dev
		if err := database.ApplySchema(db); err != nil {
			logger.Fatalf("Schema setup failed: %v", err)

		} else {
			logger.Println("Schema applied successfully")
		}

	}

	// Initialize repositories and services
	paymentRepo := payment.NewPaymentRepository(db)
	paymentService := payment.NewPaymentService(paymentRepo)

	invoiceRepo := invoice.NewInvoiceRepository(db)
	invoiceService := invoice.NewInvoiceService(*invoiceRepo)

	// Setup HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router.NewRouter(paymentService, *invoiceService, logger),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Println("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown error: %v", err)
	}
	logger.Println("server stopped")

}
