const supportedThemes = new Set(["light", "dark"]);

export const applyTheme = (theme) => {
  const nextTheme = supportedThemes.has(theme) ? theme : "light";
  document.documentElement.dataset.theme = nextTheme;
  const select = document.querySelector("[data-theme-select]");
  if (select instanceof HTMLSelectElement) {
    select.value = nextTheme;
  }
};

export const initTheme = () => {
  applyTheme(document.documentElement.dataset.theme || "light");

  document.addEventListener("change", (event) => {
    if (!(event.target instanceof HTMLSelectElement)) {
      return;
    }

    if (!event.target.matches("[data-theme-select]")) {
      return;
    }

    applyTheme(event.target.value);
  });
};
