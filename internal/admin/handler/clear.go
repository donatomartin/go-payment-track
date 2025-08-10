package handler

import (
	"net/http"
)

func (h *ApiAdminHandler) clearDatabase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.paymentRepo.Clear(r.Context()); err != nil {
		h.logger.Printf("error clearing payments: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.invoiceRepo.Clear(r.Context()); err != nil {
		h.logger.Printf("error clearing invoices: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("database cleared"))
}
