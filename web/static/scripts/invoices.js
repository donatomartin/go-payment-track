let editInvoiceModal;
let editInvoiceForm;

document.addEventListener("DOMContentLoaded", () => {
  editInvoiceModal = document.getElementById("edit-invoice-modal");
  editInvoiceForm = document.getElementById("invoice-edit-form");
  bindInvoicesTable();
});

document.body.addEventListener("htmx:afterSwap", (evt) => {
  if (evt.target.id === "invoices-fragment") {
    bindInvoicesTable();
  }
});

document.body.addEventListener("htmx:afterRequest", (evt) => {
  if (evt.target.id === "invoice-edit-form" && evt.detail.xhr.status === 200) {
    editInvoiceModal.classList.add("hidden");
    editInvoiceModal.classList.remove("flex");
    editInvoiceForm.reset();
    const msg = document.getElementById("edit-invoice-form-messages");
    if (msg) {
      msg.innerHTML = "";
    }
    htmx.trigger(document.body, "invoicesUpdated");
  }
});

function bindInvoicesTable() {
  const rows = document.querySelectorAll("#invoices-fragment tr[data-invoice-id]");
  rows.forEach((row) => {
    row.addEventListener("dblclick", () => {
      document.getElementById("edit-invoice-id").value = row.getAttribute("data-invoice-id") || "";
      document.getElementById("edit-customer-name").value = row.getAttribute("data-customer-name") || "";
      document.getElementById("edit-amount-due").value = row.getAttribute("data-amount-due") || "";
      document.getElementById("edit-payment-mean").value = row.getAttribute("data-payment-mean") || "";
      document.getElementById("edit-invoice-date").value = row.getAttribute("data-invoice-date") || "";
      document.getElementById("edit-due-date").value = row.getAttribute("data-due-date") || "";
      const msg = document.getElementById("edit-invoice-form-messages");
      if (msg) {
        msg.innerHTML = "";
      }
      editInvoiceModal.classList.remove("hidden");
      editInvoiceModal.classList.add("flex");
    });
  });
}
