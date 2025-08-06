package templates

import (
	"app/web"
	"html/template"
	"io/fs"
	"net/http"
)

func (h *DashboardHandler) getPaymentsTable(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/payments" {
		http.NotFound(w, r)
		return
	}

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		http.Error(w, "Failed to load templates: "+err.Error(), http.StatusInternalServerError)
		h.logger.Fatal("Template error:", err)
		return
	}
	t := template.Must(template.ParseFS(templateFS,
		"entrypoints/payments_table.html",
		"*.html",
	))

	pagination := Pagination{
		ShowSizeSelector: true,
		FirstPage:      1,
		PrevPage:       0,
		Page:           1,
		NextPage:       2,
		LastPage:       1000, // TODO: This should ideally be calculated based on total records
		Size:           20,
		SortBy:         "created_at",
		SorDir:         "desc",
	}

	payments, err := h.paymentRepo.GetPaged(r.Context(), "created_at", "desc", pagination.GetOffset(), pagination.Size)
	if err != nil {
		http.Error(w, "Failed to get payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	paymentViews := paymentsToPaymentViews(payments)

	data := struct {
		Title      string
		Payments   []PaymentView
		Pagination Pagination
	}{
		Title:      "Payments",
		Payments:   paymentViews,
		Pagination: pagination,
	}

	if err := t.ExecuteTemplate(w, "payments_table.html", data); err != nil {
		http.Error(w, "Failed to render payments table: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
