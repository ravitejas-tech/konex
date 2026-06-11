# Frontend Rulebook: React Router v7 + Tailwind + Framer Motion

This document defines the architectural guidelines, styling rules, and coding standards for the **Konex** frontend. The goal is to build a highly modular, strictly typed, and visually stunning web application — and to give every contributor (human or AI) an unambiguous source of truth for _where code goes_ and _how it should look_.

> **How to use this document:** When in doubt, follow the rule. When a rule blocks you, open a PR to change the rule — don't silently work around it. Consistency beats individual preference.

---

## 📑 Table of Contents

1. [Visuals & Aesthetics](#-visuals--aesthetics)
2. [Architecture & Component Structure](#-architecture--component-structure)
3. [Routing: React Router v7](#-routing-react-router-v7)
4. [Data Fetching: React Query Kit](#-data-fetching-react-query-kit)
5. [Forms & Validation](#-forms--validation)
6. [State Management](#-state-management)
7. [TypeScript Standards](#-typescript-standards)
8. [Styling Rules](#-styling-rules)
9. [Animation Guidelines](#-animation-guidelines)
10. [Accessibility](#-accessibility)
11. [Performance](#-performance)
12. [Error & Loading States](#-error--loading-states)
13. [Folder Structure](#-folder-structure)
14. [Git & PR Conventions](#-git--pr-conventions)
15. [Definition of Done](#-definition-of-done-checklist)

---

## 🎨 Visuals & Aesthetics

The application must deliver a premium, buttery-smooth user experience. Visual excellence is a top priority.

- **Motion-rich by default** — The app leans on smooth micro-interactions, page transitions, and element reveals. Motion should feel _intentional and physical_, never gratuitous.
- **Component Libraries**
  - Leverage **[21st.dev](https://21st.dev/community/components)** for modern, high-quality, pre-built animated components.
  - Use **Framer Motion** for all custom animations.
- **Styling** — **Tailwind CSS**. Avoid writing custom CSS unless absolutely necessary.
- **Typography** — The **Poppins** font family is used globally across all headings and body text.
- **Polish is part of "done."** A feature that works but feels janky is not finished. Spacing, alignment, hover/focus states, and transitions all count.

---

## 🏗️ Architecture & Component Structure

The frontend follows a modular architecture, aligning with **React Router v7** patterns.

### Principles

- **Dependencies flow one way** — Routes compose features; features compose UI primitives. A generic `ui/` component must never import a feature or a route. If a button knows about "leads," it's not a UI primitive.
- **Co-locate by feature, not by type** — Keep a feature's components, queries, and helpers near each other. Reach for a global folder only when something is genuinely shared.
- **One component, one job** — A component either _renders_ or _orchestrates_, rarely both. Pull data-fetching and logic up; keep leaf components presentational.

### 1. File & Folder Conventions

- **Kebab Case** — ALL files and folders use `kebab-case` (e.g., `user-profile.tsx`, `auth-context.ts`, `button-base.tsx`). No exceptions.
- **Component naming** — Files are `kebab-case`; the React component they export is `PascalCase` (`user-profile.tsx` → `export function UserProfile()`).
- **One component per file** — Plus its tightly-coupled tiny subcomponents if they're never used elsewhere.

### 2. File Length & Granularity

- **Small Files** — A single file should rarely exceed **150 lines**.
- **Chunking** — If a component file goes over 150 lines, it _must_ be broken into smaller sub-components placed in an adjacent `components/` folder specific to that route or module.
- **Reusability** — Generic UI elements (buttons, inputs, cards) live in a central `components/ui/` directory. Don't reinvent a primitive that already exists there.

---

## 🧭 Routing: React Router v7

- **Folder-based routes** live in `app/routes/`. Keep route files thin — they wire data and layout together, delegating real UI to feature components.
- **Loaders & actions** are the front door for route data. Use a route `loader` for read-on-navigation data and an `action` for form submissions/mutations where it fits the RR7 model; use React Query Kit for client-driven, cache-shared, or frequently-refetched data.
- **Co-locate route components** — A route's private subcomponents go in a `components/` folder next to the route (inside that route's folder), not in the global tree.
- **Type your route modules** — Use the generated route types (`Route.LoaderArgs`, `Route.ComponentProps`) rather than hand-writing param/loader types.
- **Boundaries per route** — Export an `ErrorBoundary` for routes that can fail, so one broken page never blanks the whole app.

---

## 📡 Data Fetching: React Query Kit

We use `react-query-kit` for strictly typed, modular API interactions. **Do not use standard `@tanstack/react-query` hooks (`useQuery`, `useMutation`) directly.**

### Rules for Queries

- **Centralized location** — All queries and mutations are defined in a dedicated `queries/` folder (e.g., `queries/connections.ts`, `queries/email-templates.ts`).
- **Strictly `react-query-kit`** — Use ONLY `createQuery`, `createMutation`, `createInfiniteQuery`, etc.
- **Usage** — Define the query in `queries/`, export the generated hook, import it into your component.
- **Stable, structured query keys** — Let the kit derive keys from variables; never hand-roll ad-hoc string keys that drift between files.
- **Invalidate, don't refetch by hand** — After a mutation, invalidate the relevant query keys in `onSuccess` so the cache stays the single source of truth.
- **API calls live in `api/`, not in components** — `queries/` wraps the functions from `api/`; components only ever touch the generated hooks.

**Example (`queries/connections.ts`):**

```typescript
import { createQuery, createMutation } from "react-query-kit";
import { fetchConnections, addConnection } from "../api/connections";

export const useConnectionsQuery = createQuery({
  queryKey: ["connections"],
  fetcher: fetchConnections,
});

export const useAddConnectionMutation = createMutation({
  mutationFn: addConnection,
  onSuccess: (_data, _vars, ctx) => {
    ctx?.client.invalidateQueries({ queryKey: useConnectionsQuery.getKey() });
  },
});
```

---

## 📝 Forms & Validation

All forms must be strictly typed and validated to ensure a bulletproof user experience.

- **Validation** — **Zod** for all schema definitions and type inference.
- **Form Handling** — **React Hook Form** (`react-hook-form`) for form state, wired to Zod via `@hookform/resolvers/zod`.
- **Strict typing** — Infer TypeScript types from Zod schemas (`z.infer<typeof schema>`); never declare the form type separately from its schema.
- **Schema is the source of truth** — Validate the same schema on submit; surface field-level errors from RHF, not ad-hoc `useState` flags.
- **Disable on submit** — Wire submit buttons to `formState.isSubmitting` to prevent double submits, and show inline pending feedback.

**Example:**

```typescript
const leadSchema = z.object({
  email: z.string().email(),
  name: z.string().min(1, "Name is required"),
});

type LeadForm = z.infer<typeof leadSchema>;

const form = useForm<LeadForm>({ resolver: zodResolver(leadSchema) });
```

---

## 🧠 State Management

- **Server state ≠ UI state.** Anything that comes from the API is _server state_ — it belongs to React Query Kit, not `useState`. Don't copy fetched data into local state.
- **Local UI state** (open/closed, hovered, current tab) stays in the component with `useState`/`useReducer`.
- **Shared UI state** (auth user, theme) goes in a small, typed React Context — kebab-cased (`auth-context.ts`). Keep contexts narrow; one giant global store is an anti-pattern here.
- **URL is state too.** Filters, tabs, and pagination that should survive refresh/share belong in search params, not component state.

---

## 🔒 TypeScript Standards

- **`strict` mode on**, always. No `any` — reach for `unknown` and narrow, or define the type properly.
- **No non-null assertions (`!`)** to silence the compiler; handle the nullable case.
- **Infer over annotate** where the type is obvious; annotate explicitly at public boundaries (props, query fetchers, exported functions).
- **Props are typed objects**, defined right above the component. Prefer discriminated unions over a pile of optional booleans.
- **No prop drilling past 2–3 levels** — lift to context or compose components instead.

---

## 🎨 Styling Rules

- **Tailwind utility-first.** Compose utilities in the markup; avoid custom CSS files unless a utility genuinely can't express it.
- **Design tokens, not magic numbers.** Use Tailwind's scale (spacing, colors, radius) configured in `tailwind.config.ts`. Don't sprinkle arbitrary `[13px]` values.
- **Tame long class lists.** When `className` strings get unwieldy or conditional, use a `cn()`/`clsx` + `tailwind-merge` helper in `lib/utils.ts` — never string-concatenate classes by hand.
- **Variants via a single source.** Build component variants (e.g., button sizes/intents) with one variant map, not scattered conditionals.
- **Mobile-first & responsive** by default — base styles target small screens, layer up with `md:`/`lg:`.
- **Dark mode aware** where the design calls for it; rely on the configured strategy, not one-off overrides.

---

## ✨ Animation Guidelines

- **Framer Motion for custom motion**, 21st.dev for pre-built animated components.
- **Respect `prefers-reduced-motion`.** Gate non-essential animation so motion-sensitive users get a calm experience.
- **Animate cheap properties** — `transform` and `opacity`, not `width`/`height`/`top` — to keep 60fps.
- **Consistent timing.** Define shared durations/easings (e.g., in a `lib/motion.ts`) so the whole app feels like one product, not ten different ones.
- **Purposeful, not decorative.** Motion should guide attention, signal state, or smooth a transition. If it does none of those, cut it.

---

## ♿ Accessibility

- **Semantic HTML first** — a `<button>` is a button. Don't rebuild interactive elements out of `<div>`s.
- **Keyboard reachable** — every interactive element is focusable and operable by keyboard, with visible focus states.
- **Label everything** — inputs have associated labels; icon-only buttons have `aria-label`.
- **Color is never the only signal** — pair it with text/iconography; meet contrast minimums.

---

## ⚡ Performance

- **Code-split by route** — lean on RR7's route-level splitting; lazy-load heavy, below-the-fold, or rarely-used components.
- **Memoize deliberately** — `useMemo`/`useCallback`/`React.memo` where a real, measured cost exists, not reflexively.
- **Stable keys in lists** — use real IDs, never array indices.
- **Don't ship the world** — watch bundle size; import individual icons/utilities, not entire libraries.
- **Optimize images & fonts** — preload Poppins, serve appropriately sized images.

---

## 🔄 Error & Loading States

- **Every async surface has three states**: loading, error, and empty — design all three, not just the happy path.
- **Loading** — prefer skeletons over spinners for content; keep layout stable to avoid shift.
- **Error** — show a human message and a retry affordance; route-level failures use the route `ErrorBoundary`.
- **Empty** — a deliberate empty state with guidance ("No leads yet — add your first"), never a blank screen.

---

## 📂 Folder Structure

The frontend follows a modular architecture. Generic, shared code lives in the global directories; each route lives in its own folder with its private components co-located inside it.

```text
web/
├── app/
│   ├── components/              # Global components
│   │   ├── ui/                  # Reusable UI primitives (buttons, inputs, cards)
│   │   └── layout/              # Navbars, sidebars, page wrappers
│   ├── queries/                 # React Query Kit definitions
│   │   ├── auth.ts
│   │   ├── connections.ts
│   │   ├── email-templates.ts
│   │   └── sender-accounts.ts
│   ├── api/                     # Raw API client functions (fetchers) consumed by queries/
│   │   ├── connections.ts
│   │   ├── email-templates.ts
│   │   └── sender-accounts.ts
│   ├── routes/                  # React Router v7 routes — folder per route, components co-located
│   │   ├── _index/
│   │   │   └── route.tsx                     → /
│   │   │
│   │   ├── connections/
│   │   │   ├── route.tsx                     → /connections  (listing page)
│   │   │   ├── components/
│   │   │   │   ├── connection-filters.tsx
│   │   │   │   └── connection-list.tsx
│   │   │   └── $connectionId/
│   │   │       ├── route.tsx                 → /connections/:connectionId  (detail page)
│   │   │       └── components/
│   │   │           ├── connection-details.tsx
│   │   │           └── connection-activity.tsx
│   │   │
│   │   ├── email-templates/
│   │   │   ├── route.tsx                     → /email-templates  (listing page)
│   │   │   ├── components/
│   │   │   │   └── template-grid.tsx
│   │   │   └── $templateId/
│   │   │       ├── route.tsx                 → /email-templates/:templateId
│   │   │       └── components/
│   │   │           ├── template-editor.tsx
│   │   │           └── template-preview.tsx
│   │   │
│   │   ├── sender-accounts/
│   │   │   ├── route.tsx                     → /sender-accounts
│   │   │   └── components/
│   │   │       ├── sender-account-card.tsx
│   │   │       └── connect-account-dialog.tsx
│   │   │
│   │   └── campaigns/
│   │       └── $campaignId/
│   │           ├── route.tsx                 → /campaigns/:campaignId
│   │           └── components/
│   │               └── campaign-overview.tsx
│   ├── context/                 # Narrow, typed React contexts (auth, theme)
│   │   └── auth-context.ts
│   └── lib/                     # Utilities, Zod schemas, motion presets, cn() helper
│       ├── utils.ts
│       ├── motion.ts
│       └── schemas.ts
├── tailwind.config.ts
└── tsconfig.json
```

**Rules**

- **Folder per route.** Every route is its own folder containing a `route.tsx` file. Dynamic segments use the `$param` convention as a nested folder (e.g., `connections/$connectionId/route.tsx` → `/connections/:connectionId`).
- **Co-locate route components.** A route's private subcomponents live in a `components/` folder _inside that route's folder_, never in the global tree. Nested routes get their own nested `components/` folder.
- `components/ui/` holds _only_ domain-agnostic primitives. The moment a component knows about a domain, it belongs in the relevant route's `components/` folder.
- `queries/` wraps `api/`; components import from `queries/`, never call `api/` directly.
