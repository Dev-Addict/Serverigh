const settingChecked = (name) => {
  const toggle = document.querySelector(`[data-column-toggle="${name}"]`);
  return toggle instanceof HTMLInputElement && toggle.checked;
};

const settingsPayload = () => {
  const values = new URLSearchParams();
  const theme = document.querySelector("[data-theme-select]");
  const maxPreviewBytes = document.querySelector(
    "input[data-max-preview-bytes]",
  );

  if (theme instanceof HTMLSelectElement) {
    values.set("theme", theme.value);
  }
  if (maxPreviewBytes instanceof HTMLInputElement) {
    values.set("max_preview_bytes", maxPreviewBytes.value);
  }

  values.set("column_size", String(settingChecked("size")));
  values.set("column_modified", String(settingChecked("modified")));
  values.set("column_created", String(settingChecked("created")));
  values.set("column_mode", String(settingChecked("mode")));

  return values;
};

export const saveSettings = () => {
  fetch("/settings", {
    method: "POST",
    body: settingsPayload(),
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    keepalive: true,
  }).catch(() => {});
};

export const flushSettings = () => {
  const payload = settingsPayload();
  if (navigator.sendBeacon) {
    const body = new Blob([payload.toString()], {
      type: "application/x-www-form-urlencoded",
    });
    navigator.sendBeacon("/settings", body);

    return;
  }

  saveSettings();
};
