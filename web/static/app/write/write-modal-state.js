import {
  destinationField,
  destinationInput,
  folderTree,
  setHidden,
  submitButton,
} from "./write-modal-elements.js";

export const setSubmitText = (value) => {
  const button = submitButton();
  if (button instanceof HTMLButtonElement) {
    button.textContent = value;
  }
};

export const setDangerState = (target, enabled) => {
  if (!(target instanceof HTMLElement)) {
    return;
  }
  if (enabled) {
    target.dataset.variant = "danger";

    return;
  }

  delete target.dataset.variant;
};

export const setDestinationState = (enabled) => {
  setHidden(destinationField(), !enabled);
  const input = destinationInput();
  if (input instanceof HTMLInputElement) {
    input.disabled = !enabled;
    input.value = "";
  }
  if (!enabled) {
    const tree = folderTree();
    if (tree instanceof HTMLElement) {
      tree.replaceChildren();
    }
  }
};
