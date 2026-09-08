import { initKeyboardShortcuts } from "./keyboard-shortcuts.js";
import { initTableKeyboardNavigation } from "./keyboard-table.js";

export const initKeyboardNavigation = (showToast, setKeyStatus) => {
  const context = {showToast, setKeyStatus};

  initKeyboardShortcuts(context);
  initTableKeyboardNavigation(context);
};
