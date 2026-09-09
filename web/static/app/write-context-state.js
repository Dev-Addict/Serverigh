const contextMenu = () => document.querySelector("[data-write-context-menu]");

let activeRow = null;

export const closeContextMenu = () => {
  const menu = contextMenu();
  if (!(menu instanceof HTMLElement)) {
    return;
  }

  delete menu.dataset.open;
  activeRow = null;
  menu.setAttribute("aria-hidden", "true");
};

export const contextMenuRow = () => activeRow;

const hasWriteMenu = () => Boolean(document.querySelector("[data-write-menu]"));

const positionMenu = (menu, x, y) => {
  menu.style.left = `${x}px`;
  menu.style.top = `${y}px`;
};

export const openContextMenu = (row, x, y) => {
  const menu = contextMenu();
  if (!(menu instanceof HTMLElement)) {
    return false;
  }

  activeRow = row instanceof HTMLElement ? row : null;
  if (!activeRow && !hasWriteMenu()) {
    return false;
  }

  menu.dataset.hasEntry = activeRow ? "true" : "false";
  menu.dataset.open = "true";
  menu.setAttribute("aria-hidden", "false");
  positionMenu(menu, x, y);

  return true;
};
