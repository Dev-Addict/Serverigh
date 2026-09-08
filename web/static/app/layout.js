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

const observePreviewChrome = (queueUpdate) => {
  if (!("ResizeObserver" in window)) {
    return;
  }

  const observer = new ResizeObserver(queueUpdate);
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
};

export const initPreviewLayout = () => {
  const queueUpdate = () => {
    window.requestAnimationFrame(updatePreviewOffset);
  };

  queueUpdate();
  window.addEventListener("resize", queueUpdate);
  window.addEventListener("load", queueUpdate);
  document.addEventListener("htmx:afterSettle", queueUpdate);
  document.addEventListener("htmx:historyRestore", queueUpdate);
  observePreviewChrome(queueUpdate);

  return queueUpdate;
};
