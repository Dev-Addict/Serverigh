import {
  cacheFolderChildren,
  cachedFolderChildren,
} from "./write-folder-cache.js";
import {
  fetchFolderHTML,
  isCurrentFolderSession,
} from "./write-folder-requests.js";
import {
  folderPath,
  insertChildren,
  removeDescendants,
} from "./write-folder-dom.js";
import {
  addPendingFolder,
  deletePendingFolder,
  hasPendingFolder,
} from "./write-folder-pending.js";

const childrenURL = (path) => {
  return `/partials/folders?path=${encodeURIComponent(path || "/")}`;
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
  if (hasPendingFolder(path)) {
    return;
  }

  addPendingFolder(path);
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
      deletePendingFolder(path);
    });
};
