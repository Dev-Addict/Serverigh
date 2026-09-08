let previousHelpFocus;

export const openHelpModal = (trigger) => {
  const modal = document.querySelector("[data-help-modal]");
  const dialog = document.querySelector(".help-modal");
  if (!(modal instanceof HTMLElement) || !(dialog instanceof HTMLElement)) {
    return false;
  }

  previousHelpFocus = trigger;
  modal.dataset.open = "true";
  modal.setAttribute("aria-hidden", "false");
  dialog.focus();

  return true;
};

const closeHelpModal = () => {
  const modal = document.querySelector("[data-help-modal]");
  if (!(modal instanceof HTMLElement)) {
    return;
  }

  delete modal.dataset.open;
  modal.setAttribute("aria-hidden", "true");
  if (previousHelpFocus instanceof HTMLElement) {
    previousHelpFocus.focus();
  }
};

export const initHelpModal = () => {
  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    const helpOpen = event.target.closest("[data-help-open]");
    if (helpOpen instanceof HTMLElement) {
      event.preventDefault();
      openHelpModal(helpOpen);

      return;
    }

    if (event.target.closest("[data-help-close]")) {
      event.preventDefault();
      closeHelpModal();

      return;
    }

    if (event.target === document.querySelector("[data-help-modal]")) {
      closeHelpModal();
    }
  });

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeHelpModal();
    }
  });
};
