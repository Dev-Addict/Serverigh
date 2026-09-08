import { updatePreviewOffset } from "./layout-preview.js";

const observePreviewChrome = (queueUpdate) => {
  if (!("ResizeObserver" in window)) {
    return;
  }

  const observer = new ResizeObserver(queueUpdate);
  const observedNodes = new Set();

  const observeNode = (node) => {
    if (observedNodes.has(node)) {
      return;
    }

    observedNodes.add(node);
    observer.observe(node);
  };

  const observeChrome = () => {
    observer.disconnect();
    observedNodes.clear();

    document
      .querySelectorAll(".top-bar, .breadcrumbs, .status-row")
      .forEach(observeNode);

    document
      .querySelectorAll(".preview-chrome, .preview-metadata")
      .forEach(observeNode);
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
