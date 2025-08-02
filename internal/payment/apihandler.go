package payment

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiPaymentHandler struct {
	service PaymentService
	logger  *log.Logger
}

func NewApiPaymentHandler(service PaymentService, logger *log.Logger) *ApiPaymentHandler {
	return &ApiPaymentHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ApiPaymentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/payments", h.getPayments)
}

func (h *ApiPaymentHandler) getPayments(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.Println(r.URL.Path)

	payments, err := h.service.GetAllPayments(r.Context())
	if err != nil {
		http.Error(w, "Failed to get payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)

}
