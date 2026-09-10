export const modal = () => document.querySelector("[data-write-modal]");
export const dialog = () => document.querySelector(".write-modal");
export const form = () => document.querySelector("[data-write-modal-form]");
export const title = () => document.querySelector("[data-write-modal-title]");
export const summary = () => document.querySelector("[data-write-modal-summary]");
export const nameField = () => document.querySelector("[data-write-name-field]");
export const nameLabel = () => document.querySelector("[data-write-name-label]");
export const nameInput = () => document.querySelector("[data-write-name-input]");
export const destinationField = () =>
  document.querySelector("[data-write-destination-field]");
export const destinationInput = () =>
  document.querySelector("[data-write-destination-input]");
export const folderTree = () => document.querySelector("[data-folder-tree]");
export const submitButton = () =>
  document.querySelector("[data-write-modal-submit]");

export const setHidden = (element, hidden) => {
  if (element instanceof HTMLElement) {
    element.hidden = hidden;
  }
};

export const setText = (element, value) => {
  if (element instanceof HTMLElement) {
    element.textContent = value;
  }
};
