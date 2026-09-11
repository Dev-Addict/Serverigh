import { baseActionValues } from "./write-request.js";
import { selectedEntryRows } from "./write-bulk-selection.js";

export const selectedPaths = () => {
  return selectedEntryRows()
    .map((row) => row.dataset.entryPath)
    .filter(Boolean);
};

export const bulkValues = (paths) => {
  const values = new URLSearchParams();
  const base = baseActionValues();
  values.set("path", base.path);
  values.set("sort", base.sort);
  values.set("dir", base.dir);
  paths.forEach((target) => values.append("target", target));

  return values;
};
