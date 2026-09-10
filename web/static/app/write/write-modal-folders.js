import {
  cacheFolderChildren,
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
  selectFolder,
} from "./write-folder-dom.js";
import { clearPendingFolders } from "./write-folder-pending.js";
import { folderTree } from "./write-modal-elements.js";

export { selectFolder } from "./write-folder-dom.js";
export { toggleFolder } from "./write-folder-toggle.js";

export const clearFolderTreeCache = () => {
  resetFolderRequests();
  clearFolderCache();
  clearPendingFolders();
  const target = folderTree();
  if (target instanceof HTMLElement) {
    target.textContent = "";
  }
};

const branchURL = (selectedPath) =>
  `/partials/folders?selected=${encodeURIComponent(selectedPath || "/")}`;

const seedExpandedRows = () => {
  folderRows().forEach((row) => {
    if (row.dataset.folderExpanded === "true") {
      cacheFolderChildren(folderPath(row), childRowsHTML(row));
    }
  });
};

export const loadFolderTree = (selectedPath) => {
  const target = folderTree();
  if (!(target instanceof HTMLElement)) {
    return;
  }

  resetFolderRequests();
  clearFolderCache();
  clearPendingFolders();
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
