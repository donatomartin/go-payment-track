package router

import (
	"io/fs"
	"log"
	"net/http"

	"pagos-cesar/internal/handler"
	"pagos-cesar/internal/middleware"
	"pagos-cesar/internal/service"
	"pagos-cesar/web"
)

func NewRouter(paymentService service.PaymentService, logger *log.Logger) http.Handler {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(web.WebFS, "static")
	if err != nil {
		log.Fatalf("failed to create static sub file system: %v", err)
	}

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		log.Fatalf("failed to create template sub file system: %v", err)
	}

	// API Handlers
	paymentHandler := handler.NewPaymentHandler(paymentService, logger)
	mux.Handle("/api/v1/payments", paymentHandler)

	// Web Handlers
	dashboardHandler := handler.NewDashboardHandler(logger, templateFS)
	mux.Handle("/", dashboardHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

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
