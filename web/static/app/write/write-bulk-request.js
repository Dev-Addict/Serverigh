import { syncBulkSelection } from "./write-bulk-selection.js";
import { bulkValues } from "./write-bulk-values.js";

const swapBulkRefresh = (html) => {
  window.htmx.swap("#files-region", html, {swapStyle: "innerHTML"});
  syncBulkSelection();
};

export const postBulkAction = (
  action,
  paths,
  showToast,
  extraValues = {},
) => {
  if (!window.htmx) {
    showToast("Write action unavailable");

    return;
  }

  const values = bulkValues(paths);
  Object.entries(extraValues).forEach(([key, value]) => {
    values.set(key, value);
  });

  fetch(`/actions/bulk/${action}`, {
    method: "POST",
    body: values,
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
      "HX-Request": "true",
    },
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("bulk action failed");
      }

      return response.text();
    })
    .then(swapBulkRefresh)
    .catch(() => showToast("Bulk action failed"));
};
