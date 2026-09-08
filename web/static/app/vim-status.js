let statusTimeout;

export const initVimStatus = () => {
  return (value) => {
    const status = document.querySelector("[data-vim-key-status]");
    if (!(status instanceof HTMLElement)) {
      return;
    }

    window.clearTimeout(statusTimeout);
    status.textContent = value || "--";

    if (!value) {
      return;
    }

    statusTimeout = window.setTimeout(() => {
      status.textContent = "--";
    }, 1600);
  };
};
