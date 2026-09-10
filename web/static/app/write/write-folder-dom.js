import {
  destinationInput,
  folderTree,
} from "./write-modal-elements.js";

export const folderPath = (row) => row?.dataset.folderPath || "/";

export const rowDepth = (row) => {
  const value = row?.style.getPropertyValue("--folder-depth") || "0em";

  return Number.parseFloat(value) || 0;
};

export const folderRows = () =>
  folderTree()?.querySelectorAll("[data-folder-row]") || [];

export const selectFolder = (row) => {
  const input = destinationInput();
  if (!(input instanceof HTMLInputElement)) {
    return;
  }

  input.value = folderPath(row);
  folderRows().forEach((option) => {
    delete option.dataset.selected;
    option.querySelector("[data-folder-select]")?.removeAttribute("data-selected");
  });
  row.dataset.selected = "true";
  row.querySelector("[data-folder-select]")?.setAttribute("data-selected", "true");
};

export const removeDescendants = (row) => {
  const depth = rowDepth(row);
  let next = row.nextElementSibling;
  while (next instanceof HTMLElement && rowDepth(next) > depth) {
    const removeTarget = next;
    next = next.nextElementSibling;
    removeTarget.remove();
  }
};

export const insertChildren = (row, html) => {
  row.insertAdjacentHTML("afterend", html);
};

export const childRowsHTML = (row) => {
  const depth = rowDepth(row);
  const html = [];
  let next = row.nextElementSibling;
  while (next instanceof HTMLElement && rowDepth(next) > depth) {
    if (rowDepth(next) === depth + 1) {
      const clone = next.cloneNode(true);
      clone.dataset.folderExpanded = "false";
      html.push(clone.outerHTML);
    }
    next = next.nextElementSibling;
  }

  return html.join("");
};
