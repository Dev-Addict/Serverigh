# Aria AzadiPour AI Coding Style Guide

Purpose: give AI coding agents a practical style profile to use when writing or
modifying code in Aria AzadiPour's style.

Audit basis: public GitHub account `Dev-Addict` reviewed on 2026-06-05. The
GitHub API snapshot returned 130 public owner repositories, 124 non-forks, and 6
forks. The non-fork repositories were shallow-cloned locally for source
inspection. This guide intentionally weights recent authored code more heavily
than older tutorial or learning projects.

## Best Document Type and Usage

Use this as a Markdown agent-instruction document. Markdown is better than PDF,
DOCX, or prose notes because coding agents can ingest it directly, it works in
source control, and it can be copied into bot-specific instruction files without
conversion.

Recommended setup:

1. Keep this file as the canonical personal style guide.
2. For Codex and other repo-aware agents, place a short `AGENTS.md` in the repo
   root that links to or embeds this guide.
3. For Claude Code, mirror the same instructions into `CLAUDE.md`.
4. For GitHub Copilot, copy the durable rules into
   `.github/copilot-instructions.md`.
5. For Cursor, make a thin `.cursor/rules/aria-style.mdc` wrapper that points to
   this guide.

Instruction priority for agents:

1. Follow the target repository's existing instructions, formatter, framework,
   and tests first.
2. Use this guide as the fallback style when creating new files or when local
   patterns are ambiguous.
3. Do not preserve accidental rough edges from older repos, such as committed
   dependency folders, stale debug logs, or placeholder TODO comments.

Suggested one-paragraph prompt:

> Write code in Aria AzadiPour's style. Follow the local repo conventions first.
> Prefer small domain-oriented modules with explicit role suffixes, typed data
> boundaries, practical feature-first implementation, project formatter
> settings, and focused tests for logic-heavy changes. Avoid over-abstracting,
> avoid large generic files, and do not reproduce old tutorial-code artifacts.

## Core Style Profile

Aria's durable style is pragmatic, modular, and builder-oriented. He tends to
create complete working features rather than framework-heavy skeletons. The
strongest signal across repositories is clear separation by responsibility:
components, hooks, types, enums, models, services, helpers, utils, controllers,
store slices, and domain folders.

The dominant public-code stack is TypeScript and JavaScript, especially React,
React Native, Next.js, Vite, Node, Express, GraphQL, and small frontend tools.
Newer work adds Rust, Go, Lua configuration, Deno/Hono, Biome, and systems-style
CLI/TUI projects.

When imitating the style, favor the newer version of it:

- Typed TypeScript over loose JavaScript when the project supports it.
- Small files named by role over broad catch-all modules.
- Direct implementation over elaborate architecture.
- Project tools and configs over personal preference.
- Clear domain naming over clever abstractions.

## Formatting Defaults

Always inspect and follow local formatter config first.

Common TypeScript/React fallback from many `.prettierrc` files:

- Tabs for indentation, with `tabWidth: 2`.
- Semicolons enabled.
- Single quotes for JavaScript and TypeScript.
- Respect Prettier's `bracketSpacing` setting; do not force spaces inside object
  braces when it is `false`.
- Trailing commas set to `es5`.
- Parentheses around arrow function parameters.
- Compact imports such as `import {FC, useMemo} from 'react';` appear in several
  projects.

Important exception: newer Biome-configured code, especially
`Deliverers/client`, uses tabs and double quotes. If Biome says double quotes,
use double quotes.

Rust and Go:

- Let `rustfmt` and `gofmt` own formatting.
- Use idiomatic module, package, and file naming for each language.
- Prefer clear `Result` and error-returning APIs over hidden panics.

## Naming and File Organization

Use explicit role suffixes. This is one of the clearest repeatable patterns.

Common frontend/backend suffixes:

- `.component.tsx` and `.component.ts`
- `.hook.ts`
- `.type.ts`
- `.enum.ts`
- `.helper.ts`
- `.util.ts`
- `.service.ts`
- `.model.ts`
- `.slice.ts`
- `.context.ts`
- `.provider.tsx`

Organize files by domain and role. Good examples of folder shapes:

- `components/drivers/modals/create.modal.tsx`
- `components/orders/table/columns.tsx`
- `services/auth/helpers/sign-in.helper.ts`
- `types/enums/order-status.enum.ts`
- `store/slices/drivers.slice.ts`
- `utils/validators/password.validator.ts`
- `handlers/message-handlers/print-color-scheme/...`

Naming tendencies:

- React components and classes use `PascalCase`.
- Hooks start with `use`.
- Functions, helpers, and variables use `camelCase`.
- TypeScript interfaces are direct and local: `Props`, `FormValues`, `Options`,
  or domain-specific names.
- Enums and constants are separated into their own files when shared.
- Rust modules and files use `snake_case`, with domain-specific modules like
  `directory_entry`, `input_mode`, `read_directory`, `file_result`.
- Go packages are short and plain, with internal folders for server/database
  concerns when useful.

## Frontend Style

Prefer functional React components with hooks in newer code. Older projects
include class components, but new work should use functions unless the existing
repo is class-based.

Common frontend habits:

- Local `Props` or `FormValues` interfaces near the component that consumes
  them.
- `useCallback` for event handlers and async actions passed into UI components.
- `useMemo` for modal renderers, derived column definitions, and heavier UI
  fragments.
- Context/provider files for shared state when Redux or another store is not the
  right fit.
- Store slices for app-level state in Redux-based projects.
- Typed API response and payload shapes near the API or domain folder.
- UI components are composed from simple primitives, styled-components, or the
  active UI library.

When building UI:

- Make the actual feature screen first, not a marketing wrapper.
- Keep components small and name them by role.
- Use existing UI libraries and styling systems before inventing new primitives.
- Keep props explicit and typed.
- Avoid generic "magic" component factories unless the existing repo already
  uses them.

## Backend Style

Backend code usually separates controllers/routes, services, models, helpers,
validators, and error utilities.

Common backend habits:

- Controllers parse input, validate basic shape, call models/services, and
  return JSON.
- Services compose smaller helper functions.
- Auth and validation logic gets split into helpers and tested directly.
- Errors often use explicit domain strings or app error classes when the project
  has them.
- Middleware is used for cross-cutting concerns like auth, errors, logging, and
  panic recovery.
- Environment/config handling is explicit instead of hidden in globals.

For Node/TypeScript servers:

- Keep route/controller files readable and direct.
- Put reusable domain behavior in `services/.../helpers`.
- Put request/response types under `types`, `payloads`, or domain files.
- Validate incoming body/query values before mutating data.

For Go services:

- Use `cmd`, `internal`, `business`, `foundation`, or similarly clear service
  layers when the project already has them.
- Prefer structured logging, health endpoints, graceful shutdown, and Makefile
  commands.
- Return errors with context using wrapping where useful.

For Rust:

- Break features into domain modules.
- Use enums and `match` for state machines and event handling.
- Use explicit `Result` aliases and custom error modules.
- Keep CLI/TUI app state explicit: app structs, input modes, events, windows,
  widgets, commands.

## Testing Style

Tests are present in focused pockets rather than every repository. Imitate the
useful habit, not the inconsistency.

Add tests when changing:

- Parsers, converters, validators, and algorithms.
- Auth helpers, JWT/session logic, and permission checks.
- Models and database utilities.
- API route behavior with non-trivial validation.
- Go packages with public behavior.

Testing conventions seen:

- Jest-style `describe`/`it` for JS/TS projects.
- Go `_test.go` files for Go packages.
- Focused helper/model/service tests in backend TypeScript.
- Algorithm tests in learning projects.

Do not create a large new test framework just to add one small change. Use the
project's existing runner and style.

## Error Handling and Logging

Prefer explicit error paths.

- In Rust, return `Result` and convert errors into domain-specific error types.
- In Go, return `error`, wrap context, and keep shutdown/server errors visible.
- In JS/TS servers, use existing `AppError`, middleware, or JSON error response
  conventions.
- In frontend code, expose `isLoading` and `isError` state for async flows.
- Use optimistic UI updates only when rollback behavior is clear.

Avoid adding casual `console.log` calls. Existing older repos include debug
logs, but new code should use the project's logger, toast, error middleware, or
test assertions instead.

## Comments and Documentation

Keep comments sparse and useful.

Good comments:

- Explain non-obvious command sequences, shutdown flows, or formatter/linter
  exceptions.
- Mark a real TODO with a concrete next action.
- Document public service methods if the surrounding project already uses JSDoc.

Avoid comments that restate code. Avoid placeholder linter comments such as
`<explanation>`; if a linter suppression is necessary, write the actual reason.

## Repository Hygiene

For new code, improve on the public-history rough edges:

- Do not commit `node_modules`, generated build folders, or IDE metadata unless
  the repo explicitly requires it.
- Keep README/setup instructions current when adding commands.
- Prefer one lockfile per package manager.
- Use `package.json`, `Makefile`, `Cargo.toml`, `go.mod`, or project-native
  commands as the source of truth.
- Keep generated files marked or isolated.

## AI Bot Checklist

Before editing:

- Inspect formatter config, package scripts, test commands, and nearby files.
- Identify whether the repo is older tutorial-style code or newer typed/modular
  code.
- Match local naming and folder patterns.

While editing:

- Create small role-named files.
- Keep data types close to their domain.
- Use explicit async/loading/error states.
- Prefer direct implementation over broad abstraction.
- Add focused tests for logic-heavy behavior.
- Preserve project-specific quote, semicolon, tab, and import style.

Before finishing:

- Run the narrowest relevant formatter/test/build command.
- Remove accidental debug logs.
- Check that new file names follow the role-suffix convention.
- Mention any command you could not run.

## Do and Do Not

Do:

- Use typed boundaries for API payloads, props, context values, models, and
  enums.
- Split reusable logic into helpers or utils with role suffixes.
- Keep UI components direct and feature-oriented.
- Follow project formatters exactly.
- Prefer practical local state, context, or Redux slices depending on the
  existing project.
- Use idiomatic Rust and Go when in those ecosystems.

Do not:

- Generate large generic architecture layers that the repo does not need.
- Collapse many roles into one broad file.
- Add unexplained linter suppressions.
- Copy old committed dependency-folder habits.
- Add noisy comments or debug logging.
- Override local repo conventions with this guide.

## Evidence Snapshot

Local audit artifacts:

- `github_audit/repository_audit_summary.csv`: CSV inventory of 130 public owner
  repositories from the GitHub API snapshot.
- `github_audit/repos/`: shallow source snapshot of 124 non-fork public
  repositories.

High-signal repositories for this guide:

- `Arfima`: recent Rust CLI/TUI file manager with modular
  app/window/widget/input structure.
- `RRacone`: Rust interpreter project with scanner/parser/error/result modules.
- `Invaders`: Rust TUI game.
- `Go-Service`: Go HTTP service base with structured config, logging, shutdown,
  Makefile, and deployment folders.
- `Siget` and `GoIgnore`: smaller Go projects with package tests.
- `Deliverers`: newer TypeScript/Deno/React delivery dashboard with Biome, Hono,
  Redux, TanStack Router, modal/table/form components, and domain services.
- `jreel`, `apricity`, `figma-coloria`: TypeScript UI/tooling projects showing
  component/type/enum/helper naming.
- `tactical-edge-project-back-end`: TypeScript backend with auth services,
  helpers, validators, models, GraphQL, and focused tests.
- Older JS/React/Node projects such as `online-schools`, `portfolio`,
  `react-testing-starter`, and `blockchain`: useful for recurring
  controller/helper/test patterns, but lower weight for modern style.

Confidence:

- High for TypeScript/JavaScript/React/Node style.
- Medium-high for newer Rust and Go style.
- Medium for testing habits.
- Low for Java/Android style because those public repos are older and less
  representative of current style.

## Markdown

Use 80 character line length for markdown files
