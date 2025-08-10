package handler

import (
	"net/http"

	"github.com/xuri/excelize/v2"

	"app/internal/admin/data"
)

func (h *ApiAdminHandler) importData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "failed to parse file", http.StatusBadRequest)
		return
	}

	if err := data.ImportFacturas(r.Context(), f, h.invoiceRepo); err != nil {
		h.logger.Printf("import facturas: %v", err)
		http.Error(w, "failed to import invoices", http.StatusInternalServerError)
		return
	}

	if err := data.ImportPagos(r.Context(), f, h.paymentRepo); err != nil {
		h.logger.Printf("import pagos: %v", err)
		http.Error(w, "failed to import payments", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("import completed"))
}
