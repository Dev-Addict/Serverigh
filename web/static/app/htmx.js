export const initHtmxStatus = (queuePreviewOffsetUpdate) => {
  document.addEventListener("htmx:beforeRequest", () => {
    document.documentElement.dataset.loading = "true";
  });

  document.addEventListener("htmx:afterRequest", () => {
    document.documentElement.dataset.loading = "false";
    queuePreviewOffsetUpdate();
  });

  document.addEventListener("htmx:responseError", () => {
    document.documentElement.dataset.loading = "false";
    queuePreviewOffsetUpdate();
  });

  document.addEventListener("htmx:beforeSwap", (event) => {
    const status = event.detail.xhr.status;
    if (status >= 400 && event.detail.xhr.responseText) {
      event.detail.shouldSwap = true;
      event.detail.isError = false;
    }
  });
};
