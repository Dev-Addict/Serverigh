import {
  cacheFolderChildren,
  cachedFolderChildren,
  clearFolderCache,
} from "./write-folder-cache.js";
import {
  currentFolderSession,
  fetchFolderHTML,
  isCurrentFolderSession,
  resetFolderRequests,
} from "./write-folder-requests.js";
import {
  childRowsHTML,
  folderPath,
  folderRows,
  insertChildren,
  removeDescendants,
  selectFolder,
} from "./write-folder-dom.js";
import { folderTree } from "./write-modal-elements.js";

export { selectFolder } from "./write-folder-dom.js";

const pending = new Set();

export const clearFolderTreeCache = () => {
  resetFolderRequests();
  clearFolderCache();
  pending.clear();
  const target = folderTree();
  if (target instanceof HTMLElement) {
    target.textContent = "";
  }
};

const childrenURL = (path) =>
  `/partials/folders?path=${encodeURIComponent(path || "/")}`;

const branchURL = (selectedPath) =>
  `/partials/folders?selected=${encodeURIComponent(selectedPath || "/")}`;

const seedExpandedRows = () => {
  folderRows().forEach((row) => {
    if (row.dataset.folderExpanded === "true") {
      cacheFolderChildren(folderPath(row), childRowsHTML(row));
    }
  });
};

export const toggleFolder = (row) => {
  const path = folderPath(row);
  if (row.dataset.folderExpanded === "true") {
    row.dataset.folderExpanded = "false";
    removeDescendants(row);

    return;
  }

  const cached = cachedFolderChildren(path);
  row.dataset.folderExpanded = "true";
  if (cached !== undefined) {
    insertChildren(row, cached);

    return;
  }
  if (pending.has(path)) {
    return;
  }

  pending.add(path);
  fetchFolderHTML(childrenURL(path))
    .then(({html, session}) => {
      if (!isCurrentFolderSession(session)) {
        return;
      }
      cacheFolderChildren(path, html);
      if (row.dataset.folderExpanded === "true") {
        removeDescendants(row);
        insertChildren(row, html);
      }
    })
    .catch(() => {
      if (row.isConnected) {
        row.dataset.folderExpanded = "false";
      }
    })
    .finally(() => {
      pending.delete(path);
    });
};

export const loadFolderTree = (selectedPath) => {
  const target = folderTree();
  if (!(target instanceof HTMLElement)) {
    return;
  }

  resetFolderRequests();
  clearFolderCache();
  pending.clear();
  target.textContent = "Loading folders...";
  const session = currentFolderSession();
  fetchFolderHTML(branchURL(selectedPath))
    .then(({html}) => {
      if (!isCurrentFolderSession(session)) {
        return;
      }
      target.innerHTML = html;
      seedExpandedRows();
      const selected = target.querySelector("[data-selected='true']");
      const row = selected?.closest("[data-folder-row]");
      if (row instanceof HTMLElement) {
        selectFolder(row);
      }
    })
    .catch(() => {
      if (isCurrentFolderSession(session)) {
        target.textContent = "Unable to load folders.";
      }
    });
};
