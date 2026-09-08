export const isTypingTarget = (target) =>
  target instanceof HTMLInputElement ||
  target instanceof HTMLTextAreaElement ||
  target instanceof HTMLSelectElement ||
  target?.isContentEditable;

export const focusSearch = () => {
  const search = document.querySelector("#global-search");
  if (!(search instanceof HTMLInputElement)) {
    return false;
  }

  search.focus();
  search.select();

  return true;
};
