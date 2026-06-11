# Backend Rulebook: Go + PocketBase Architecture

This document defines the architectural guidelines and coding standards for the **Konex** backend. It exists to keep the codebase clean, modular, highly readable, and maintainable as it scales — and to give every contributor (human or AI) an unambiguous source of truth for _where code goes_ and _how it should look_.

> **How to use this document:** When in doubt, follow the rule. When a rule blocks you, open a PR to change the rule — don't silently work around it. Consistency beats individual preference.

---

## 📑 Table of Contents

1. [Core Philosophy](#-core-philosophy-modular--decoupled)
2. [Directory Structure](#-directory-structure)
3. [Layer Responsibilities](#-layer-responsibilities)
4. [Request Lifecycle](#-request-lifecycle-the-golden-path)
5. [Coding Standards](#-coding-standards)
6. [Error Handling](#-error-handling)
7. [Configuration & Secrets](#-configuration--secrets)
8. [Concurrency & Background Work](#-concurrency--background-work)
9. [Database & Migrations](#-database--migrations)
10. [Testing](#-testing)
11. [Logging & Observability](#-logging--observability)
12. [Security](#-security)
13. [Git & PR Conventions](#-git--pr-conventions)
14. [Definition of Done](#-definition-of-done-checklist)

---

## 🏗️ Core Philosophy: Modular & Decoupled

The backend follows a strictly **Domain-Driven Modular Architecture**. Instead of stuffing logic into giant files or a massive `main.go`, the application is split into small, focused, independently testable packages organized by business domain (e.g., `leads`, `templates`, `campaigns`, `accounts`).

### Principles

- **Small Files, Small Functions** — Files should rarely exceed 200–300 lines. Functions should do _one_ thing well (Single Responsibility Principle). If a function grows past 30–40 lines, extract helpers.
- **Separation of Concerns** — HTTP parsing (Handlers) must be strictly separated from business logic (Services) and data access (Models/Repositories).
- **Readability over Cleverness** — Write explicit, boring Go. Avoid deep abstraction layers and "clever" one-liners that the next reader has to decode.
- **Dependencies Point Inward** — Handlers depend on Services; Services depend on Repositories/Models. Nothing inner ever imports something from an outer layer. A `service` must never import a `handler`.
- **Make the Right Thing the Easy Thing** — If contributors keep breaking a rule, the structure is wrong. Fix the structure, not the people.

---

## 📂 Directory Structure

The `internal/` directory is the heart of the application. It is organized by technical layer, and within each layer, divided by domain.

```text
api/
├── cmd/
│   └── api/
│       └── main.go              # Entry point: init PocketBase, load env, register hooks/routes, app.Start()
├── internal/
│   ├── config/                  # Environment variables & app-wide configuration (loaded once, injected everywhere)
│   ├── database/
│   │   ├── migrations/          # Programmatic PocketBase schema definitions & seed data
│   │   └── schema/              # Collection name + field name constants (single source of truth for strings)
│   ├── handlers/                # HTTP layer (Controllers)
│   │   ├── routes.go            # Central router registration & middleware wiring
│   │   ├── leads_handler.go     # HTTP logic for Leads
│   │   └── campaign_handler.go  # HTTP logic for Campaigns
│   ├── hooks/                   # PocketBase event hooks (OnRecordCreate, etc.)
│   │   ├── hooks.go             # RegisterHooks(app) — single registration entry point
│   │   └── email_hooks.go       # Lifecycle hooks, decoupled from handlers
│   ├── models/                  # DTOs (request/response structs) — type safety + JSON tags only
│   │   ├── lead_dto.go
│   │   └── template_dto.go
│   ├── repositories/            # (Optional) DB query logic when it grows beyond trivial Record calls
│   │   └── lead_repo.go
│   ├── services/                # Core business logic
│   │   ├── batcher_service.go   # Logic for the 1 email / 10 min cadence
│   │   └── sender_service.go    # Google auth & email dispatch
│   └── pkg/                     # Cross-cutting helpers (errors, validation, response writers) — no domain logic
│       ├── apierror/
│       └── response/
├── pb_data/                     # Auto-managed SQLite DB (git-ignored)
└── pb_migrations/               # JS migrations if the Admin UI generates any (prefer Go migrations)
```

**Rules**

- `internal/` is enforced privacy. Nothing outside `api/` can import it — keep it that way.
- `pkg/` holds _only_ domain-agnostic utilities. The moment a helper knows about "leads," it belongs in a domain service.
- Never reference a collection or field by a raw string literal. Define it once in `database/schema/` and import the constant. This makes renames a one-line change and kills typo bugs.

---

## 🧩 Layer Responsibilities

### 1. `handlers/` — Controllers / Routers

**Responsibility:** Translate between HTTP and the domain. Nothing more.

- **No business logic.** A handler may only: parse/bind the request (`e.BindBody()`), validate input, call a service, and format the response.
- **Granular files.** Split by domain (`leads_handler.go`, `auth_handler.go`) — never one mega `handlers.go`.
- **Central routing.** Register all routes in `routes.go`, where middleware (e.g., `apis.RequireAuth()`) is applied in a single, auditable place.
- **Constructor injection.** Build handlers as structs that receive their service dependencies, e.g. `NewLeadsHandler(leadSvc *services.LeadService)`. No reaching into globals.

### 2. `services/` — Business Logic

**Responsibility:** Execute the rules of Konex.

- **Own the rules.** Cadence, validation beyond shape, orchestration across domains — all here.
- **Typed errors.** Return clear, wrapped Go errors (see [Error Handling](#-error-handling)). Never return a raw DB error to a handler.
- **Compose, don't duplicate.** If a service needs another domain, inject and call that service rather than copying logic.
- **No HTTP, no `core.RequestEvent`.** A service must never know it was triggered by an HTTP request. Pass plain data in, return plain data/errors out — this keeps it callable from hooks, workers, and tests alike.

### 3. `models/` — Schemas & DTOs

**Responsibility:** Define the shape of data crossing the boundary.

- Hold request payloads (`CreateLeadRequest`) and response bodies (`LeadResponse`) with JSON + validation tags.
- **No query execution.** Models are for type safety and serialization only. DB access lives in repositories/services.
- Keep request and response structs distinct — never reuse one struct for both input and output.

### 4. `repositories/` — Data Access (optional but encouraged as it grows)

**Responsibility:** Encapsulate all PocketBase/SQLite query logic.

- Start with PocketBase's `app.FindRecordById` etc. inline in services. **Once a domain accumulates non-trivial or repeated queries, extract them into a repository** so services stay focused on rules, not SQL/Dao plumbing.
- Repositories return models or `core.Record`s — never HTTP types.

### 5. `hooks/` — PocketBase Lifecycle

**Responsibility:** React to record lifecycle events.

- Register everything via a single `RegisterHooks(app core.App)` — keep `main.go` clean.
- **Keep hooks fast.** If a hook kicks off something long-running (an email blast), hand it to a goroutine/worker; never block the response. See [Concurrency](#-concurrency--background-work).
- Hooks should delegate to services, not contain business logic themselves.

### 6. `database/` — Source of Truth for Schema

**Responsibility:** Define collections programmatically.

- Write Go migrations in `database/migrations/` so collections (`leads`, `templates`, `campaigns`) are created consistently across all environments — never rely on hand-edits in the Admin UI as the source of truth.

---

## 🛤️ Request Lifecycle: The Golden Path

Every request should flow in exactly one direction. If you find yourself skipping a layer, stop and reconsider.

```text
HTTP Request
   │
   ▼
[Handler]   bind body → validate shape → map to service input
   │
   ▼
[Service]   apply business rules → orchestrate → call repo(s)
   │
   ▼
[Repository / PocketBase Dao]   read/write records
   │
   ▼
[Service]   build domain result
   │
   ▼
[Handler]   map result → DTO → write JSON response
   │
   ▼
HTTP Response
```

**Never** let a handler touch the database directly, and **never** let a service format an HTTP response.

---

## 🚦 Coding Standards

1. **Naming**
   - `camelCase` for local variables and functions.
   - `PascalCase` for exported structs, interfaces, and methods.
   - Files in `snake_case.go`.
   - Interfaces describe behavior (`Sender`, `LeadFinder`), not implementations.
2. **Formatting** — `gofmt`/`goimports` is non-negotiable and enforced in CI. No manually formatted code.
3. **Linting** — `golangci-lint` must pass with zero warnings before merge.
4. **Context** — Pass `context.Context` as the first argument to any function that does I/O, and respect cancellation.
5. **Dependency Injection** — Pass the PocketBase `core.App` (and other deps) explicitly into services and hooks. No global mutable state. This is what makes the code testable.
6. **No magic values** — Timeouts, batch sizes (e.g., the 10-minute cadence), and limits live in `config/`, named and documented.
7. **Comments explain _why_, not _what_** — The code already says what it does; comment the non-obvious reasoning.

---

## ❗ Error Handling

- **Always handle errors explicitly.** Never `_ =` an error you haven't consciously decided to ignore (and if you do ignore one, comment why).
- **No `panic`** outside one-time initialization in `main.go`.
- **Wrap with context** as errors travel up:
  ```go
  if err != nil {
      return fmt.Errorf("send email to lead %s: %w", leadID, err)
  }
  ```
- **Sentinel / typed errors** for cases the caller must branch on:
  ```go
  var ErrLeadNotFound = errors.New("lead not found")
  // caller: if errors.Is(err, services.ErrLeadNotFound) { ... }
  ```
- **One place maps errors to HTTP.** A small helper in `pkg/apierror` translates domain errors → status codes so handlers stay clean and responses stay consistent. Services never set HTTP status codes.

---

## ⚙️ Configuration & Secrets

- Load all config **once** at startup in `internal/config/` into a typed struct; inject it. No scattered `os.Getenv()` calls deep in the codebase.
- **Fail fast.** If a required env var (Google credentials, DB path) is missing, the app should refuse to start with a clear message — not crash later mid-request.
- Secrets never get logged, committed, or returned in API responses. `.env` is git-ignored; commit a `.env.example` with placeholder keys.

---

## ⚡ Concurrency & Background Work

- **Never block the HTTP response** on slow work (email sends, third-party API calls). Acknowledge fast, process in the background.
- Every goroutine must have a clear owner and a way to stop — respect `context` cancellation and app shutdown.
- **Recover from panics in background goroutines** so one failed job can't take down the process:
  ```go
  go func() {
      defer func() {
          if r := recover(); r != nil {
              app.Logger().Error("worker panic", "err", r)
          }
      }()
      // ...work...
  }()
  ```
- The email cadence (1 / 10 min) is a domain rule — implement it in `batcher_service.go`, configurable via `config/`, not hard-coded in a hook.

---

## 🗄️ Database & Migrations

- **Go migrations are the source of truth.** Every schema change ships as a migration so environments stay in lockstep.
- Migrations are **forward-only and append-only** once merged — never edit a migration that has run in production; write a new one.
- Reference collection and field names through constants in `database/schema/`, never raw strings.
- `pb_data/` is git-ignored and disposable in development; production data is backed up separately.

---

## 🧪 Testing

- **Services are the priority target.** Because they hold business logic and have no HTTP/DB coupling, they should be unit-testable with fakes/mocks of their dependencies.
- Use table-driven tests — the idiomatic Go pattern — for functions with multiple input cases.
- Handlers get lightweight tests for binding/validation and correct status codes.
- A bug fix lands with a test that fails before the fix and passes after.
- Aim for meaningful coverage of business rules over a vanity percentage.

---

## 📈 Logging & Observability

- Use PocketBase's structured logger (`app.Logger()`); log key/value pairs, not formatted prose.
- Log at boundaries: request received (debug), business decisions (info), recoverable issues (warn), failures (error).
- **Never log secrets or full PII.** Log identifiers (lead IDs), not raw email bodies or tokens.
- Every error log should carry enough context (IDs, operation) to debug without a reproduction.

---

## 🔐 Security

- **Validate all input at the handler boundary** before it reaches a service.
- Apply `apis.RequireAuth()` (and role checks) centrally in `routes.go`; default to _locked down_, open up explicitly.
- Trust nothing from the client for authorization — re-derive the acting user server-side.
- Keep dependencies patched; run `govulncheck` in CI.

---

## ✅ Definition of Done (Checklist)

A change is done when:

- [ ] Logic lives in the correct layer (handler / service / repo / model).
- [ ] No file exceeds ~300 lines; no function exceeds ~40.
- [ ] Errors are wrapped with context; no swallowed errors; no stray `panic`.
- [ ] No raw collection/field strings — constants used.
- [ ] No secrets logged, committed, or returned.
- [ ] Config values are injected, not hard-coded.
- [ ] Slow work runs in the background, not on the request path.
- [ ] Tests cover the new business rules and pass.
- [ ] `gofmt` + `golangci-lint` are clean.
