import { selectedEntryRow } from "../keyboard/keyboard-entries.js";
import { initWriteCreateActions } from "./write-create.js";
import { initWriteModal } from "./write-modal.js";
import {
  closeContextMenu,
  contextMenuRow,
} from "./write-context-state.js";
import { initContextMenu } from "./write-context-menu.js";
import { runEntryWriteAction } from "./write-entry-actions.js";
import { initWriteMenu } from "./write-menu.js";
import { initWriteUploads } from "./write-upload.js";

const copyContextPath = (row, absolute, showToast) => {
  const path = absolute ? row?.dataset.absolutePath : row?.dataset.relativePath;
  if (!path || !navigator.clipboard) {
    showToast("Unable to copy path");

    return;
  }

  navigator.clipboard
    .writeText(path)
    .then(() =>
      showToast(absolute ? "Absolute path copied" : "Relative path copied"),
    )
    .catch(() => showToast("Unable to copy path"));
};

const initContextActions = (showToast) => {
  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    const writeAction = event.target.closest("[data-context-write]");
    const copyAction = event.target.closest("[data-context-copy]");
    if (!writeAction && !copyAction) {
      return;
    }

    event.preventDefault();
    event.stopPropagation();
    const row = contextMenuRow() || selectedEntryRow();
    if (writeAction instanceof HTMLButtonElement) {
      runEntryWriteAction(writeAction.dataset.contextWrite, row, showToast);
    } else if (copyAction instanceof HTMLButtonElement) {
      copyContextPath(
        row,
        copyAction.dataset.contextCopy === "absolute",
        showToast,
      );
    }
    closeContextMenu();
  });
};

export const initWriteActions = (showToast) => {
  initWriteModal();
  initContextMenu();
  initContextActions(showToast);
  initWriteMenu();
  initWriteCreateActions(showToast);
  initWriteUploads(showToast);
};
