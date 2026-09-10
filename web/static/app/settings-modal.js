import { settingsModal } from "./settings-dom.js";

let previousSettingsFocus;

const settingsDialog = () => {
  return document.querySelector(".settings-modal");
};

export const openSettings = (trigger) => {
  const modal = settingsModal();
  const dialog = settingsDialog();
  if (!(modal instanceof HTMLElement) || !(dialog instanceof HTMLElement)) {
    return;
  }

  previousSettingsFocus = trigger;
  modal.dataset.open = "true";
  modal.setAttribute("aria-hidden", "false");
  dialog.focus();
};

export const closeSettings = () => {
  const modal = settingsModal();
  if (!(modal instanceof HTMLElement)) {
    return;
  }

  delete modal.dataset.open;
  modal.setAttribute("aria-hidden", "true");
  if (previousSettingsFocus instanceof HTMLElement) {
    previousSettingsFocus.focus();
  }
};
