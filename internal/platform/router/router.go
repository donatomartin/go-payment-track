package router

import (
	"log"
	"net/http"

	invoiceRepo "app/internal/invoice/repository"
	paymentHandler "app/internal/payment/handler"
	paymentRepo "app/internal/payment/repository"
	"app/internal/platform/middleware"
	"app/internal/web/dashboard"
	"app/internal/web/static"
)

func NewRouter(paymentRepo paymentRepo.PaymentRepository, invoiceRepo invoiceRepo.InvoiceRepository, logger *log.Logger) http.Handler {

	mux := http.NewServeMux()

	// API Handlers
	paymentHandler.NewApiPaymentHandler(paymentRepo, logger, mux).RegisterRoutes()

	// Web Handlers
	dashboard.NewDashboardHandler(paymentRepo, invoiceRepo, logger).RegisterRoutes(mux)

	// Static Handlers
	static.NewStaticHandler(logger).RegisterRoutes(mux)

	// Wrap with fallback 404 handler
	return middleware.Logging(mux, logger)

}
