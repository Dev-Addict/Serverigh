import { modal } from "./write-modal-elements.js";
import {
  selectFolder,
  toggleFolder,
} from "./write-modal-folders.js";

const handleFolderAction = (event) => {
  const toggle = event.target.closest("[data-folder-toggle]");
  const folder = event.target.closest("[data-folder-select]");
  if (toggle instanceof HTMLButtonElement) {
    event.preventDefault();
    const row = toggle.closest("[data-folder-row]");
    if (row instanceof HTMLElement) {
      toggleFolder(row);
    }

    return true;
  }
  if (folder instanceof HTMLButtonElement) {
    event.preventDefault();
    const row = folder.closest("[data-folder-row]");
    if (row instanceof HTMLElement) {
      selectFolder(row);
    }

    return true;
  }

  return false;
};

export const initWriteModalEvents = (closeWriteModal) => {
  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    if (event.target.closest("[data-write-modal-close]")) {
      event.preventDefault();
      closeWriteModal();

      return;
    }
    if (handleFolderAction(event)) {
      return;
    }
    if (event.target === modal()) {
      closeWriteModal();
    }
  });
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeWriteModal();
    }
  });
};
