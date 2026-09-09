import { createCommandState } from "./keyboard-command-state.js";
import { focusSearch } from "./keyboard-dom.js";
import {
  handleBackShortcut,
  handleEscape,
  shouldSkipShortcut,
} from "./keyboard-events.js";
import { runBufferedCommand, runSingleKeyCommand } from "./keyboard-runner.js";

const isCountKey = (key, state) =>
  /^[1-9]$/.test(key) || (state.hasCount() && key === "0");

const isBufferedPrefix = (key) =>
  key === "c" || key === "d" || key === "g" || key === "m" || key === "y";

const handleCommandKey = (event, context, state) => {
  if (isCountKey(event.key, state)) {
    state.addCount(event.key);

    return true;
  }

  if (state.hasCommand() && runBufferedCommand(event.key, context, state)) {
    state.reset(`${state.value()}${event.key}`);

    return true;
  }

  if (state.hasCommand()) {
    state.reset();

    return true;
  }

  if (isBufferedPrefix(event.key)) {
    state.addCommand(event.key);

    return true;
  }

  if (runSingleKeyCommand(event.key, context, state.count())) {
    state.reset(`${state.countPrefix()}${event.key}`);

    return true;
  }

  state.reset();

  return false;
};

export const initKeyboardShortcuts = (context) => {
  const state = createCommandState(context.setKeyStatus);

  document.addEventListener("keydown", (event) => {
    if (handleEscape(event)) {
      return;
    }

    if (shouldSkipShortcut(event)) {
      return;
    }

    if (event.key === "/" && focusSearch()) {
      event.preventDefault();
      return;
    }

    if (handleBackShortcut(event, state)) {
      return;
    }

    if (event.altKey) {
      return;
    }

    if (handleCommandKey(event, context, state)) {
      event.preventDefault();
    }
  });
};
