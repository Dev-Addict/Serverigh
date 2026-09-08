import { htmxNavigate } from "./keyboard-htmx-navigation.js";

const pathURL = (prefix, path) => {
  const values = new URLSearchParams(window.location.search);
  values.set("path", path);
  values.delete("file");

  return `${prefix}?${values.toString()}`;
};

const parentPath = (activePath, count) => {
  const segments = activePath.split("/").filter(Boolean);
  segments.splice(Math.max(0, segments.length - count), count);

  return segments.length === 0 ? "/" : `/${segments.join("/")}`;
};

const activePath = () => {
  const panel = document.querySelector(".file-panel[data-path]");
  if (panel instanceof HTMLElement && panel.dataset.path) {
    return panel.dataset.path;
  }

  return new URLSearchParams(window.location.search).get("path") || "/";
};

export const goToParent = (count) => {
  const targetPath = parentPath(activePath(), count);
  const browseURL = pathURL("/browse", targetPath);
  if (window.htmx) {
    htmxNavigate(browseURL, pathURL("/partials/files", targetPath));

    return true;
  }

  window.location.assign(browseURL);

  return true;
};

export const goBack = () => {
  const backLink = document.querySelector("[data-history-back]");
  if (!(backLink instanceof HTMLElement)) {
    return false;
  }

  backLink.click();

  return true;
};

export const refreshListing = () => {
  const search = document.querySelector("#global-search");
  if (search instanceof HTMLInputElement && search.value.trim()) {
    search.dispatchEvent(new Event("search", {bubbles: true}));

    return true;
  }

  const refresh = document.querySelector(".file-panel .quiet-link");
  if (!(refresh instanceof HTMLElement)) {
    return false;
  }

  refresh.click();

  return true;
};
