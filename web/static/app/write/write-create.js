import { baseActionValues, postWriteAction } from "./write-request.js";
import { openNameModal } from "./write-modal.js";

export const createItem = (kind, showToast) => {
  if (!document.querySelector("[data-write-menu]")) {
    showToast("Write mode is disabled");

    return true;
  }

  const action = kind === "file" ? "file" : "mkdir";
  openNameModal({
    nameLabel: kind === "file" ? "File name" : "Folder path",
    submit: ({name}) => {
      if (!postWriteAction(action, {...baseActionValues(), name})) {
        showToast("Write action unavailable");

        return;
      }

      showToast("Write action submitted");
    },
    submitLabel: "Create",
    title: kind === "file" ? "New file" : "New folder",
  });

  return true;
};

export const initWriteCreateActions = (showToast) => {
  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    const button = event.target.closest("[data-write-create]");
    if (!(button instanceof HTMLButtonElement)) {
      return;
    }

    event.preventDefault();
    createItem(button.dataset.writeCreate, showToast);
  });
};
