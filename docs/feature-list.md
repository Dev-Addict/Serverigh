# Serverigh Feature List

Serverigh should be a local-first web file explorer built for fast browsing,
readable previews, and careful filesystem boundaries. The Go and HTMX stack
should keep the product small: server-rendered HTML, focused partial updates,
minimal client state, and no heavy frontend build step.

## Product Principles

- Make the first screen the actual file explorer.
- Keep local filesystem access explicit, scoped, and visible.
- Default to read-only behavior until write operations are deliberately enabled.
- Treat Markdown and CSV previews as first-class workflows.
- Keep URLs meaningful so reload, back, forward, and sharing a local view all
  preserve context.
- Prefer standard browser behavior and Go server rendering over custom client
  framework state.

## MVP Features

- Root-scoped browsing.
  Priority: must have. Browse from a configured filesystem root. Start with one
  root, and prevent path traversal and symlink escapes.
- Split-pane explorer.
  Priority: must have. Use a dense, desktop-first layout with folders and files
  on the left and preview/details on the right.
- Directory listing.
  Priority: must have. Show folders and files with name, type, size, modified
  time, and permissions. Sort by name, type, size, and modified date.
- Breadcrumb navigation.
  Priority: must have. Provide clickable path segments from root to the current
  folder, while respecting the configured root.
- URL state.
  Priority: must have. Encode the active folder or file in the URL with query
  params such as `?view=folder&path=...` and `?view=file&path=...`.
- Markdown preview.
  Priority: must have. Render Markdown with headings, lists, links, code blocks,
  blockquotes, and GitHub-style pipe tables.
- CSV preview.
  Priority: must have. Render CSV as a readable table with horizontal scrolling
  and sane column sizing for wide files and text-heavy cells.
- Text/code preview.
  Priority: must have. Preview text, config, log, JSON, YAML, Go, JS, TS, CSS,
  HTML, and shell files. Truncate large files with bounded reads.
- Media preview.
  Priority: should have. Preview images, PDFs, audio, and video with native
  browser rendering where possible.
- Binary file details.
  Priority: must have. Show metadata and download/open actions for unknown
  binary files instead of trying to render arbitrary binary content.
- Copy path.
  Priority: should have. Copy relative and absolute paths from a compact icon
  button with a tooltip.
- Fuzzy filename search.
  Priority: must have. Search as the user types, rank current-folder matches
  first, and open file results in their parent directory with preview selected.
- Config-driven hidden-file visibility.
  Priority: should have. Show or hide dotfiles based only on server
  configuration.
- Read-only mode.
  Priority: must have. Disable rename, delete, move, upload, and edit by
  default.
- Health endpoint.
  Priority: must have. Provide a small endpoint for smoke checks and future
  packaging.

## Version 1 Features

- File metadata drawer.
  Priority: high. Show owner, permissions, mime type, line count, dimensions,
  and checksums where cheap. Compute expensive fields lazily.
- Multi-root bookmarks.
  Priority: high. Configure named roots such as Projects, Downloads, and Home.
  Each root must have its own path boundary.
- File operations.
  Priority: high. Add new folder, rename, move, copy, delete, and upload only
  when write mode is enabled. Confirm destructive actions.
- Trash-first delete.
  Priority: high. Prefer moving to trash when supported. Fall back to permanent
  delete only with explicit confirmation.
- Bulk selection.
  Priority: medium. Select multiple files for download, move, copy, or delete
  with predictable keyboard and checkbox behavior.
- Archive download.
  Priority: medium. Download selected files or folders as a streamed zip
  archive.
- File watcher refresh.
  Priority: medium. Refresh visible directory contents when files change on
  disk. Use polling first if native watchers add too much complexity.
- Better code preview.
  Priority: medium. Add syntax highlighting and line numbers with server-side
  highlighting or a small static highlighter.
- Preview size controls.
  Priority: medium. Switch between compact, comfortable, and full-width preview
  density, especially for CSV and Markdown.
- Keyboard navigation.
  Priority: medium. Support arrow keys, enter to open, backspace to parent, and
  slash to search. Keep shortcuts discoverable through tooltips and menus.

## Later Features

- Git-aware folders.
  Priority: medium. Show branch, modified files, untracked files, and quick
  diffs only inside Git repositories.
- Text diff viewer.
  Priority: medium. Compare two selected text files for config and generated
  document review.
- Favorites and recent files.
  Priority: medium. Pin folders/files and show recent locations from a small
  local config file.
- Tagging and notes.
  Priority: low. Store user-defined notes for files and folders outside the
  original files.
- Archive browsing.
  Priority: low. Browse zip/tar contents as virtual read-only folders without
  extracting.
- Remote adapters.
  Priority: low. Add optional SFTP, SMB, S3, or WebDAV adapters after the local
  MVP is stable.
- Plugin previews.
  Priority: low. Register custom previewers for specialized file types after
  core preview contracts are stable.

## Features To Avoid Early

- Browser-based terminal execution.
- Multi-user permissions before the local single-user workflow is excellent.
- Full in-browser IDE behavior.
- Heavy SPA routing or a frontend build system unless the product outgrows HTMX.
- AI file indexing before search, preview, and safety basics are finished.
