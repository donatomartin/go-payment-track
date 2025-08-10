package handler

import (
	"app/internal/invoice"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *ApiInvoiceHandler) updateInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	customerName := strings.TrimSpace(r.FormValue("customer-name"))
	amountStr := r.FormValue("amount-due")
	paymentMean := strings.TrimSpace(r.FormValue("payment-mean"))
	invoiceDateStr := r.FormValue("invoice-date")
	dueDateStr := r.FormValue("due-date")

	if id == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="text-red-400">ID inválido.</div>`))
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="text-red-400">Cantidad inválida. Por favor, ingresa un número válido.</div>`))
		return
	}

	invoiceDate, err := time.Parse("2006-01-02", invoiceDateStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="text-red-400">Formato de fecha de factura inválido. Por favor, usa el formato AAAA-MM-DD.</div>`))
		return
	}

	dueDate, err := time.Parse("2006-01-02", dueDateStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="text-red-400">Formato de fecha de vencimiento inválido. Por favor, usa el formato AAAA-MM-DD.</div>`))
		return
	}

	err = h.repo.UpdateInvoice(r.Context(), &invoice.Invoice{
		ID:           id,
		CustomerName: customerName,
		AmountDue:    amount,
		PaymentMean:  paymentMean,
		InvoiceDate:  invoiceDate,
		DueDate:      dueDate,
	})
	if err != nil {
		h.logger.Printf("Error updating invoice %s: %v", id, err)
		w.Header().Set("Content-Type", "text/html")
		buf := fmt.Sprintf(`<div class="text-red-400">Error al actualizar la factura. %v. Inténtalo de nuevo.</div>`, err.Error())
		w.Write([]byte(buf))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div class="text-green-400">Factura actualizada exitosamente.</div>`))
}
