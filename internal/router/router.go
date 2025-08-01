package router

import (
	"io/fs"
	"log"
	"net/http"
	"os"

	"pagos-cesar/internal/handler"
	"pagos-cesar/internal/middleware"
	"pagos-cesar/internal/service"
	"pagos-cesar/web"
)

// noListFileSystem is a custom file system that prevents directory listing.
type noListFileSystem struct {
	fs http.FileSystem
}

// Open opens the file with the given name, but returns os.ErrNotExist for directories.
func (nlfs noListFileSystem) Open(name string) (http.File, error) {
	f, err := nlfs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

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

	fileServer := http.FileServer(noListFileSystem{http.FS(staticFS)})
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving static file: %s", r.URL.Path)
		fileServer.ServeHTTP(w, r)
	})))

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
