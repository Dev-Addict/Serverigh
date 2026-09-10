const appendUploadFiles = (data, files) => {
  [...files].forEach((file) => {
    data.append("relative_path", file.webkitRelativePath || file.name);
    data.append("file", file, file.name);
  });
};

const uploadFiles = (input, showToast) => {
  const form = input.form;
  if (!form || !input.files || input.files.length === 0 || !window.htmx) {
    return;
  }

  const data = new FormData(form);
  data.delete(input.name);
  appendUploadFiles(data, input.files);

  fetch(form.action || "/actions/upload", {
    method: "POST",
    body: data,
    headers: {"HX-Request": "true"},
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("upload failed");
      }

      return response.text();
    })
    .then((html) => {
      window.htmx.swap("#files-region", html, {swapStyle: "innerHTML"});
      input.value = "";
      showToast("Upload complete");
    })
    .catch(() => showToast("Upload failed"));
};

export const initWriteUploads = (showToast) => {
  document.addEventListener("change", (event) => {
    if (!(event.target instanceof HTMLInputElement)) {
      return;
    }
    if (!event.target.matches("[data-write-upload]")) {
      return;
    }
    if (!event.target.files || event.target.files.length === 0) {
      return;
    }

    showToast("Upload submitted");
    uploadFiles(event.target, showToast);
  });
};
