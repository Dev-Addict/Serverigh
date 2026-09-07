(() => {
  document.documentElement.dataset.serverigh = "ready";

  const updatePreviewOffset = () => {
    const topBar = document.querySelector(".top-bar");
    const breadcrumbs = document.querySelector(".breadcrumbs");
    const footer = document.querySelector(".status-row");
    const previewChrome = document.querySelector(".preview-chrome");
    const topBarHeight = topBar?.offsetHeight || 0;
    const breadcrumbsHeight = breadcrumbs?.offsetHeight || 0;
    const footerHeight = footer?.offsetHeight || 0;
    const previewChromeHeight = previewChrome?.offsetHeight || 0;
    const offset = topBarHeight + breadcrumbsHeight;

    document.documentElement.style.setProperty(
      "--serverigh-top-bar-height",
      `${topBarHeight}px`,
    );

    document.documentElement.style.setProperty(
      "--serverigh-footer-height",
      `${footerHeight}px`,
    );

    document.documentElement.style.setProperty(
      "--serverigh-preview-chrome-height",
      `${previewChromeHeight}px`,
    );

    document.documentElement.style.setProperty(
      "--serverigh-preview-offset",
      `${offset}px`,
    );
  };

  const queuePreviewOffsetUpdate = () => {
    window.requestAnimationFrame(updatePreviewOffset);
  };

  let toastTimeout;
  const showToast = (message) => {
    const toast = document.querySelector("#toast-region");
    if (!(toast instanceof HTMLElement)) {
      return;
    }

    window.clearTimeout(toastTimeout);
    toast.textContent = message;
    toast.dataset.visible = "true";
    toastTimeout = window.setTimeout(() => {
      delete toast.dataset.visible;
    }, 1800);
  };

  queuePreviewOffsetUpdate();
  window.addEventListener("resize", queuePreviewOffsetUpdate);
  window.addEventListener("load", queuePreviewOffsetUpdate);

  if ("ResizeObserver" in window) {
    const observer = new ResizeObserver(queuePreviewOffsetUpdate);
    const observeChrome = () => {
      document
        .querySelectorAll(".top-bar, .breadcrumbs, .status-row")
        .forEach((node) => {
          observer.observe(node);
        });
      document.querySelectorAll(".preview-chrome").forEach((node) => {
        observer.observe(node);
      });
    };

    observeChrome();
    document.addEventListener("htmx:afterSettle", observeChrome);
  }

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

  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }

    const copyButton = event.target.closest("[data-copy-text]");
    if (copyButton) {
      event.preventDefault();
      if (!(copyButton instanceof HTMLElement) || !navigator.clipboard) {
        return;
      }

      navigator.clipboard.writeText(copyButton.dataset.copyText || "").then(() => {
        copyButton.dataset.copied = "true";
        showToast("Path copied");
        window.setTimeout(() => {
          delete copyButton.dataset.copied;
        }, 1200);
      });

      return;
    }

    const backLink = event.target.closest("[data-history-back]");
    if (!backLink) {
      return;
    }

    event.preventDefault();
    if (window.history.length > 1) {
      window.history.back();

      return;
    }

    if (backLink instanceof HTMLAnchorElement) {
      window.location.assign(backLink.href);
    }
  });

  document.addEventListener("htmx:afterSettle", queuePreviewOffsetUpdate);
  document.addEventListener("htmx:historyRestore", queuePreviewOffsetUpdate);
})();
