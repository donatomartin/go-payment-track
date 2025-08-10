let editPaymentModal;
let editPaymentForm;

document.addEventListener("DOMContentLoaded", () => {
  editPaymentModal = document.getElementById("edit-payment-modal");
  editPaymentForm = document.getElementById("payment-edit-form");
  bindPaymentsTable();
});

document.body.addEventListener("htmx:afterSwap", (evt) => {
  if (evt.target.id === "payments-fragment") {
    bindPaymentsTable();
  }
});

document.body.addEventListener("htmx:afterRequest", (evt) => {
  if (evt.target.id === "payment-edit-form" && evt.detail.xhr.status === 200) {
    editPaymentModal.classList.add("hidden");
    editPaymentModal.classList.remove("flex");
    editPaymentForm.reset();
    const msg = document.getElementById("edit-payment-form-messages");
    if (msg) {
      msg.innerHTML = "";
    }
    htmx.trigger(document.body, "paymentsUpdated");
  }
});

function bindPaymentsTable() {
  const rows = document.querySelectorAll("#payments-fragment tr[data-payment-id]");
  rows.forEach((row) => {
    row.addEventListener("dblclick", () => {
      document.getElementById("edit-id").value = row.getAttribute("data-payment-id") || "";
      document.getElementById("edit-invoice-id").value = row.getAttribute("data-invoice-id") || "";
      document.getElementById("edit-amount").value = row.getAttribute("data-amount") || "";
      document.getElementById("edit-date").value = row.getAttribute("data-date") || "";

      const msg = document.getElementById("edit-payment-form-messages");
      if (msg) {
        msg.innerHTML = "";
      }
      editPaymentModal.classList.remove("hidden");
      editPaymentModal.classList.add("flex");
    });
  });
}
