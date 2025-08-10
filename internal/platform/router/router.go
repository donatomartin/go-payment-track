package router

import (
	"log"
	"net/http"

	invoiceHandler "app/internal/invoice/handler"
	invoiceRepo "app/internal/invoice/repository"
	paymentHandler "app/internal/payment/handler"
	paymentRepo "app/internal/payment/repository"
	"app/internal/platform/middleware"
	"app/internal/web/static"
	"app/internal/web/templates"
)

func NewRouter(paymentRepo paymentRepo.PaymentRepository, invoiceRepo invoiceRepo.InvoiceRepository, logger *log.Logger) http.Handler {

	mux := http.NewServeMux()

	// API Handlers
	paymentHandler.NewApiPaymentHandler(paymentRepo, logger, mux).RegisterRoutes()
	invoiceHandler.NewApiInvoiceHandler(invoiceRepo, logger, mux).RegisterRoutes()

	// Template Handlers
	templates.NewDashboardHandler(paymentRepo, invoiceRepo, logger).RegisterRoutes(mux)

	// Static Handlers
	static.NewStaticHandler(logger).RegisterRoutes(mux)

	// Wrap with fallback 404 handler
	return middleware.Logging(mux, logger)

}
