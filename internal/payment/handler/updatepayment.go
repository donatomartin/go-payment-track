package handler

import (
	"app/internal/payment"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *ApiPaymentHandler) updatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	invoiceID := strings.TrimSpace(r.FormValue("invoice-id"))
	amountStr := r.FormValue("amount")
	dateStr := r.FormValue("date")

	id, err := strconv.Atoi(idStr)
	if err != nil {
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

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="text-red-400">Formato de fecha inválido. Por favor, usa el formato AAAA-MM-DD.</div>`))
		return
	}

	err = h.repo.UpdatePayment(r.Context(), &payment.Payment{
		ID:        id,
		InvoiceID: invoiceID,
		Amount:    amount,
		Date:      parsedDate,
	})
	if err != nil {
		h.logger.Printf("Error updating payment %d: %v", id, err)
		w.Header().Set("Content-Type", "text/html")
		buf := fmt.Sprintf(`<div class="text-red-400">Error al actualizar el pago. %v. Inténtalo de nuevo.</div>`, err.Error())
		w.Write([]byte(buf))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div class="text-green-400">Pago actualizado exitosamente.</div>`))
}
