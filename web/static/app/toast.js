export const createToast = () => {
  let toastTimeout;

  return (message) => {
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
};
