# Serverigh Project Specification

## Summary

Serverigh is a local web file explorer written in Go with HTMX-enhanced
server-rendered HTML. It lets a user browse a configured filesystem root,
inspect folders, and preview common file types from a clean
browser UI.

The first production-quality version should be local-first, read-only by
default, desktop-focused, and strict about filesystem boundaries.

## Goals

- Provide a fast web UI for browsing local files and folders.
- Render Markdown and CSV in a highly readable form.
- Preserve active folder/file state in the URL.
- Keep implementation small enough to understand and maintain.
- Use Go for filesystem access, routing, streaming, and templates.
- Use HTMX for partial page updates without a full SPA framework.
- Make write operations opt-in and clearly separated from the read-only core.

## Non-Goals

- Do not build a cloud storage product in the MVP.
- Do not expose the app to the public internet by default.
- Do not implement multi-user auth before the local single-user workflow is
  stable.
- Do not require Node, React, Vite, or a frontend bundler for the core app.

## Target User

The primary user is a developer or operator who wants a browser-based view over
local project folders, generated documents, CSVs, Markdown files, logs, images,
and miscellaneous workspace assets.

Key jobs:

- Browse project and document folders quickly.
- Preview Markdown and CSV without opening another app.
- Keep a stable URL for the current folder or file.
- Safely inspect system files without accidental mutation.

## Technology Stack

- Language: Go.
- HTTP server: Fiber.
- Templates: Go `html/template`.
- UI updates: HTMX.
- Styling: static CSS with no required build step.
- Client scripting: small vanilla JavaScript only where browser APIs are needed,
  such as copy-to-clipboard, keyboard handling, and persisted UI preferences.
- Data storage: none for MVP; optional local config file later for bookmarks,
  favorites, and preferences.

## Runtime Configuration

Serverigh should be configured through flags and environment variables.

Suggested flags:

- `--root`
  Default: current working directory. Filesystem root that Serverigh may expose.
- `--host`
  Default: `127.0.0.1`. Bind host.
- `--port`
  Default: `4173`. Bind port.
- `--write`
  Default: `false`. Enables mutating file operations.
- `--show-hidden`
  Default: `false`. Initial hidden-file visibility.
- `--max-preview-bytes`
  Default: `1048576`. Maximum bytes read for text-like preview.

Suggested environment aliases:

- `SERVERIGH_ROOT`
- `SERVERIGH_HOST`
- `SERVERIGH_PORT`
- `SERVERIGH_WRITE`
- `SERVERIGH_MAX_PREVIEW_BYTES`

## Routing

- `GET /`: redirect or render the current root folder.
- `GET /browse?path=...`: render the full browser page for a folder.
- `GET /partials/files?path=...`: render only the file table for HTMX
  navigation.
- `GET /partials/breadcrumbs?path=...`: render breadcrumbs for the active path.
- `GET /preview?path=...`: render a file preview partial.
- `GET /download?path=...`: stream a file download.
- `GET /raw?path=...`: stream raw file content when browser-safe.
- `GET /healthz`: health check.

Future write-mode routes:

- `POST /actions/mkdir`: create a folder.
- `POST /actions/rename`: rename a file or folder.
- `POST /actions/move`: move selected entries.
- `POST /actions/copy`: copy selected entries.
- `POST /actions/delete`: delete selected entries.
- `POST /actions/upload`: upload files to the active folder.

## UI Layout

The main screen should be a dense file-explorer workspace:

- Top bar: root selector, current path, search field, view controls.
- Left/sidebar area: folder tree or quick roots.
- Main area: sortable file table.
- Preview pane: selected file preview, metadata, and actions.
- Footer/status row: item count, selected item count, root boundary, and
  read-only/write mode.

HTMX behavior:

- Folder clicks update the file table, breadcrumbs, preview pane, and URL.
- File clicks update the preview pane and URL.
- Sorting and filtering update the file table without full reload.
- Back and forward navigation restore the folder/file view.

## Preview Rules

Preview type detection should use extension plus content sniffing where cheap.

- Markdown: render safe HTML with headings, lists, links, code blocks,
  blockquotes, and pipe tables.
- CSV: render an HTML table with sampled width estimation and horizontal
  scrolling.
- Text/code: show escaped text with truncation notice for large files.
- JSON: pretty-print when valid; fall back to text.
- Images: render with natural dimensions and fit controls.
- PDF: use browser-native embed/object preview.
- Audio/video: use browser-native media controls.
- Binary/unknown: show metadata, mime guess, size, modified time, and download
  action.

Large-file policy:

- Never read a whole large file just to decide whether it can be previewed.
- Use `os.Open` plus bounded reads for preview.
- Use streaming responses for downloads and raw file serving.

## Filesystem Safety

Path safety is a core requirement.

- Resolve all requested paths against the configured root.
- Reject traversal attempts that escape root.
- Reject or carefully resolve symlinks that escape root.
- Bind to `127.0.0.1` by default.
- Default to read-only mode.
- Hide write controls unless write mode is enabled.
- Confirm destructive operations.
- Enforce preview byte limits.
- Avoid logging full sensitive paths unless debug logging is enabled.

## Error Handling

User-facing errors should be specific and calm:

- File not found.
- Permission denied.
- Path is outside configured root.
- Preview is too large.
- File type cannot be previewed.
- Operation requires write mode.

Server errors should wrap context in Go and log enough detail for local
debugging.

## Testing Plan

Unit tests:

- Path resolver prevents traversal and symlink escape.
- Directory listing sorts folders/files correctly.
- Markdown pipe-table renderer emits table view models.
- CSV parser handles quoted cells, multiline rows, wide rows, and empty fields.
- Preview type detector handles common extensions and unknown files.

Handler tests:

- `/browse` renders root and nested folders.
- `/preview` renders correct preview partials.
- Read-only mode blocks write routes.

Manual checks:

- Large CSV remains readable with horizontal scroll.
- Markdown tables render as tables.
- Reload restores active folder/file.
- Browser back/forward restores previous folder/file.
- Hidden-file toggle does not resize the layout unexpectedly.

## MVP Acceptance Criteria

- App starts with one command and serves on localhost.
- User can browse inside the configured root.
- User cannot escape the configured root through `..`, encoded paths, or symlink
  tricks.
- Directory view shows name, type, size, modified time, and permissions.
- Clicking folders updates the view without a full page reload.
- Clicking files renders the correct preview without a full page reload.
- Markdown files render pipe tables correctly.
- CSV files with many columns remain readable.
- Large text previews do not load entire files into memory.
- Active folder/file survives reload through URL state.
- App is read-only unless explicitly started in write mode.
