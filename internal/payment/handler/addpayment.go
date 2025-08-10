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

	paymentAmount, _ := strconv.ParseFloat(bodyAmount, 64)
	paymentInvoiceID := strings.Trim(bodyInvoiceID, " \t\n\r")
	paymentDate, _ := time.Parse("2006-01-02", bodyDate)

	h.logger.Printf("Adding payment: InvoiceID=%s, Amount=%.2f, Date=%s", paymentInvoiceID, paymentAmount, paymentDate)

	err := h.repo.AddPayment(r.Context(), &payment.Payment{
		InvoiceID: paymentInvoiceID,
		Amount:    paymentAmount,
		Date:      paymentDate,
	})

	if err != nil {
		h.logger.Printf("Error adding payment: %v", err)
		w.Header().Set("Content-Type", "text/html")
		buf := fmt.Sprintf(`<div class="text-red-400">Error al agregar el pago. %v. Inténtalo de nuevo.</div>`, err.Error())
		w.Write([]byte(buf))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div class="text-green-400">Pago agregado exitosamente.</div>`))

}
