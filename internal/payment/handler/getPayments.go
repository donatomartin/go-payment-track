package handler

import (
	"encoding/json"
	"net/http"
)

func (h *ApiPaymentHandler) getPayments(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payments, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)

}
