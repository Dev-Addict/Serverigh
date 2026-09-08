const heightOf = (selector) => document.querySelector(selector)?.offsetHeight || 0;

const setSizeProperty = (name, value) => {
  document.documentElement.style.setProperty(name, `${value}px`);
};

export const updatePreviewOffset = () => {
  const topBarHeight = heightOf(".top-bar");
  const breadcrumbsHeight = heightOf(".breadcrumbs");
  const footerHeight = heightOf(".status-row");
  const previewChromeHeight = heightOf(".preview-chrome");
  const previewMetadataHeight = heightOf(".preview-metadata");

  setSizeProperty("--serverigh-top-bar-height", topBarHeight);
  setSizeProperty("--serverigh-footer-height", footerHeight);
  setSizeProperty("--serverigh-preview-chrome-height", previewChromeHeight);
  setSizeProperty("--serverigh-preview-metadata-height", previewMetadataHeight);
  setSizeProperty(
    "--serverigh-preview-offset",
    topBarHeight + breadcrumbsHeight,
  );
};
