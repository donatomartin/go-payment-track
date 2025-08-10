package templates

import (
	"app/internal/platform/util"
	"app/web"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

func (h *DashboardHandler) getDashboard(w http.ResponseWriter, r *http.Request) {

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		http.Error(w, "Failed to load templates: "+err.Error(), http.StatusInternalServerError)
		h.logger.Fatal("Template error:", err)
		return
	}

	t := template.Must(template.ParseFS(templateFS,
		"entrypoints/dashboard.html",
		"*.html",
	))

	if r.URL.Path != "/" && r.URL.Path != "/dashboard" {
		http.NotFound(w, r)
		return
	}

	pagination := Pagination{
		ShowSizeSelector: false,
		HtmxFragmentName: "payments",
		FirstPage:        1,
		PrevPage:         0,
		Page:             1,
		NextPage:         2,
		LastPage:         1000, // TODO: This should ideally be calculated based on total records,
		Size:             6,
		SortBy:           "created_at",
		SorDir:           "desc",
	}

	payments, err := h.paymentRepo.GetAll(r.Context(), "created_at", "desc", pagination.GetOffset(), pagination.Size)
	if err != nil {
		http.Error(w, "Failed to get payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	paymentViews := paymentsToPaymentViews(payments)

	currentDate := time.Now().Format("2006-01-02")
	currentYear := time.Now().Year()

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
		Pagination             Pagination
		Status                 string
		CurrentDate            string
		CurrentYear            int
	}{
		Title:                  "Pagos",
		Payments:               paymentViews,
		DelayedInvoicesCount:   delayedInvoicesCount,
		DelayedInvoicesAmount:  util.Float64ToEuros(delayedInvoicesAmount),
		PendingInvoicesCount:   pendingInvoicesCount,
		PendingInvoicesAmount:  util.Float64ToEuros(pendingInvoicesAmount),
		PartialInvoicesCount:   partialInvoicesCount,
		CompletedInvoicesCount: completedInvoicesCount,
		Pagination:             pagination,
		Status:                 "all",
		CurrentDate:            currentDate,
		CurrentYear:            currentYear,
	}

	if err := t.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		http.Error(w, "Failed to render dashboard: ", http.StatusInternalServerError)
		h.logger.Fatal("Template error:", err)
	}

}
