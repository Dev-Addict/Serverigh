import { closeContextMenu } from "./write-context-state.js";
import { downloadBulk } from "./write-bulk-download.js";
import { postBulkAction } from "./write-bulk-request.js";
import { selectedPaths } from "./write-bulk-values.js";
import { openDestinationModal } from "./write-modal.js";
import { currentPath } from "./write-request.js";

const bulkRefreshActions = {
  delete: "Trash submitted",
  duplicate: "Duplicate submitted",
};

const openBulkDestination = (action, paths, showToast) => {
  const label = action === "copy" ? "Copy" : "Move";
  openDestinationModal({
    selectedPath: currentPath(),
    submit: (values) => {
      postBulkAction(action, paths, showToast, {
        destination: values.destination || "/",
      });
      showToast(`${label} submitted`);
    },
    submitLabel: label,
    summary: `${label} ${paths.length} selected entries to:`,
    title: `${label} selected entries`,
  });
};

export const runBulkAction = (action, showToast) => {
  const paths = selectedPaths();
  if (paths.length === 0) {
    showToast("Select entries first");

    return true;
  }

  if (action === "download") {
    downloadBulk(paths, showToast);

    return true;
  }
  if (action === "copy" || action === "move") {
    openBulkDestination(action, paths, showToast);

    return true;
  }

  postBulkAction(action, paths, showToast);
  showToast(bulkRefreshActions[action] || "Bulk action submitted");

  return true;
};

export const initBulkActions = (showToast) => {
  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    const button = event.target.closest("[data-bulk-action]");
    if (!(button instanceof HTMLButtonElement)) {
      return;
    }

    event.preventDefault();
    runBulkAction(button.dataset.bulkAction, showToast);
    closeContextMenu();
  });
};
