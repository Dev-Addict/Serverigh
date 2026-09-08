let previousSettingsFocus;

const settingsModal = () => document.querySelector("[data-settings-modal]");
const settingsDialog = () => document.querySelector(".settings-modal");

const openSettings = (trigger) => {
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

const closeSettings = () => {
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

export const initSettings = () => {
  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    const settingsOpen = event.target.closest("[data-settings-open]");
    if (settingsOpen instanceof HTMLElement) {
      event.preventDefault();
      openSettings(settingsOpen);

      return;
    }

    const settingsClose = event.target.closest("[data-settings-close]");
    if (settingsClose) {
      event.preventDefault();
      closeSettings();

      return;
    }

    if (event.target === settingsModal()) {
      closeSettings();
    }
  });

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeSettings();
    }
  });
};
