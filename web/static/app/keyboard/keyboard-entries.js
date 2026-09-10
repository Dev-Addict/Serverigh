const entryLinks = () => [
  ...document.querySelectorAll(".file-table tbody .entry-link"),
];

const linkIndex = (links, activeElement) => {
  const activeRow = activeElement?.closest?.("[data-entry-row]");
  return links.findIndex((link) => {
    if (link === activeElement || link.contains(activeElement)) {
      return true;
    }

    return activeRow && link.closest("[data-entry-row]") === activeRow;
  });
};

const focusEntry = (links, index) => {
  const link = links[index];
  if (!(link instanceof HTMLElement)) {
    return;
  }
  link.focus();
  link.scrollIntoView({block: "nearest", inline: "nearest"});
};

export const selectedEntryLink = () => {
	const links = entryLinks();
	const index = linkIndex(links, document.activeElement);

	return index < 0 ? null : links[index];
};

export const selectedEntryRow = () =>
  selectedEntryLink()?.closest("[data-entry-row]");

export const moveEntryFocus = (delta) => {
  const links = entryLinks();
  if (links.length === 0) {
    return false;
  }
  const activeIndex = linkIndex(links, document.activeElement);
  if (activeIndex < 0) {
    focusEntry(links, delta < 0 ? links.length - 1 : 0);

    return true;
  }
  const targetIndex = Math.max(
    0,
    Math.min(activeIndex + delta, links.length - 1),
  );
  focusEntry(links, targetIndex);

  return true;
};

export const focusFirstEntry = () => {
  const links = entryLinks();
  if (links.length === 0) {
    return false;
  }
  focusEntry(links, 0);

  return true;
};

export const focusLastEntry = () => {
  const links = entryLinks();
  if (links.length === 0) {
    return false;
  }
  focusEntry(links, links.length - 1);

  return true;
};

export const openSelectedEntry = () => {
  const link = selectedEntryLink();
  if (!(link instanceof HTMLElement)) {
    return false;
  }
  link.click();

  return true;
};
