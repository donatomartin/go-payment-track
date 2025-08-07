package templates

import (
	"app/internal/platform/util"
	"app/web"
	"html/template"
	"io/fs"
	"net/http"
)

func (h *DashboardHandler) getInvoicesTable(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/invoices" {
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
		"entrypoints/invoices_table.html",
		"*.html",
	))

	pagination := Pagination{
		ShowSizeSelector: false,
		FirstPage:        1,
		PrevPage:         0,
		Page:             1,
		NextPage:         2,
		LastPage:         1000, // TODO: This should ideally be calculated based on total records,
		Size:             6,
		SortBy:           "created_at",
		SorDir:           "desc",
	}

	invoices, err := h.invoiceRepo.GetPaged(r.Context(), "invoice_date", "desc", pagination.GetOffset(), pagination.Size)
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

	if err := t.ExecuteTemplate(w, "invoices_table.html", data); err != nil {
		http.Error(w, "Failed to render invoices table: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
