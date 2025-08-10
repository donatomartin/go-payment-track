package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"app/internal/admin/data"
)

func (h *ApiAdminHandler) exportData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		http.Error(w, "missing year", http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}

	f, err := data.ExportData(r.Context(), h.paymentRepo, h.invoiceRepo, year)
	if err != nil {
		h.logger.Printf("export data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=data-%d.xlsx", year))
	if err := f.Write(w); err != nil {
		h.logger.Printf("write excel: %v", err)
	}
}
