package dashboard

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"pagos-cesar/internal/payment"
	"pagos-cesar/web"
)

type DashboardHandler struct {
	logger    *log.Logger
	templates *template.Template
	service   payment.PaymentService
}

func NewDashboardHandler(service payment.PaymentService, logger *log.Logger) *DashboardHandler {

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		log.Fatalf("failed to create template sub file system: %v", err)
	}
	t := template.Must(template.ParseFS(templateFS, "*.html"))

	return &DashboardHandler{
		logger:    logger,
		templates: t,
		service:   service,
	}
}

func (h *DashboardHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/dashboard", h.getDashboard)
	mux.HandleFunc("/db", h.getDashboard) // Redirect root to dashboard
}

func (h *DashboardHandler) getDashboard(w http.ResponseWriter, r *http.Request) {

	payments, err := h.service.GetAllPayments(r.Context())
	if err != nil {
		http.Error(w, "Failed to get payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		Payments []payment.Payment
	}{
		Title:    "Dashboard",
		Payments: payments,
	}

	if err := h.templates.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		http.Error(w, "Failed to render dashboard: ", http.StatusInternalServerError)
		h.logger.Fatal("Template error:", err)
	}

}
