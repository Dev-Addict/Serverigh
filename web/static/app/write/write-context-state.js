import { bulkSelectionForRow } from "./write-bulk-selection.js";

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

const syncBulkContext = (menu, row) => {
  const selection = bulkSelectionForRow(row);
  menu.dataset.bulk = String(selection.bulk);
  menu.dataset.bulkCount = String(selection.count);
  menu
    .querySelectorAll("[data-bulk-count]")
    .forEach((target) => {
      target.textContent = String(selection.count);
    });
};

export const openContextMenu = (row, x, y) => {
  const menu = contextMenu();
  if (!(menu instanceof HTMLElement)) {
    return false;
  }

  activeRow = row instanceof HTMLElement ? row : null;
  if (activeRow?.dataset.trashRoot === "true") {
    activeRow = null;

    return false;
  }
  if (!activeRow && !hasWriteMenu()) {
    return false;
  }

  menu.dataset.hasEntry = activeRow ? "true" : "false";
  syncBulkContext(menu, activeRow);
  menu.dataset.open = "true";
  menu.setAttribute("aria-hidden", "false");
  positionMenu(menu, x, y);

  return true;
};
