export const debounce = (callback, delay) => {
  let timeout;

  return () => {
    window.clearTimeout(timeout);
    timeout = window.setTimeout(callback, delay);
  };
};
