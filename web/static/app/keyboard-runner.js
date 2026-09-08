import { openHelpModal } from "./help.js";
import { copySelectedPath } from "./keyboard-copy-path.js";
import {
  focusFirstEntry,
  focusLastEntry,
  moveEntryFocus,
  openSelectedEntry,
} from "./keyboard-entries.js";
import { goBack, goToParent, refreshListing } from "./keyboard-navigation.js";
import { focusFilesRegion, focusPreviewRegion } from "./keyboard-panes.js";

export const runSingleKeyCommand = (key, context, count) => {
  switch (key) {
  case "j":
  case "ArrowDown":
    return moveEntryFocus(count);
  case "k":
  case "ArrowUp":
    return moveEntryFocus(-count);
  case "h":
    return goToParent(count);
  case "l":
  case "o":
  case "Enter":
    return openSelectedEntry();
  case "G":
  case "End":
    return focusLastEntry();
  case "Home":
    return focusFirstEntry();
  case "r":
  case "R":
    return refreshListing();
  case "Y":
    return copySelectedPath(true, context.showToast);
  case "p":
    return focusPreviewRegion();
  case "f":
    return focusFilesRegion();
  case "b":
    return goBack();
  case "?":
    return openHelpModal(document.activeElement);
  default:
    return false;
  }
};

export const runBufferedCommand = (key, context, state) => {
  const command = `${state.value()}${key}`;
  if (command === "gg" || command.endsWith("gg")) {
    return focusFirstEntry();
  }

  if (command === "yy" || command.endsWith("yy")) {
    return copySelectedPath(false, context.showToast);
  }

  return false;
};
