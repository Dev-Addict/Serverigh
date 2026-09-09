const openMenus = () => document.querySelectorAll("[data-write-menu][open]");

const closeUploadMenus = () => {
  openMenus().forEach((menu) => {
    menu.removeAttribute("open");
  });
};

export const initWriteMenu = () => {
  document.addEventListener("click", (event) => {
    const clickedMenu = event.target instanceof Element
      && event.target.closest("[data-write-menu]");
    if (clickedMenu) {
      return;
    }

    closeUploadMenus();
  });

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeUploadMenus();
    }
  });
};
