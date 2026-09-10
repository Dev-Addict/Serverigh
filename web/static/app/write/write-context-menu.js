import {
  closeContextMenu,
  openContextMenu,
} from "./write-context-state.js";

const buttonPoint = (button) => {
  const rect = button.getBoundingClientRect();

  return {
    x: rect.right,
    y: rect.bottom,
  };
};

export const initContextMenu = () => {
  document.addEventListener("contextmenu", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    const row = event.target.closest("[data-entry-row]");
    const panel = event.target.closest(".file-panel");
    if (!panel) {
      return;
    }

    event.preventDefault();
    openContextMenu(row, event.clientX, event.clientY);
  });

  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    const button = event.target.closest("[data-entry-menu-button]");
    if (!(button instanceof HTMLButtonElement)) {
      return;
    }

    const row = button.closest("[data-entry-row]");
    const point = buttonPoint(button);
    event.preventDefault();
    event.stopPropagation();
    openContextMenu(row, point.x, point.y);
  });

  document.addEventListener("click", (event) => {
    const clickedMenu = event.target instanceof Element
      && event.target.closest("[data-write-context-menu]");
    const clickedButton = event.target instanceof Element
      && event.target.closest("[data-entry-menu-button]");
    if (clickedMenu || clickedButton) {
      return;
    }

    closeContextMenu();
  });
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeContextMenu();
    }
  });
};
