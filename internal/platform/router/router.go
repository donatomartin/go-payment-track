package router

import (
	"log"
	"net/http"

	"app/internal/invoice"
	"app/internal/payment"
	"app/internal/platform/middleware"
	"app/internal/web/dashboard"
	"app/internal/web/static"
)

func NewRouter(paymentService payment.PaymentService, invoiceService invoice.InvoiceService, logger *log.Logger) http.Handler {

	mux := http.NewServeMux()

	// API Handlers
	payment.NewApiPaymentHandler(paymentService, logger).RegisterRoutes(mux)

	// Web Handlers
	dashboard.NewDashboardHandler(paymentService, invoiceService, logger).RegisterRoutes(mux)

	// Static Handlers
	static.NewStaticHandler(logger).RegisterRoutes(mux)

	// Wrap with fallback 404 handler
	return middleware.Logging(mux, logger)

}
