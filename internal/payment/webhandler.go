package payment

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"pagos-cesar/web"
)

type WebPaymentHandler struct {
	service   PaymentService
	logger    *log.Logger
	templates *template.Template
}

func NewWebPaymentHandler(service PaymentService, logger *log.Logger) *WebPaymentHandler {

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		log.Fatalf("failed to create template sub file system: %v", err)
	}
	t := template.Must(template.ParseFS(templateFS, "*.html"))

	return &WebPaymentHandler{
		service:   service,
		logger:    logger,
		templates: t,
	}
}

func (h *WebPaymentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/payments", h.getPayments)
}

func (h *WebPaymentHandler) getPayments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payments, err := h.service.GetAllPayments(r.Context())
	if err != nil {
		http.Error(w, "Failed to get payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = h.templates.ExecuteTemplate(w, "payments.html", payments)
	if err != nil {
		http.Error(w, "Failed to render payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
