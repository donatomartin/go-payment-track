package dashboard

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

	invoices, err := h.invoiceRepo.GetPaged(r.Context(), "invoice_date", "desc", 0, 10)
	if err != nil {
		http.Error(w, "Failed to get invoices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var invoiceViews []InvoiceView
	for _, invoice := range invoices {
		invoiceViews = append(invoiceViews, InvoiceView{
			ID:           invoice.ID,
			CustomerName: invoice.CustomerName,
			AmountDue:    util.Float64ToEuros(invoice.AmountDue),
			PaymentMean:  invoice.PaymentMean,
			InvoiceDate:  invoice.InvoiceDate.Format("2006-01-02"),
			DueDate:      invoice.DueDate.Format("2006-01-02"),
		})
	}

	data := struct {
		Title    string
		Invoices []InvoiceView
	}{
		Title:    "Invoices",
		Invoices: invoiceViews,
	}

	if err := t.ExecuteTemplate(w, "invoices_table.html", data); err != nil {
		http.Error(w, "Failed to render invoices table: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

