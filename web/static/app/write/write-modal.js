import {
  destinationInput,
  form,
  nameInput,
} from "./write-modal-elements.js";
import {
  clearFolderTreeCache,
  loadFolderTree,
} from "./write-modal-folders.js";
import { initWriteModalEvents } from "./write-modal-events.js";
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
  initWriteModalEvents(closeWriteModal);
};
