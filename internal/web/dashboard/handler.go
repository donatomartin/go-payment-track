package dashboard

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"app/internal/invoice"
	"app/internal/payment"
	"app/web"
)

type DashboardHandler struct {
	logger         *log.Logger
	templates      *template.Template
	paymentService payment.PaymentService
	invoiceService invoice.InvoiceService
}

func NewDashboardHandler(

	paymentService payment.PaymentService,
	invoiceService invoice.InvoiceService,
	logger *log.Logger,

) *DashboardHandler {

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		log.Fatalf("failed to create template sub file system: %v", err)
	}
	t := template.Must(template.ParseFS(templateFS, "*.html"))

	return &DashboardHandler{
		logger:         logger,
		templates:      t,
		paymentService: paymentService,
		invoiceService: invoiceService,
	}
}

func (h *DashboardHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/dashboard", h.getDashboard)
	mux.HandleFunc("/", h.getDashboard) // Redirect root to dashboard
}

func (h *DashboardHandler) getDashboard(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/dashboard" {
		http.NotFound(w, r)
		return
	}

	payments, err := h.paymentService.GetPagedPayments(r.Context(), "date", "desc", 5, 6)
	if err != nil {
		http.Error(w, "Failed to get payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var paymentViews []PaymentView
	for _, payment := range payments {
		paymentViews = append(paymentViews, PaymentView{
			ID:         payment.ID,
			InvoiceID:  payment.InvoiceID,
			Amount:     fmt.Sprintf("%v €", payment.Amount),
			Date:       payment.Date.Format("2006-01-02"),
			ClientName: payment.ClientName,
		})
	}

	delayedInvoicesCount, err := h.invoiceService.GetDelayedInvoicesCount(r.Context())
	pendingInvoicesCount, err := h.invoiceService.GetPendingInvoicesCount(r.Context())

	data := struct {
		Title                string
		Payments             []PaymentView
		DelayedInvoicesCount int
		PendingInvoicesCount int
	}{
		Title:                "Dashboard",
		Payments:             paymentViews,
		DelayedInvoicesCount: delayedInvoicesCount,
		PendingInvoicesCount: pendingInvoicesCount,
	}

	if err := h.templates.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		http.Error(w, "Failed to render dashboard: ", http.StatusInternalServerError)
		h.logger.Fatal("Template error:", err)
	}

}

type PaymentView struct {
	ID         int
	InvoiceID  string
	Amount     string
	Date       string
	ClientName string
}
