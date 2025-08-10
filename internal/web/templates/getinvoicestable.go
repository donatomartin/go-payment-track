package templates

import (
	"app/internal/invoice"
	"app/web"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

func (h *DashboardHandler) getInvoicesTable(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/invoices" {
		http.NotFound(w, r)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		status = "all" // Default to "all" if no status is provided
	}

	templateFS, err := fs.Sub(web.WebFS, "templates")
	if err != nil {
		http.Error(w, "Failed to load templates: "+err.Error(), http.StatusInternalServerError)
		h.logger.Fatal("Template error:", err)
		return
	}
	t := template.Must(template.ParseFS(templateFS,
		"entrypoints/invoices_table.html",
		"*.html",
	))

	pagination := Pagination{
		ShowSizeSelector: true,
		HtmxFragmentName: "invoices",
		FirstPage:        1,
		PrevPage:         0,
		Page:             1,
		NextPage:         2,
		LastPage:         1000, // TODO: This should ideally be calculated based on total records,
		Size:             20,
		SortBy:           "created_at",
		SorDir:           "desc",
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
		title = "Facturas en Proceso"
	default:
		invoices, err = h.invoiceRepo.GetAll(r.Context(), "invoice_date", "desc", pagination.GetOffset(), pagination.Size)
		title = "Todas las facturas"
	}
	if err != nil {
		http.Error(w, "Failed to get invoices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	invoiceViews := invoicesToInvoiceViews(invoices)

	currentDate := time.Now().Format("2006-01-02")
	currentYear := time.Now().Year()

	data := struct {
		Title       string
		Invoices    []InvoiceView
		Pagination  Pagination
		Status      string
		CurrentDate string
		CurrentYear int
	}{
		Title:       title,
		Invoices:    invoiceViews,
		Pagination:  pagination,
		Status:      status,
		CurrentDate: currentDate,
		CurrentYear: currentYear,
	}

	if err := t.ExecuteTemplate(w, "invoices_table.html", data); err != nil {
		http.Error(w, "Failed to render invoices table: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
