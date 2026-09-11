const refreshIntervalMS = 4000;

const state = {
  url: "",
  version: "",
  inFlight: false,
};

const filesPanel = () => document.querySelector(".file-panel[data-path]");

const refreshLinkForPanel = (panel) =>
  panel?.querySelector(".quiet-link[hx-get]");

const versionURLForLink = (link) => {
  const refreshURL = link?.getAttribute("hx-get");
  if (!refreshURL) {
    return "";
  }

  const url = new URL(refreshURL, window.location.origin);
  url.pathname = "/partials/files/version";

  return url.pathname + url.search;
};

const resetRefreshVersion = () => {
  state.url = "";
  state.version = "";
};

const settledFilesRegion = (event) =>
  event.target?.id === "files-region" ||
  event.detail?.target?.id === "files-region";

const pollVisibleFiles = async () => {
  if (
    document.hidden ||
    state.inFlight ||
    document.documentElement.dataset.loading === "true"
  ) {
    return;
  }

  const panel = filesPanel();
  const refreshLink = refreshLinkForPanel(panel);
  const url = versionURLForLink(refreshLink);
  if (!panel || !refreshLink || !url) {
    resetRefreshVersion();

    return;
  }

  state.inFlight = true;
  try {
    const response = await fetch(url, {
      cache: "no-store",
      headers: {Accept: "application/json"},
    });
    if (!response.ok) {
      return;
    }

    const body = await response.json();
    if (state.url !== url || !state.version) {
      state.url = url;
      state.version = body.version;

      return;
    }

    if (body.version && body.version !== state.version) {
      state.version = body.version;
      refreshLink.click();
    }
  } catch {
    // Keep polling; transient filesystem errors should not interrupt browsing.
  } finally {
    state.inFlight = false;
  }
};

export const initFileRefresh = () => {
  window.setInterval(pollVisibleFiles, refreshIntervalMS);
  document.addEventListener("visibilitychange", () => {
    if (!document.hidden) {
      resetRefreshVersion();
      pollVisibleFiles();
    }
  });
  document.addEventListener("htmx:afterSettle", (event) => {
    if (settledFilesRegion(event)) {
      resetRefreshVersion();
    }
  });
  pollVisibleFiles();
};
