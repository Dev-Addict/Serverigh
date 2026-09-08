export const focusFilesRegion = () => {
  const filesRegion = document.querySelector("#files-region");
  if (!(filesRegion instanceof HTMLElement)) {
    return false;
  }

  filesRegion.tabIndex = -1;
  filesRegion.focus();

  return true;
};

export const focusPreviewRegion = () => {
  const previewRegion = document.querySelector("#preview-region");
  if (!(previewRegion instanceof HTMLElement)) {
    return false;
  }

  previewRegion.tabIndex = -1;
  previewRegion.focus();

  return true;
};
