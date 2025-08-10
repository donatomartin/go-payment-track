package handler

import (
	"app/internal/invoice"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *ApiInvoiceHandler) addInvoice(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get from body
	bodyID := r.FormValue("invoice-id")
	bodyCustomerName := r.FormValue("customer-name")
	bodyAmount := r.FormValue("amount-due")
	bodyDate := r.FormValue("invoice-date")
	bodyPaymentMean := r.FormValue("payment-mean")
	bodyDueDate := r.FormValue("due-date")

	// Parse the body
	parsedID := strings.TrimSpace(bodyID)
	parsedCustomerName := strings.TrimSpace(bodyCustomerName)
	parsedPaymentMean := strings.TrimSpace(bodyPaymentMean)

	parsedAmount, err := strconv.ParseFloat(bodyAmount, 64)
	if err != nil {
		h.logger.Printf("Invalid amount: %v", err)
		w.Header().Set("Content-Type", "text/html")
		buf := `<div class="text-red-400">Monto inválido %s. Por favor, ingresa un número válido.</div>`
		w.Write([]byte(buf))
		return
	}
	parsedDate, err := time.Parse("2006-01-02", bodyDate)
	if err != nil {
		h.logger.Printf("Invalid date format: %v", err)
		w.Header().Set("Content-Type", "text/html")
		buf := `<div class="text-red-400">Formato de fecha inválido. Por favor, usa el formato AAAA-MM-DD.</div>`
		w.Write([]byte(buf))
		return
	}
	parsedDueDate, err := time.Parse("2006-01-02", bodyDueDate)
	if err != nil {
		h.logger.Printf("Invalid due date format: %v", err)
		w.Header().Set("Content-Type", "text/html")
		buf := `<div class="text-red-400">Formato de fecha de vencimiento inválido. Por favor, usa el formato AAAA-MM-DD.</div>`
		w.Write([]byte(buf))
		return
	}

	err = h.repo.AddInvoice(r.Context(), &invoice.Invoice{
		ID:           parsedID,
		CustomerName: parsedCustomerName,
		AmountDue:    parsedAmount,
		PaymentMean:  parsedPaymentMean,
		InvoiceDate:  parsedDate,
		DueDate:      parsedDueDate,
	})
	if err != nil {
		h.logger.Printf("Error adding payment: %v", err)
		w.Header().Set("Content-Type", "text/html")
		buf := fmt.Sprintf(`<div class="text-red-400">Error al agregar el factura. %v. Inténtalo de nuevo.</div>`, err.Error())
		w.Write([]byte(buf))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div class="text-green-400">factura agregada exitosamente.</div>`))

}
