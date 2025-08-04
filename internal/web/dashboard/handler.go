package dashboard

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	invoiceRepo "app/internal/invoice/repository"
	paymentRepo "app/internal/payment/repository"
	"app/internal/platform/util"
	"app/web"
)

type DashboardHandler struct {
	logger      *log.Logger
	templates   *template.Template
	paymentRepo paymentRepo.PaymentRepository
	invoiceRepo invoiceRepo.InvoiceRepository
}

func NewDashboardHandler(

	paymentRepository paymentRepo.PaymentRepository,
	invoiceRepository invoiceRepo.InvoiceRepository,
	logger *log.Logger,

) *DashboardHandler {

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		log.Fatalf("failed to create template sub file system: %v", err)
	}
	t := template.Must(template.ParseFS(templateFS, "*.html"))

	return &DashboardHandler{
		logger:      logger,
		templates:   t,
		paymentRepo: paymentRepository,
		invoiceRepo: invoiceRepository,
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

	payments, err := h.paymentRepo.GetPaged(r.Context(), "date", "desc", 5, 6)
	if err != nil {
		http.Error(w, "Failed to get payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var paymentViews []PaymentView
	for _, payment := range payments {
		paymentViews = append(paymentViews, PaymentView{
			ID:         payment.ID,
			InvoiceID:  payment.InvoiceID,
			Amount:     util.Float64ToEuros(payment.Amount),
			Date:       payment.Date.Format("2006-01-02"),
			ClientName: payment.ClientName,
		})
	}

	delayedInvoicesCount, err := h.invoiceRepo.GetDelayedInvoicesCount(r.Context())
	delayedInvoicesAmount, err := h.invoiceRepo.GetDelayedInvoicesAmount(r.Context())
	pendingInvoicesCount, err := h.invoiceRepo.GetPendingInvoicesCount(r.Context())
	pendingInvoicesAmount, err := h.invoiceRepo.GetPendingInvoicesAmount(r.Context())
	partialInvoicesCount, err := h.invoiceRepo.GetPartialInvoicesCount(r.Context())
	completedInvoicesCount, err := h.invoiceRepo.GetCompletedInvoicesCount(r.Context())

	data := struct {
		Title                  string
		Payments               []PaymentView
		DelayedInvoicesCount   int
		DelayedInvoicesAmount  string
		PendingInvoicesCount   int
		PendingInvoicesAmount  string
		PartialInvoicesCount   int
		CompletedInvoicesCount int
	}{
		Title:                  "Dashboard",
		Payments:               paymentViews,
		DelayedInvoicesCount:   delayedInvoicesCount,
		DelayedInvoicesAmount:  util.Float64ToEuros(delayedInvoicesAmount),
		PendingInvoicesCount:   pendingInvoicesCount,
		PendingInvoicesAmount:  util.Float64ToEuros(pendingInvoicesAmount),
		PartialInvoicesCount:   partialInvoicesCount,
		CompletedInvoicesCount: completedInvoicesCount,
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
