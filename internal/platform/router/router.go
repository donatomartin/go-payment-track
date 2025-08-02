package router

import (
	"log"
	"net/http"

	"pagos-cesar/internal/payment"
	"pagos-cesar/internal/platform/middleware"
	"pagos-cesar/internal/web/dashboard"
	"pagos-cesar/internal/web/static"
)

func NewRouter(paymentService payment.PaymentService, logger *log.Logger) http.Handler {

	mux := http.NewServeMux()

	// API Handlers
	payment.NewApiPaymentHandler(paymentService, logger).RegisterRoutes(mux)

	// Web Handlers
	dashboard.NewDashboardHandler(paymentService, logger).RegisterRoutes(mux)
	payment.NewWebPaymentHandler(paymentService, logger).RegisterRoutes(mux)

	// Static Handlers
	static.NewStaticHandler(logger).RegisterRoutes(mux)

	// Wrap with fallback 404 handler
	return middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		handler, pattern := mux.Handler(r)
		if pattern == "" {
			logger.Printf("404 Not Found: %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	}), logger)

}
