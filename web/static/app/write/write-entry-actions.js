import { baseActionValues } from "./write-request.js";
import {
  openConfirmModal,
  openDestinationModal,
  openNameModal,
} from "./write-modal.js";
import {
  submitEntryAction,
  writeControlsAvailable,
} from "./write-entry-submit.js";
import {
  openDeleteTrash,
  openRestoreTrash,
} from "./write-entry-trash-actions.js";

const openRename = (row, showToast) => {
  openNameModal({
    nameLabel: "Name",
    nameValue: row.dataset.entryName || "",
    submit: (values) => submitEntryAction("rename", row, values, showToast),
    submitLabel: "Rename",
    title: "Rename entry",
  });
};

const openDestinationAction = (action, row, showToast) => {
  const label = action === "copy" ? "Copy" : "Move";
  const current = baseActionValues().path;
  const selectedPath = current === "/trash" ? "/" : current;
  openDestinationModal({
    selectedPath,
    submit: (values) => submitEntryAction(action, row, values, showToast),
    submitLabel: label,
    summary: `${label} ${row.dataset.entryName || "this entry"} to:`,
    title: `${label} entry`,
  });
};

const openDelete = (row, showToast) => {
	openConfirmModal({
		danger: true,
		submit: () => submitEntryAction("delete", row, {}, showToast),
		submitLabel: "Trash",
    summary: `Trash ${row.dataset.entryName || "this entry"}?`,
    title: "Trash entry",
  });
};

const runDirectAction = (action, row, showToast) => {
  submitEntryAction(action, row, {}, showToast);
};

export const runEntryWriteAction = (action, row, showToast) => {
  if (!writeControlsAvailable()) {
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
  case "trash/restore":
    openRestoreTrash(row, showToast);
    break;
  case "trash/delete":
    openDeleteTrash(row, showToast);
    break;
  case "duplicate":
    runDirectAction(action, row, showToast);
    break;
  }

  return true;
};
