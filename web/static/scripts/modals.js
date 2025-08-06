// Payment Modal
const paymentModal = document.getElementById("payment-modal");

let addPaymentButtons = [];

addPaymentButtons.push(document.querySelector("button:nth-child(2)"));
addPaymentButtons.push(
  document.querySelector(".quick-access button:nth-child(2)"),
);

const closePaymentModal = document.getElementById("close-payment-modal");

addPaymentButtons.forEach((button) => {
  button.addEventListener("click", () => {
    paymentModal.classList.remove("hidden");
    paymentModal.classList.add("flex");
  });
});

closePaymentModal.addEventListener("click", () => {
  paymentModal.classList.add("hidden");
  paymentModal.classList.remove("flex");
});

// Invoice Modal
const invoiceModal = document.getElementById("invoice-modal");
const addInvoiceButtons = [];
addInvoiceButtons.push(document.querySelector("button:nth-child(3)"));
addInvoiceButtons.push(
  document.querySelector(".quick-access button:nth-child(1)"),
);
const closeInvoiceModal = document.getElementById("close-invoice-modal");

addInvoiceButtons.forEach((button) => {
  button.addEventListener("click", () => {
    paymentModal.classList.remove("hidden");
    paymentModal.classList.add("flex");
  });
});

closeInvoiceModal.addEventListener("click", () => {
  paymentModal.classList.add("hidden");
  paymentModal.classList.remove("flex");
});
