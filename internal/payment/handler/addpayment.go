package handler

import (
	"app/internal/payment"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *ApiPaymentHandler) addPayment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get from body
	bodyAmount := r.FormValue("amount")
	bodyInvoiceID := r.FormValue("invoice-id")
	bodyDate := r.FormValue("date")

	// Parse the body
	parsedAmount, err := strconv.ParseFloat(bodyAmount, 64)
	if err != nil {
		h.logger.Printf("Invalid amount: %v", err)
		w.Header().Set("Content-Type", "text/html")
		buf := `<div class="text-red-400">Cantidad inválida. Por favor, ingresa un número válido.</div>`
		w.Write([]byte(buf))
		return
	}
	parsedInvoiceID := strings.Trim(bodyInvoiceID, " \t\n\r")
	parsedDate, err := time.Parse("2006-01-02", bodyDate)
	if err != nil {
		h.logger.Printf("Invalid date format: %v", err)
		w.Header().Set("Content-Type", "text/html")
		buf := `<div class="text-red-400">Formato de fecha inválido. Por favor, usa el formato AAAA-MM-DD.</div>`
		w.Write([]byte(buf))
		return
	}

	// Add the payment
	err = h.repo.AddPayment(r.Context(), &payment.Payment{
		InvoiceID: parsedInvoiceID,
		Amount:    parsedAmount,
		Date:      parsedDate,
	})
	if err != nil {
		h.logger.Printf("Error adding payment: %v", err)
		w.Header().Set("Content-Type", "text/html")
		buf := fmt.Sprintf(`<div class="text-red-400">Error al agregar el pago. %v. Inténtalo de nuevo.</div>`, err.Error())
		w.Write([]byte(buf))
		return
	}

	h.logger.Printf("Payment added successfully for invoice ID: %s", parsedInvoiceID)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div class="text-green-400">Pago agregado exitosamente.</div>`))

}
