export const downloadBulk = (paths, showToast) => {
  const form = document.createElement("form");
  form.method = "POST";
  form.action = "/download/bulk";
  form.hidden = true;
  paths.forEach((target) => {
    const input = document.createElement("input");
    input.name = "target";
    input.value = target;
    form.append(input);
  });
  document.body.append(form);
  form.submit();
  form.remove();
  showToast("Download started");
};
