import {
  destinationInput,
  form,
  modal,
  nameInput,
} from "./write-modal-elements.js";
import {
  clearFolderTreeCache,
  loadFolderTree,
  selectFolder,
  toggleFolder,
} from "./write-modal-folders.js";
import {
  closeModalView,
  openModalView,
} from "./write-modal-view.js";

let activeModal = null;

const openModal = (options) => {
  if (!openModalView(options)) {
    return false;
  }
  activeModal = options;

  return true;
};

export const closeWriteModal = () => {
  activeModal = null;
  clearFolderTreeCache();
  closeModalView();
};

export const openNameModal = (options) => {
  openModal({
    ...options,
    destination: false,
    name: true,
  });
};

export const openConfirmModal = (options) => {
  openModal({
    ...options,
    destination: false,
    name: false,
  });
};

export const openDestinationModal = (options) => {
  const opened = openModal({
    ...options,
    destination: true,
    name: false,
  });
  if (opened) {
    loadFolderTree(options.selectedPath);
  }
};

const submitModal = (event) => {
  event.preventDefault();
  if (!activeModal) {
    return;
  }

  const input = nameInput();
  const destination = destinationInput();
  activeModal.submit({
    destination: destination instanceof HTMLInputElement ? destination.value : "",
    name: input instanceof HTMLInputElement ? input.value : "",
  });
  closeWriteModal();
};

export const initWriteModal = () => {
  form()?.addEventListener("submit", submitModal);
  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    if (event.target.closest("[data-write-modal-close]")) {
      event.preventDefault();
      closeWriteModal();

      return;
    }

    const toggle = event.target.closest("[data-folder-toggle]");
    const folder = event.target.closest("[data-folder-select]");
    if (toggle instanceof HTMLButtonElement) {
      event.preventDefault();
      const row = toggle.closest("[data-folder-row]");
      if (row instanceof HTMLElement) {
        toggleFolder(row);
      }
    } else if (folder instanceof HTMLButtonElement) {
      event.preventDefault();
      const row = folder.closest("[data-folder-row]");
      if (row instanceof HTMLElement) {
        selectFolder(row);
      }
    } else if (event.target === modal()) {
      closeWriteModal();
    }
  });
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeWriteModal();
    }
  });
};
