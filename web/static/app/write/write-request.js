const actionPath = (action) => `/actions/${action}`;

const hiddenValue = (name, fallback) => {
  const input = document.querySelector(`#search-region input[name='${name}']`);

  return input?.value || fallback;
};

export const currentPath = () => hiddenValue("path", "/");

export const baseActionValues = (targetPath = "") => ({
  dir: hiddenValue("dir", "asc"),
  path: currentPath(),
  sort: hiddenValue("sort", "name"),
  target: targetPath,
});

export const postWriteAction = (action, values) => {
  if (!window.htmx) {
    return false;
  }

  window.htmx.ajax("POST", actionPath(action), {
    target: "#files-region",
    swap: "innerHTML",
    values,
  });

  return true;
};
