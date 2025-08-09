package templates

import (
	"app/web"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
)

func (h *DashboardHandler) getInvoicesFragment(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/invoices/fragment" {
		http.NotFound(w, r)
		return
	}

	requestShowSizeSelector := r.URL.Query().Get("showSizeSelector")
	requestSortBy := r.URL.Query().Get("sortBy")
	requestSortDir := r.URL.Query().Get("sortDir")
	requestPage := r.URL.Query().Get("page")
	requestPageSize := r.URL.Query().Get("pageSize")

	paginationShowSizeSelector := requestShowSizeSelector == "true"
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
		ShowSizeSelector: paginationShowSizeSelector,
		HtmxFragmentName: "invoices",
		FirstPage:        1,
		PrevPage:         int(paginationPage - 1),
		Page:             int(paginationPage),
		NextPage:         int(paginationPage + 1),
		LastPage:         1000, // TODO: This should ideally be calculated based on total records
		Size:             int(paginationSize),
		SortBy:           paginationSortBy,
		SorDir:           paginationSortDir,
	}

	invoices, err := h.invoiceRepo.GetPaged(r.Context(), "created_at", "desc", pagination.GetOffset(), pagination.Size)
	if err != nil {
		http.Error(w, "Failed to get invoices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	invoiceViews := invoicesToInvoiceViews(invoices)

	data := struct {
		Title      string
		Invoices   []InvoiceView
		Pagination Pagination
	}{
		Title:      "Invoices",
		Invoices:   invoiceViews,
		Pagination: pagination,
	}

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		http.Error(w, "Failed to load templates: "+err.Error(), http.StatusInternalServerError)
		h.logger.Fatal("Template error:", err)
		return
	}
	t := template.Must(template.ParseFS(templateFS,
		"fragments/invoices.html",
		"*.html",
	))

	if err := t.ExecuteTemplate(w, "invoices.html", data); err != nil {
		http.Error(w, "Failed to render invoices table: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
