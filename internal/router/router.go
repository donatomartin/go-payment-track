package router

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"pagos-cesar/internal/handler"
	"pagos-cesar/internal/middleware"
	"pagos-cesar/internal/service"
)

var staticFiles embed.FS

func NewRouter(paymentService service.PaymentService, logger *log.Logger) http.Handler {
	mux := http.NewServeMux()

	fs, err := fs.Sub(staticFiles, "web/static")
	if err != nil {
		log.Fatalf("failed to create sub file system: %v", err)
	}

	// API Handlers
	apiRouter := http.NewServeMux()
	paymentHandler := handler.NewPaymentHandler(paymentService, logger)
	apiRouter.Handle("/api/v1/payments", paymentHandler)

	// Web Handlers
	webRouter := http.NewServeMux()
	dashboardHandler := handler.NewDashboardHandler(logger)
	webRouter.Handle("/", dashboardHandler)
	webRouter.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(fs))))

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
