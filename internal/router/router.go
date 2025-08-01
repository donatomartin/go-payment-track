package router

import (
	"log"
	"net/http"

	"pagos-cesar/internal/handler"
	"pagos-cesar/internal/middleware"
	"pagos-cesar/internal/service"
)

func NewRouter(paymentService service.PaymentService, logger *log.Logger) http.Handler {
	mux := http.NewServeMux()

	// API Handlers
	paymentHandler := handler.NewPaymentHandler(paymentService, logger)
	mux.Handle("/api/v1/payments", paymentHandler)

	// Dashboard Html
	dashboardHandler := handler.NewDashboardHandler(logger)
	mux.Handle("/", dashboardHandler)

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

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
