import { baseActionValues, postWriteAction } from "./write-request.js";

export const entryActionValues = (action, row, formValues = {}) => {
  const values = baseActionValues(row.dataset.entryPath || "");
  switch (action) {
  case "rename":
    values.name = formValues.name;
    break;
  case "copy":
  case "move":
    values.destination = formValues.destination || "/";
    break;
  }

  return values;
};

export const submitEntryAction = (action, row, formValues, showToast) => {
  if (!postWriteAction(action, entryActionValues(action, row, formValues))) {
    showToast("Write action unavailable");

    return;
  }

  showToast("Write action submitted");
};

export const writeControlsAvailable = () => {
  return Boolean(
    document.querySelector("[data-write-menu]") ||
    document.querySelector(".file-panel[data-trash-view='true']"),
  );
};
