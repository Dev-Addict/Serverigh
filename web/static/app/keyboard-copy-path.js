import { selectedEntryLink } from "./keyboard-entries.js";

export const copySelectedPath = (absolute, showToast) => {
  const row = selectedEntryLink()?.closest("[data-entry-row]");
  const path = absolute ? row?.dataset.absolutePath : row?.dataset.relativePath;
  if (!path || !navigator.clipboard) {
    return false;
  }

  navigator.clipboard
    .writeText(path)
    .then(() => {
      row.dataset.copied = "true";
      showToast(absolute ? "Absolute path copied" : "Relative path copied");
      window.setTimeout(() => {
        delete row.dataset.copied;
      }, 1200);
    })
    .catch(() => {
      showToast("Unable to copy path");
    });

  return true;
};
