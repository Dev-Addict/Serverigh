import { baseActionValues, postWriteAction } from "./write-request.js";
import {
  openConfirmModal,
  openDestinationModal,
  openNameModal,
} from "./write-modal.js";

const actionValues = (action, row, formValues = {}) => {
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

const submitAction = (action, row, formValues, showToast) => {
  if (!postWriteAction(action, actionValues(action, row, formValues))) {
    showToast("Write action unavailable");

    return;
  }

  showToast("Write action submitted");
};

const openRename = (row, showToast) => {
  openNameModal({
    nameLabel: "Name",
    nameValue: row.dataset.entryName || "",
    submit: (values) => submitAction("rename", row, values, showToast),
    submitLabel: "Rename",
    title: "Rename entry",
  });
};

const openDestinationAction = (action, row, showToast) => {
  const label = action === "copy" ? "Copy" : "Move";
  const current = baseActionValues().path;
  openDestinationModal({
    selectedPath: current,
    submit: (values) => submitAction(action, row, values, showToast),
    submitLabel: label,
    summary: `${label} ${row.dataset.entryName || "this entry"} to:`,
    title: `${label} entry`,
  });
};

const openDelete = (row, showToast) => {
  openConfirmModal({
    submit: () => submitAction("delete", row, {}, showToast),
    submitLabel: "Trash",
    summary: `Trash ${row.dataset.entryName || "this entry"}?`,
    title: "Trash entry",
  });
};

const runDirectAction = (action, row, showToast) => {
  submitAction(action, row, {}, showToast);
};

export const runEntryWriteAction = (action, row, showToast) => {
  if (!document.querySelector("[data-write-menu]")) {
    showToast("Write mode is disabled");

    return true;
  }

  if (!(row instanceof HTMLElement)) {
    showToast("Select an entry first");

    return true;
  }

  switch (action) {
  case "rename":
    openRename(row, showToast);
    break;
  case "copy":
  case "move":
    openDestinationAction(action, row, showToast);
    break;
  case "delete":
    openDelete(row, showToast);
    break;
  case "duplicate":
    runDirectAction(action, row, showToast);
    break;
  }

  return true;
};
