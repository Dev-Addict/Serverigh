# Serverigh First Version Plan

## Purpose

Build the first production-quality version of Serverigh as a local-first web
file explorer. The product should stay small, fast, and understandable: Go for
the server and filesystem layer, Fiber for HTTP, Go templates for HTML, HTMX
for partial updates, and minimal vanilla JavaScript for browser-only behavior.

This plan treats the first version as the first usable release, not a full
future product. The release should make read-only browsing, safe previewing,
and URL-backed navigation feel complete before adding broad write workflows.

## Source Documents

- `docs/feature-list.md`
- `docs/project-specification.md`

Keep this plan aligned with those files when scope changes.

## Current Baseline

The repository already has a useful foundation:

- Go module and Makefile commands for build, test, vet, format, and run.
- CLI flags and environment variables for root, host, port, write mode,
  hidden files, and max preview bytes.
- Fiber server setup with route registration.
- Go template rendering with embedded templates.
- Handler package split by route.
- Filesystem service for root resolution, directory listing, file metadata,
  bounded previews, raw file serving, and downloads.
- Operational errors with stable application error codes.
- Unit and handler tests for the existing server, routes, and filesystem
  behavior.

The baseline is still not a product-quality first release because the main UI,
HTMX navigation, preview rendering, breadcrumbs, sorting, search, and manual
browser checks are not complete.

## First Version Scope

Ship these as required first-version capabilities:

- Root-scoped browsing with traversal and symlink escape prevention.
- Read-only behavior by default.
- Directory table with folders first, file metadata, empty states, and errors.
- HTMX folder navigation without full page reloads.
- Preview pane updates for selected files without full page reloads.
- URL state for the active folder and active file.
- Breadcrumb navigation for the active path.
- Text and code preview with bounded reads and truncation notices.
- Markdown preview with safe HTML and readable pipe tables.
- CSV preview with horizontal scrolling and readable wide-table behavior.
- Binary file details with raw and download actions where appropriate.
- Native previews for images, PDFs, audio, and video where browser support is
  enough.
- Hidden-file behavior that is explicit and stable.
- Health endpoint for smoke checks.
- Focused tests for filesystem safety, handlers, previews, and URL behavior.

Defer these until after the first version unless the read-only release gates
are already met:

- Mutating file actions such as rename, move, copy, delete, mkdir, and upload.
- Trash-first delete behavior.
- Multi-root bookmarks.
- Bulk selection and archive downloads.
- File watching.
- Git-aware folders and diff views.
- Remote filesystem adapters.
- AI indexing or semantic search.

## Milestone 1: Stabilize The Baseline

Goal: make the current foundation reliable enough to extend.

Tasks:

- Run `make check` and fix any failing tests, vet warnings, or tidy changes.
- Confirm every filesystem boundary error is an operational error.
- Confirm static assets are served correctly from embedded or stable paths.
- Add missing tests around root config normalization and invalid paths.
- Review all existing routes and remove any stale route assumptions from docs.
- Keep Markdown files wrapped at 80 characters.

Acceptance criteria:

- `make check` passes.
- No known `fmt.Errorf` application errors remain.
- Running `make run` starts the app on `127.0.0.1:4173`.
- The app can list the configured root and reject invalid paths.

## Milestone 2: Build The Explorer Shell

Goal: replace the current minimal page with the real file explorer workspace.

Tasks:

- Create a dense desktop-first layout with top bar, file table, preview pane,
  and status row.
- Render the configured root, current path, mode, item count, and root boundary.
- Add an empty preview state for folders and initial page loads.
- Add loading and error states for partial updates.
- Make the layout responsive enough to remain usable on narrow screens.
- Avoid a marketing or landing page. The first screen is the explorer.

Acceptance criteria:

- `/browse?path=/` renders the actual explorer.
- Empty folders, inaccessible paths, and missing paths have clear UI states.
- Text does not overflow buttons, table cells, panels, or status areas.

## Milestone 3: Wire HTMX Navigation

Goal: make browsing feel like an app while keeping server-rendered HTML.

Tasks:

- Convert folder links to update the file table, breadcrumbs, status row, and
  preview pane through HTMX.
- Convert file links to update only the preview pane and URL state.
- Implement `/partials/breadcrumbs`.
- Add a stable response contract for partials and operational errors.
- Use `history.pushState` or HTMX history support for folder and file views.
- Restore folder and file state on reload, back, and forward navigation.

Acceptance criteria:

- Folder clicks do not trigger a full page reload.
- File clicks update the preview pane without replacing the explorer shell.
- Reloading a URL restores the same folder or file preview.
- Browser back and forward restore previous folder and file views.

## Milestone 4: Complete Core Previewers

Goal: make the files Serverigh is built for readable in the browser.

Tasks:

- Add preview type detection using extension plus cheap content sniffing.
- Keep bounded reads for text-like previews.
- Render plain text and common code/config files safely.
- Render Markdown through a safe renderer with headings, lists, code blocks,
  links, blockquotes, and pipe tables.
- Render CSV with a real parser, preserved quoting behavior, and horizontal
  scrolling for wide tables.
- Pretty-print valid JSON and fall back to text for invalid JSON.
- Use native browser rendering for images, PDFs, audio, and video.
- Show metadata and download actions for binary or unsupported files.

Acceptance criteria:

- Large text previews are truncated without reading the whole file.
- Markdown tables render as tables.
- Wide CSV files remain readable.
- Binary files do not produce broken text output.
- Preview tests cover text, Markdown, CSV, JSON, media, binary, and truncation.

## Milestone 5: Add Table Controls

Goal: make the directory table useful for real project folders.

Tasks:

- Add sorting by name, size, modified time, and created time when known.
- Add client-visible sort state in URLs or HTMX request parameters.
- Add copy actions for relative and absolute paths.
- Add keyboard basics: arrow navigation, enter to open, and backspace to parent.

Acceptance criteria:

- Sorting is deterministic and tested.
- Hidden files are never shown by accident when the option is disabled.
- Copy actions use small controls and work without a frontend build step.

## Milestone 6: Search Within Root

Goal: provide useful search while keeping safety and performance predictable.

Tasks:

- Start with filename search under the configured root.
- Add result limits and clear truncation messaging.
- Exclude hidden files unless hidden-file mode is enabled.
- Add bounded text-content search only after filename search is stable.
- Avoid indexing or long-running background work in the first version.

Acceptance criteria:

- Search cannot escape the configured root.
- Search results link back into the explorer and preview pane.
- Large folders do not freeze the UI or consume unbounded memory.

## Milestone 7: Package The First Release

Goal: make the product easy to run, test, and inspect.

Tasks:

- Add a short `README.md` with install, run, flags, and safety notes.
- Add example commands for common roots and preview-byte limits.
- Document read-only defaults and write-mode limitations.
- Add a release checklist that includes automated and manual checks.
- Build a local binary with `make build`.

Acceptance criteria:

- A new user can run the app from the README without reading the source.
- `make test`, `make vet`, and `make build` pass.
- Manual checks cover traversal, symlink boundaries, HTMX navigation, reload,
  back and forward navigation, Markdown tables, CSV width, and large previews.

## Testing Strategy

Keep tests close to the behavior they protect:

- Filesystem tests live in `internal/filesystem`.
- Handler tests live beside each handler in `internal/httpserver/handlers`.
- Route tests live in `internal/httpserver`.
- CLI tests live in `cmd/serverigh`.
- Preview parser tests should be table-driven and file-type specific.

Add broader integration tests only when behavior crosses package boundaries,
such as URL state, HTMX partial contracts, and full browse-to-preview flows.

## Release Gates

Do not tag the first version until all gates pass:

- `make check` passes.
- The app binds to `127.0.0.1` by default.
- The configured root is visible in the UI.
- Traversal and symlink escape attempts are rejected.
- Read-only mode exposes no mutating controls.
- Core previews work for Markdown, CSV, text, JSON, image, PDF, and binary
  files.
- Large files are handled with bounded reads or streaming.
- URLs preserve active folder and file state.
- Browser back, forward, and reload work for normal navigation.
- Error responses are stable, specific, and based on operational error codes.

## Scope Control

The first version should stay focused. When choosing between polish and broad
new surface area, prefer the choice that makes browsing, previewing, and path
safety more reliable.

Write actions should remain disabled or clearly unfinished until the read-only
workflow is solid. Search should start with filenames. Remote adapters,
multi-user auth, and AI indexing should wait until the local single-user
product is excellent.
