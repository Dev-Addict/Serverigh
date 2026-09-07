# Serverigh

Serverigh is a local-first web file explorer written in Go. It exposes a
configured filesystem root through a browser UI so developers and operators can
browse folders, inspect file metadata, preview text-like files, and download or
view raw files without leaving the browser.

The product is designed to stay small: Fiber for HTTP, Go `html/template` for
server-rendered HTML, embedded static assets, HTMX-ready partial routes, and
minimal vanilla JavaScript where browser APIs are needed.

Serverigh binds to `127.0.0.1` by default and is read-only unless write mode is
explicitly enabled.

## Current Status

Serverigh is in early first-version development. The backend foundation and
main explorer shell are in place, including CLI configuration, filesystem
browsing, safe path handling, bounded previews, HTMX folder/file navigation,
Markdown/CSV/JSON/code/media previewers, raw/download responses, and
operational errors.

Search, table controls, and write workflows are tracked below as unchecked
items.

## Feature Checklist

Implemented:

- [x] Fiber HTTP server with CLI and environment configuration.
- [x] Configurable root, host, port, write mode, hidden files, and preview
  limits.
- [x] Embedded templates and static assets.
- [x] Root-scoped browsing with traversal and symlink escape protection.
- [x] Directory listing with folder-first sorting and bounded result counts.
- [x] File metadata, including size, mode, MIME type, creation time, and
  modification time.
- [x] Bounded text preview with truncation and binary-file handling.
- [x] Markdown, CSV, JSON, code/config, media, and binary preview states.
- [x] Raw and download routes with validated file handles and safety headers.
- [x] Operational error codes with JSON or HTML responses as appropriate.
- [x] Explorer shell with top bar, file table, preview pane, and status row.
- [x] HTMX navigation for folder browsing, file preview updates, and
  breadcrumbs.
- [x] URL-backed state for reload, back, and forward navigation.
- [x] Table sorting, config-driven hidden-file visibility, and copy-path
  actions.
- [x] Health endpoint, Makefile workflow, and backend test coverage.

Planned for the first version:

- [ ] Search within the configured root, starting with filenames.
- [ ] Keyboard navigation and improved metadata views.
- [ ] Opt-in write-mode workflows with confirmations and trash-first delete.
- [ ] Bulk selection, archive downloads, and file refresh behavior.

Later or explicitly deferred:

- [ ] Bookmarks, favorites, recent files, tagging, and notes.
- [ ] Git-aware folders, text diffs, and archive browsing.
- [ ] Remote filesystem adapters and plugin previewers.
- [ ] AI indexing, multi-user authentication, and browser terminal execution.

## Installation

Serverigh currently builds from source.

Requirements:

- Go `1.23` or newer.
- `make`, for the provided development commands.

Build the binary:

```sh
make build
```

The binary is written to `bin/serverigh`.

## Usage

Run against the current directory:

```sh
make run
```

Run against a specific root:

```sh
make run ROOT=/path/to/workspace
```

Run the compiled binary:

```sh
bin/serverigh --root /path/to/workspace
```

Then open:

```text
http://127.0.0.1:4173
```

## Configuration

Serverigh can be configured with CLI flags or environment variables.

- `--root`
  Environment: `SERVERIGH_ROOT`. Default: current directory. Filesystem root.
- `--host`
  Environment: `SERVERIGH_HOST`. Default: `127.0.0.1`. Bind host.
- `--port`
  Environment: `SERVERIGH_PORT`. Default: `4173`. Bind port.
- `--write`
  Environment: `SERVERIGH_WRITE`. Default: `false`. Enables write mode.
- `--show-hidden`
  Environment: `SERVERIGH_SHOW_HIDDEN`. Default: `false`. Shows hidden files.
- `--max-preview-bytes`
  Environment: `SERVERIGH_MAX_PREVIEW_BYTES`. Default: `1048576`. Preview byte
  limit.

## Routes

Current routes:

- `GET /`
  Redirects to the root browse view.
- `GET /browse?path=...`
  Renders the full browse page for a folder.
- `GET /partials/files?path=...`
  Renders the directory listing partial.
- `GET /partials/breadcrumbs?path=...`
  Renders the breadcrumb partial.
- `GET /preview?path=...`
  Renders the file preview partial.
- `GET /raw?path=...`
  Streams raw file content with active-content safety handling.
- `GET /download?path=...`
  Streams a file download.
- `GET /healthz`
  Returns a health payload for smoke checks.

Reserved but not implemented yet:

- `POST /actions/mkdir`
- `POST /actions/rename`
- `POST /actions/move`
- `POST /actions/copy`
- `POST /actions/delete`
- `POST /actions/upload`

## Safety Model

Serverigh is intentionally local-first and root-scoped.

- It binds to `127.0.0.1` by default.
- It resolves requested paths against the configured root.
- It rejects traversal attempts.
- It rejects symlinks that escape the configured root.
- It limits preview reads with `--max-preview-bytes`.
- It streams raw and download responses from validated file handles.
- It serves active raw content, such as HTML and SVG, as plain text.
- It sends `X-Content-Type-Options: nosniff` for raw and download responses.
- It keeps mutating actions unavailable until write workflows are implemented.

Do not expose Serverigh to the public internet unless authentication,
authorization, and transport security have been designed for that deployment.

## Development

Common commands:

```sh
make help
make fmt
make tidy
make test
make vet
make check
make build
make clean
```

`make check` runs formatting, module tidy, vet, and tests.

## Testing

The current test suite covers:

- CLI parsing.
- Config normalization.
- Operational errors.
- Filesystem path safety.
- Directory listing.
- Hidden-file behavior.
- Creation-time behavior.
- Bounded previews.
- HTTP route wiring.
- Handler responses.
- Raw and download safety headers.
- Embedded static asset serving.

Run:

```sh
make test
```

## Project Documents

Additional planning and specification documents live in:

- `docs/feature-list.md`
- `docs/project-specification.md`
- `.agents/PLANS.md`
