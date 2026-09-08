const themeStorageKey = "serverigh.theme";
const supportedThemes = new Set(["light", "dark"]);

const storedValue = (key) => {
  try {
    return window.localStorage?.getItem(key);
  } catch {
    return null;
  }
};

const storeValue = (key, value) => {
  try {
    window.localStorage?.setItem(key, value);
  } catch {
    return;
  }
};

export const applyTheme = (theme) => {
  const nextTheme = supportedThemes.has(theme) ? theme : "light";
  document.documentElement.dataset.theme = nextTheme;
  const select = document.querySelector("[data-theme-select]");
  if (select instanceof HTMLSelectElement) {
    select.value = nextTheme;
  }
};

export const initTheme = () => {
  applyTheme(storedValue(themeStorageKey) || "light");

  document.addEventListener("change", (event) => {
    if (!(event.target instanceof HTMLSelectElement)) {
      return;
    }

    if (!event.target.matches("[data-theme-select]")) {
      return;
    }

    applyTheme(event.target.value);
    storeValue(themeStorageKey, event.target.value);
  });
};
