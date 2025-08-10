// Payment Modal
const paymentModal = document.getElementById("payment-modal");

let addPaymentButtons = [];

addPaymentButtons.push(document.querySelector("button:nth-child(4)"));
addPaymentButtons.push(
  document.querySelector(".quick-access button:nth-child(2)"),
);

const closePaymentModal = document.getElementById("close-payment-modal");

addPaymentButtons.forEach((button) => {
  if (button === null) {
    return;
  }
  button.addEventListener("click", () => {
    paymentModal.classList.remove("hidden");
    paymentModal.classList.add("flex");
  });
});

  closePaymentModal.addEventListener("click", () => {
    paymentModal.classList.add("hidden");
    paymentModal.classList.remove("flex");
  });

// Edit Payment Modal
const closeEditPaymentModal = document.getElementById("close-edit-payment-modal");
if (closeEditPaymentModal) {
  const editPaymentModalEl = document.getElementById("edit-payment-modal");
  closeEditPaymentModal.addEventListener("click", () => {
    editPaymentModalEl.classList.add("hidden");
    editPaymentModalEl.classList.remove("flex");
  });
}

// Invoice Modal
const invoiceModal = document.getElementById("invoice-modal");
const addInvoiceButtons = [];
addInvoiceButtons.push(document.querySelector("button:nth-child(3)"));
addInvoiceButtons.push(
  document.querySelector(".quick-access button:nth-child(1)"),
);
const closeInvoiceModal = document.getElementById("close-invoice-modal");

addInvoiceButtons.forEach((button) => {
  if (button === null) {
    return;
  }
  button.addEventListener("click", () => {
    invoiceModal.classList.remove("hidden");
    invoiceModal.classList.add("flex");
  });
});

closeInvoiceModal.addEventListener("click", () => {
  invoiceModal.classList.add("hidden");
  invoiceModal.classList.remove("flex");
});

// Edit Invoice Modal
const closeEditInvoiceModal = document.getElementById("close-edit-invoice-modal");
if (closeEditInvoiceModal) {
  const editInvoiceModalEl = document.getElementById("edit-invoice-modal");
  closeEditInvoiceModal.addEventListener("click", () => {
    editInvoiceModalEl.classList.add("hidden");
    editInvoiceModalEl.classList.remove("flex");
  });
}

// Settings Modal
const settingsModal = document.getElementById("settings-modal");
const openSettingsModal = document.getElementById("open-settings-modal");
const closeSettingsModal = document.getElementById("close-settings-modal");

if (openSettingsModal) {
  openSettingsModal.addEventListener("click", () => {
    settingsModal.classList.remove("hidden");
    settingsModal.classList.add("flex");
  });
}

if (closeSettingsModal) {
  closeSettingsModal.addEventListener("click", () => {
    settingsModal.classList.add("hidden");
    settingsModal.classList.remove("flex");
  });
}

const importFile = document.getElementById("import-file");
const importDataButton = document.querySelector(".quick-access button:nth-child(3)");
if (importDataButton && importFile) {
  importDataButton.addEventListener("click", () => importFile.click());
}

const exportModal = document.getElementById("export-modal");
const openExportModal = document.querySelector(".quick-access button:nth-child(4)");
const closeExportModal = document.getElementById("close-export-modal");

if (openExportModal && exportModal) {
  openExportModal.addEventListener("click", () => {
    exportModal.classList.remove("hidden");
    exportModal.classList.add("flex");
  });
}

if (closeExportModal && exportModal) {
  closeExportModal.addEventListener("click", () => {
    exportModal.classList.add("hidden");
    exportModal.classList.remove("flex");
  });
}
