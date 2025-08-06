package templates

import (
	"app/web"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
)

func (h *DashboardHandler) getPaymentsFragment(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/payments/fragment" {
		http.NotFound(w, r)
		return
	}

	requestShowPagination := r.URL.Query().Get("showPagination")
	requestSortBy := r.URL.Query().Get("sortBy")
	requestSortDir := r.URL.Query().Get("sortDir")
	requestPage := r.URL.Query().Get("page")
	requestPageSize := r.URL.Query().Get("pageSize")

	paginationShowPagination := requestShowPagination == "true"
	paginationSortBy := requestSortBy
	paginationSortDir := requestSortDir
	paginationPage, err := strconv.ParseInt(requestPage, 10, 64)
	if err != nil || paginationPage < 1 {
		paginationPage = 1
	}
	paginationSize, err := strconv.ParseInt(requestPageSize, 10, 64)
	if err != nil || paginationSize < 1 {
		paginationSize = 10
	}

	pagination := Pagination{
		ShowSizeSelector: paginationShowPagination,
		FirstPage:        1,
		PrevPage:         int(paginationPage - 1),
		Page:             int(paginationPage),
		NextPage:         int(paginationPage + 1),
		LastPage:         1000, // TODO: This should ideally be calculated based on total records
		Size:             int(paginationSize),
		SortBy:           paginationSortBy,
		SorDir:           paginationSortDir,
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

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		http.Error(w, "Failed to load templates: "+err.Error(), http.StatusInternalServerError)
		h.logger.Fatal("Template error:", err)
		return
	}
	t := template.Must(template.ParseFS(templateFS,
		"fragments/payments.html",
		"*.html",
	))

	if err := t.ExecuteTemplate(w, "payments.html", data); err != nil {
		http.Error(w, "Failed to render payments table: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
