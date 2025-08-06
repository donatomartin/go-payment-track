package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *ApiPaymentHandler) getPagedPayments(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	limit, _ := strconv.Atoi(pageSize)
	// Calculate offset based on page and pageSize
	offset := 0
	if page != "" && pageSize != "" {
		// Assuming page is 1-indexed
		offset, _ = strconv.Atoi(page)
		offset = (offset - 1) * limit
	} else {
		offset = 0
		limit = 10
	}

	if sortBy == "" {
		sortBy = "created_at" // Default sort by created_at
	}
	if sortDir == "" {
		sortDir = "desc" // Default sort direction
	}

	payments, err := h.repo.GetPaged(r.Context(), sortBy, sortDir, offset, limit)
	if err != nil {
		http.Error(w, "Failed to get payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)

}
