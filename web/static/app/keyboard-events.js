import { isTypingTarget } from "./keyboard-dom.js";
import { goBack } from "./keyboard-navigation.js";

export const handleEscape = (event) => {
  if (event.key !== "Escape" || !isTypingTarget(event.target)) {
    return false;
  }

  event.target.blur();

  return true;
};

export const shouldSkipShortcut = (event) =>
  isTypingTarget(event.target) || event.metaKey || event.ctrlKey;

export const handleBackShortcut = (event, state) => {
  const wantsBack = event.key === "Backspace" ||
    (event.altKey && event.key === "ArrowLeft");
  if (!wantsBack || !goBack()) {
    return false;
  }

  state.reset(event.key);
  event.preventDefault();

  return true;
};
