export const settingsModal = () => {
  return document.querySelector("[data-settings-modal]");
};

const columnKey = (name) => {
  return `column${name[0].toUpperCase()}${name.slice(1)}`;
};

export const syncColumnToggle = (toggle) => {
  const column = toggle.dataset.columnToggle;
  if (!column) {
    return;
  }

  document.body.dataset[columnKey(column)] = String(toggle.checked);
};

export const syncColumnToggles = () => {
  document.querySelectorAll("[data-column-toggle]").forEach((toggle) => {
    if (toggle instanceof HTMLInputElement) {
      syncColumnToggle(toggle);
    }
  });
};

export const syncMaxPreviewBytes = (input) => {
  document.body.dataset.maxPreviewBytes = input.value;
};
