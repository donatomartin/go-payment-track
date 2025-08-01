package webhandler

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"pagos-cesar/internal/service"
	"pagos-cesar/web"
)

type PaymentHandler struct {
	service   service.PaymentService
	logger    *log.Logger
	templates *template.Template
}

func NewPaymentHandler(service service.PaymentService, logger *log.Logger) *PaymentHandler {

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		log.Fatalf("failed to create template sub file system: %v", err)
	}
	t := template.Must(template.ParseFS(templateFS, "*.html"))

	return &PaymentHandler{
		service:   service,
		logger:    logger,
		templates: t,
	}
}

func (h *PaymentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/payments", h.getPayments)
}

func (h *PaymentHandler) getPayments(w http.ResponseWriter, r *http.Request) {
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
