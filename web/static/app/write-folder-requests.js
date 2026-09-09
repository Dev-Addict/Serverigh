let controller = new AbortController();
let session = 0;
const inFlight = new Map();

export const resetFolderRequests = () => {
  controller.abort();
  controller = new AbortController();
  session += 1;
  inFlight.clear();
};

const requestKey = (url) => `${session}:${url}`;

export const isCurrentFolderSession = (value) => value === session;

export const currentFolderSession = () => session;

export const fetchFolderHTML = (url) => {
  const key = requestKey(url);
  const cached = inFlight.get(key);
  if (cached) {
    return cached;
  }

  const requestSession = session;
  const request = fetch(url, {signal: controller.signal})
    .then((response) => {
      if (!response.ok) {
        throw new Error("folder tree request failed");
      }

      return response.text();
    })
    .then((html) => ({
      html,
      session: requestSession,
    }))
    .finally(() => {
      inFlight.delete(key);
    });

  inFlight.set(key, request);

  return request;
};
