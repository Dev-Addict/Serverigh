const selectionInputs = () => [
  ...document.querySelectorAll("[data-entry-select]"),
];

const filesPanel = () => document.querySelector(".file-panel");

const setRowSelected = (input) => {
  const row = input.closest("[data-entry-row]");
  if (!(row instanceof HTMLElement)) {
    return;
  }

  if (input.checked) {
    row.dataset.selected = "true";
  } else {
    delete row.dataset.selected;
  }
};

const selectedInputs = () => {
  return selectionInputs().filter((input) => input.checked);
};

export const selectedEntryRows = () => {
  return selectedInputs()
    .map((input) => input.closest("[data-entry-row]"))
    .filter((row) => row instanceof HTMLElement);
};

export const bulkSelectionForRow = (row) => {
  const rows = selectedEntryRows();
  const count = rows.length;

  return {
    count,
    rows,
    bulk: count > 0,
  };
};

export const syncBulkSelection = () => {
  selectionInputs().forEach(setRowSelected);
  const panel = filesPanel();
  if (!(panel instanceof HTMLElement)) {
    return;
  }

  const count = selectedInputs().length;
  panel.dataset.bulkCount = String(count);
  panel.dataset.hasBulkSelection = String(count > 0);
};

export const initBulkSelection = () => {
  syncBulkSelection();
  document.addEventListener("change", (event) => {
    if (
      event.target instanceof HTMLInputElement
      && event.target.matches("[data-entry-select]")
    ) {
      syncBulkSelection();
    }
  });
  document.body.addEventListener("htmx:afterSwap", syncBulkSelection);
};
