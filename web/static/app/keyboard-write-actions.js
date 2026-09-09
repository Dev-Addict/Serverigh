import { selectedEntryRow } from "./keyboard-entries.js";
import { createItem } from "./write-create.js";
import { runEntryWriteAction } from "./write-entry-actions.js";

export const createFolderFromKeyboard = (showToast) =>
  createItem("folder", showToast);

export const copySelectedFromKeyboard = (showToast) =>
  runEntryWriteAction("copy", selectedEntryRow(), showToast);

export const deleteSelectedFromKeyboard = (showToast) =>
  runEntryWriteAction("delete", selectedEntryRow(), showToast);

export const moveSelectedFromKeyboard = (showToast) =>
  runEntryWriteAction("move", selectedEntryRow(), showToast);

export const renameSelectedFromKeyboard = (showToast) =>
  runEntryWriteAction("rename", selectedEntryRow(), showToast);
