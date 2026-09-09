const cache = new Map();

export const cachedFolderChildren = (path) => cache.get(path);

export const cacheFolderChildren = (path, html) => {
  cache.set(path, html);
};

export const clearFolderCache = () => {
  cache.clear();
};
