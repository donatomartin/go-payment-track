package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"pagos-cesar/internal/service"
)

type PaymentHandler struct {
	service service.PaymentService
	logger  *log.Logger
}

func NewPaymentHandler(service service.PaymentService, logger *log.Logger) *PaymentHandler {
	return &PaymentHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
