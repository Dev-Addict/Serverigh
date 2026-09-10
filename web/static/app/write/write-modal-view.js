import {
  destinationField,
  dialog,
  form,
  modal,
  nameField,
  nameInput,
  nameLabel,
  setHidden,
  setText,
  submitButton,
  summary,
  title,
} from "./write-modal-elements.js";

let previousFocus = null;

const setSubmitText = (value) => {
  const button = submitButton();
  if (button instanceof HTMLButtonElement) {
    button.textContent = value;
  }
};

const focusModal = () => {
  const input = nameInput();
  const field = nameField();
  if (
    input instanceof HTMLInputElement
    && field instanceof HTMLElement
    && !field.hidden
  ) {
    input.focus();
    input.select();

    return;
  }

  dialog()?.focus();
};

export const openModalView = (options) => {
  const target = modal();
  const targetForm = form();
  if (
    !(target instanceof HTMLElement)
    || !(targetForm instanceof HTMLFormElement)
  ) {
    return false;
  }

  previousFocus = document.activeElement;
  targetForm.reset();
  setText(title(), options.title);
  setText(summary(), options.summary || "");
  setText(nameLabel(), options.nameLabel || "Name");
  setSubmitText(options.submitLabel || "Submit");
  setHidden(nameField(), !options.name);
  setHidden(destinationField(), !options.destination);

  const input = nameInput();
  if (input instanceof HTMLInputElement) {
    input.value = options.nameValue || "";
    input.required = Boolean(options.name);
  }

  target.dataset.open = "true";
  target.setAttribute("aria-hidden", "false");
  focusModal();

  return true;
};

export const closeModalView = () => {
  const target = modal();
  if (!(target instanceof HTMLElement)) {
    return;
  }

  delete target.dataset.open;
  target.setAttribute("aria-hidden", "true");
  if (previousFocus instanceof HTMLElement) {
    previousFocus.focus();
  }
};
