export const initHistoryBack = () => {
  document.addEventListener("click", (event) => {
    if (!(event.target instanceof Element)) {
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
};
