package templates

import (
	"app/internal/invoice"
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

	status := r.URL.Query().Get("status")

	requestShowSizeSelector := r.URL.Query().Get("showSizeSelector")
	requestSortBy := r.URL.Query().Get("sortBy")
	requestSortDir := r.URL.Query().Get("sortDir")
	requestPage := r.URL.Query().Get("page")
	requestPageSize := r.URL.Query().Get("pageSize")

	paginationShowSizeSelector := requestShowSizeSelector == "true"
	paginationSortBy := requestSortBy
	if paginationSortBy == "" {
		paginationSortBy = "created_at"
	}
	paginationSortDir := requestSortDir
	if paginationSortDir == "" {
		paginationSortDir = "desc"
	}
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

	var (
		invoices []invoice.Invoice
		title    string
	)

	switch status {
	case "completed":
		invoices, err = h.invoiceRepo.GetCompletedInvoices(r.Context(), pagination.GetOffset(), pagination.Size)
		title = "Facturas Completadas"
	case "pending":
		invoices, err = h.invoiceRepo.GetPendingInvoices(r.Context(), pagination.GetOffset(), pagination.Size)
		title = "Facturas Pendientes"
	case "delayed":
		invoices, err = h.invoiceRepo.GetDelayedInvoices(r.Context(), pagination.GetOffset(), pagination.Size)
		title = "Facturas Demoradas"
	case "partial":
		invoices, err = h.invoiceRepo.GetPartialInvoices(r.Context(), pagination.GetOffset(), pagination.Size)
		title = "Facturas en proceso de pago"
	default:
		invoices, err = h.invoiceRepo.GetAll(r.Context(), pagination.SortBy, pagination.SorDir, pagination.GetOffset(), pagination.Size)
		title = "Todas las facturas"
	}
	if err != nil {
		http.Error(w, "Failed to get invoices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	invoiceViews := invoicesToInvoiceViews(invoices)

	data := struct {
		Title      string
		Invoices   []InvoiceView
		Pagination Pagination
		Status     string
	}{
		Title:      title,
		Invoices:   invoiceViews,
		Pagination: pagination,
		Status:     status,
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
