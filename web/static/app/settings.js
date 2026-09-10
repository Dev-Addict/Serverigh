import { debounce } from "./debounce.js";
import { flushSettings, saveSettings } from "./settings-save.js";
import {
  settingsModal,
  syncColumnToggle,
  syncColumnToggles,
  syncMaxPreviewBytes,
} from "./settings-dom.js";
import { closeSettings, openSettings } from "./settings-modal.js";

const saveSettingsSoon = debounce(saveSettings, 250);

const handleSettingsClick = (event) => {
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
};

const handleSettingsChange = (event) => {
  if (
    event.target instanceof HTMLInputElement
    && event.target.matches("[data-column-toggle]")
  ) {
    syncColumnToggle(event.target);
    saveSettings();

    return;
  }

  if (
    event.target instanceof HTMLInputElement
    && event.target.matches("[data-max-preview-bytes]")
  ) {
    syncMaxPreviewBytes(event.target);
    saveSettings();

    return;
  }

  if (
    event.target instanceof HTMLSelectElement
    && event.target.matches("[data-theme-select]")
  ) {
    saveSettings();
  }
};

const handleSettingsInput = (event) => {
  if (
    event.target instanceof HTMLInputElement
    && event.target.matches("[data-max-preview-bytes]")
  ) {
    syncMaxPreviewBytes(event.target);
    saveSettingsSoon();
  }
};

export const initSettings = () => {
  syncColumnToggles();

  document.addEventListener("click", handleSettingsClick);

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeSettings();
    }
  });

  document.addEventListener("change", handleSettingsChange);
  document.addEventListener("input", handleSettingsInput);

  window.addEventListener("pagehide", flushSettings);
};
