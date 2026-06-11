# Git & Pull Request Rulebook

This document defines the strict version control conventions for the **Konex** repository. Following these rules ensures a clean, searchable git history and makes code reviews efficient and effective.

> **How to use this document:** When in doubt, follow the rule. When a rule blocks you, open a PR (or an issue) to change the rule — don't silently work around it. A consistent history is worth more than any individual shortcut.

---

## 📑 Table of Contents

1. [Commit Conventions](#-commit-conventions)
2. [Branching Strategy](#-branching-strategy)
3. [Pull Request Conventions](#-pull-request-conventions)
4. [PR Description Template](#-pr-description-template)
5. [Code Review Etiquette](#-code-review-etiquette)
6. [Merging](#-merging)
7. [What Never Goes In Git](#-what-never-goes-in-git)
8. [Definition of Done](#-definition-of-done-checklist)

---

## 🌳 Commit Conventions

To maintain absolute clarity on exactly what changed and why, we enforce the following strict commit rules.

### 1. One File Per Commit

- **Strict rule:** Commit **one file per commit**.
- **Why?** This guarantees that if a bug is introduced, tracking it down via `git bisect` or reverting a specific change is surgical and won't drag along unrelated changes.
- **Exception to be sane about:** A file and its tightly-bound generated counterpart (e.g., a snapshot or a lockfile that _must_ move with it) may share a commit when separating them would break the build at that commit. Note the reason in the message body.

### 2. Descriptive Prefixes

Every commit message must start with a standardized prefix followed by a clear, descriptive summary of the change.

**Allowed prefixes:**

- `feat:` — A new feature
- `fix:` — A bug fix
- `refactor:` — A code change that neither fixes a bug nor adds a feature
- `chore:` — Changes to the build process or auxiliary tools/libraries
- `style:` — Changes that don't affect meaning (whitespace, formatting, semicolons)
- `test:` — Adding missing tests or correcting existing ones
- `docs:` — Documentation-only changes
- `hotfix:` — An urgent fix pushed directly to production/main
- `perf:` — A code change that improves performance
- `typo:` — Fixing typos in code or documentation

**Example commit messages:**

- `feat: add Google OAuth integration to sender service`
- `fix: resolve null pointer exception in campaign batcher`
- `style: format items handler to match 150-line rule`

### 3. Writing the Message

- **Imperative mood, present tense** — write `add`, not `added` or `adds`. The message should complete the sentence "If applied, this commit will…".
- **Keep the summary line under ~72 characters.** It's the one line everyone reads in `git log --oneline`.
- **Lowercase after the prefix**, no trailing period on the summary line.
- **Describe the _why_ in the body** when the change isn't self-evident. Wrap the body at ~72 characters and separate it from the summary with a blank line.
- **One concern per commit.** If you find yourself writing "and" in a commit summary, it's probably two commits.

```text
feat: add retry backoff to email sender

The Gmail API intermittently returns 429s during a blast. Without a
backoff the batcher dropped those leads silently. This adds an
exponential backoff capped at 3 attempts before marking the lead failed.

Closes #58
```

---

## 🌿 Branching Strategy

- **Never commit directly to `main`** (except `hotfix:` where the process allows it, and even then prefer a fast-tracked PR).
- **Branch off `main`** and name branches with the same prefix vocabulary as commits:
  - `feat/google-oauth`
  - `fix/batcher-null-pointer`
  - `chore/upgrade-pocketbase`
- **Keep branches short-lived.** Rebase on `main` regularly to avoid painful merges. A branch that lives for weeks is a merge conflict waiting to happen.
- **One branch = one PR = one logical change.** Don't pile unrelated work onto a single branch.

---

## 🔀 Pull Request Conventions

Pull Requests are not just for merging code; they are the permanent documentation of _why_ a decision was made. Every PR must be thoroughly documented.

- **Small and reviewable.** A PR should represent one logical change. If it's growing past a few hundred lines of meaningful diff, consider splitting it.
- **Self-review first.** Read your own diff before requesting review — half the comments you'd get, you can catch yourself.
- **Green before review.** The PR must pass formatting, linting, type-checks, and tests before a human is asked to look at it.
- **The title follows the same prefix convention** as commits, e.g. `feat: add lead CRUD endpoints`.
- **Visual changes include a screenshot or recording** so reviewers can verify the result without checking out the branch.

---

## 📝 PR Description Template

Your PR description **must** include the following headings and answer the associated questions.

### 📝 Description

A high-level overview of the feature, bug fix, or refactor.

**Answer here:**

- **What does this do?** (e.g., _Adds a new Lead table and the associated CRUD endpoints in the backend._)
- **Why did you do this?** (e.g., _Users need a way to store their network contacts before launching an email blast._)

### 🔄 Changes

A bulleted list of the exact technical changes made. Because you commit one file at a time, this should align perfectly with your commit history.

**Answer here:**

- **Who/what does this impact?** (e.g., _Impacts the database schema and adds a new route to the frontend router. Does not affect existing email templates._)

### 🧪 Testing

**Answer here:**

- **How did you test this?** (e.g., _Tested locally by creating 5 leads through the UI and verifying they appear in the PocketBase Admin UI._)

### 📎 Closes

List any related issues or tickets this PR resolves.

- _Closes #42_
- _Fixes #15_

> **Tip:** Save this as `.github/pull_request_template.md` so the template auto-populates on every new PR.

---

## 👀 Code Review Etiquette

**For the author:**

- Respond to every comment — resolve it with a change, or explain why not. Don't silently dismiss.
- Push fixes as new commits during review (don't force-push mid-review) so reviewers can see what changed.

**For the reviewer:**

- Review promptly; a stale PR blocks the author and grows conflicts.
- Be specific and kind — comment on the code, not the coder. Suggest, don't command.
- Distinguish blocking concerns from nits; prefix optional suggestions with `nit:`.
- Approve only when you'd be comfortable owning the change yourself.

---

## ✅ Merging

- **All conversations resolved** and at least one approval before merge.
- **CI is green** — no merging on red or "I'll fix it after."
- **Prefer squash merge** so `main` reads as one clean commit per PR, while the branch keeps the granular one-file-per-commit history.
- **Delete the branch** after merge to keep the remote tidy.
- **The merger is responsible** for confirming the PR still applies cleanly to current `main`.

---

## 🚫 What Never Goes In Git

- **Secrets** — API keys, tokens, Google credentials, `.env` files. Commit a `.env.example` with placeholders instead.
- **Generated artifacts & data** — `pb_data/`, build output, `node_modules/`. Keep `.gitignore` authoritative.
- **Large binaries** — use proper storage, not the repo.
- **Commented-out dead code** — delete it; git remembers it for you.

> If a secret is ever committed, treat it as compromised: rotate the credential immediately, then scrub history. Removing it in a later commit is **not** enough.

---

## ✅ Definition of Done (Checklist)

A change is ready to merge when:

- [ ] Each commit touches one file and starts with an approved prefix.
- [ ] Commit messages are imperative, lowercase, and under ~72 chars.
- [ ] The branch is named with the prefix convention and rebased on current `main`.
- [ ] The PR title follows the prefix convention.
- [ ] The PR description fills out Description, Changes, Testing, and Closes.
- [ ] Visual changes include a screenshot or recording.
- [ ] Formatting, lint, type-checks, and tests are all green.
- [ ] No secrets, generated data, or dead code in the diff.
- [ ] All review conversations are resolved with at least one approval.
