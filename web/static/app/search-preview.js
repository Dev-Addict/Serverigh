const resetSearchPreviewLinks = (activeLink) => {
  document
    .querySelectorAll("[data-search-preview-link][data-previewed='true']")
    .forEach((link) => {
      if (link !== activeLink) {
        delete link.dataset.previewed;
        delete link.closest("tr")?.dataset.previewed;
      }
    });
};

export const initSearchPreview = () => {
  document.addEventListener(
    "click",
    (event) => {
      if (
        event.defaultPrevented ||
        event.button !== 0 ||
        event.metaKey ||
        event.ctrlKey ||
        event.shiftKey ||
        event.altKey ||
        !(event.target instanceof Element)
      ) {
        return;
      }

      const link = event.target.closest("[data-search-preview-link]");
      if (!(link instanceof HTMLAnchorElement)) {
        return;
      }

      if (link.dataset.previewed === "true") {
        return;
      }

      const previewURL = link.dataset.previewUrl;
      if (!previewURL || !window.htmx) {
        return;
      }

      event.preventDefault();
      event.stopPropagation();
      resetSearchPreviewLinks(link);
      link.dataset.previewed = "true";
      const row = link.closest("tr");
      if (row instanceof HTMLElement) {
        row.dataset.previewed = "true";
      }

      window.htmx.ajax("GET", previewURL, {
        target: "#preview-region",
        swap: "innerHTML",
      });
    },
    true,
  );
};
