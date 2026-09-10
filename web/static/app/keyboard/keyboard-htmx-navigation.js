export const htmxNavigate = (browseURL, filesURL) => {
  const link = document.createElement("a");
  link.href = browseURL;
  link.setAttribute("hx-get", filesURL);
  link.setAttribute("hx-target", "#files-region");
  link.setAttribute("hx-swap", "innerHTML");
  link.setAttribute("hx-push-url", browseURL);
  link.hidden = true;
  document.body.append(link);
  window.htmx.process(link);
  link.addEventListener(
    "htmx:afterRequest",
    () => {
      link.remove();
    },
    {once: true},
  );
  link.click();
};
