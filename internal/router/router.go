package router

import (
	"log"
	"net/http"

	v1apihandler "pagos-cesar/internal/handler/api/v1"
	statichandler "pagos-cesar/internal/handler/static"
	webhandler "pagos-cesar/internal/handler/web"
	"pagos-cesar/internal/middleware"
	"pagos-cesar/internal/service"
)

func NewRouter(paymentService service.PaymentService, logger *log.Logger) http.Handler {

	mux := http.NewServeMux()

	// API Handlers
	v1apihandler.NewPaymentHandler(paymentService, logger).RegisterRoutes(mux)

	// Web Handlers
	webhandler.NewPaymentHandler(paymentService, logger).RegisterRoutes(mux)
	webhandler.NewDashboardHandler(paymentService, logger).RegisterRoutes(mux)

	// Static Handlers
	statichandler.NewStaticHandler(logger).RegisterRoutes(mux)

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
