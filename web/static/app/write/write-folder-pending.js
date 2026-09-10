const pending = new Set();

export const hasPendingFolder = (path) => pending.has(path);

export const addPendingFolder = (path) => {
  pending.add(path);
};

export const deletePendingFolder = (path) => {
  pending.delete(path);
};

export const clearPendingFolders = () => {
  pending.clear();
};
