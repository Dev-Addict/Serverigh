import { openConfirmModal } from "./write-modal.js";
import { submitEntryAction } from "./write-entry-submit.js";

export const openRestoreTrash = (row, showToast) => {
  openConfirmModal({
    submit: () => submitEntryAction("trash/restore", row, {}, showToast),
    submitLabel: "Restore",
    summary: `Restore ${row.dataset.entryName || "this entry"}?`,
    title: "Restore entry",
  });
};

export const openDeleteTrash = (row, showToast) => {
	openConfirmModal({
		danger: true,
		submit: () => submitEntryAction("trash/delete", row, {}, showToast),
    submitLabel: "Delete permanently",
    summary: `Permanently delete ${row.dataset.entryName || "this entry"}?`,
    title: "Delete permanently",
  });
};
