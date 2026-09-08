export const initClipboardActions = (showToast) => {
  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    const copyButton = event.target.closest("[data-copy-text]");
    if (!copyButton) {
      return;
    }

    event.preventDefault();
    if (!(copyButton instanceof HTMLElement) || !navigator.clipboard) {
      return;
    }

    navigator.clipboard
      .writeText(copyButton.dataset.copyText || "")
      .then(() => {
        copyButton.dataset.copied = "true";
        showToast("Path copied");
        window.setTimeout(() => {
          delete copyButton.dataset.copied;
        }, 1200);
      })
      .catch(() => {
        showToast("Unable to copy path");
      });
  });
};
