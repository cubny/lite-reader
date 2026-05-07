# Implementation Plan: Lite Reader Frontend Modernization

## Overview

Replace the legacy jQuery 1.x frontend with a Preact + htm + signals static SPA, embedded in the Go binary via `//go:embed`. Backend stays JSON-only. Greenfield rewrite, single PR, ordered commits. Full reference: `SPEC.md` and `/Users/cubny/.claude/plans/ideate-how-we-can-eager-pudding.md`.

## Architecture Decisions

- **No build step.** Vendored ESM modules served as-is. Browser import map resolves bare specifiers.
- **Preact (10.24.3) + htm (3.1.1) + @preact/signals (1.3.0)** — pinned, vendored at `internal/web/public/vendor/`.
- **Embedded FS** via `//go:embed all:public` in new `internal/web/web.go`. `PUBLIC_DIR` env var overrides for dev hot-reload.
- **Two PRs**: PR-A = Phases 1+2 (foundation + embed FS, legacy UI still served from embed). PR-B = Phases 3–7 (full Preact rewrite + cleanup). Smaller reviewable units, de-risks the move.
- **Vendor path**: `internal/web/public/assets/vendor/...`, served at `/assets/vendor/...`.
- **Untrusted item HTML**: rendered inside `<iframe sandbox srcdoc>` in `ItemDetail`. Backend does no sanitization (`internal/app/feed/service.go:91` stores raw `gofeed` Content). No DOMPurify needed.
- **Vertical slicing** — each phase delivers one working flow end-to-end (auth → feeds → items).
- **Co-located tests** — `Foo.test.js` next to `Foo.js`. Vitest aliases bare specifiers to `node_modules`.
- **`data-testid` everywhere** — no CSS-class or text-based selectors.

---

## Task List

### Phase 1: Foundation (no UI changes visible to users yet)

#### Task 1: Vendor Preact + htm + signals; add import map smoke test
**Description:** Download pinned ESM builds, place under `internal/web/public/vendor/`, wire an import map into a throwaway smoke HTML, confirm the browser executes a `<script type="module">` that imports from `preact` and renders a `Hello` node.

**Acceptance criteria:**
- [ ] `internal/web/public/assets/vendor/preact@10.24.3/{preact,hooks,signals}.module.js` exist (signals subdir under preact for grouping is fine).
- [ ] `internal/web/public/assets/vendor/htm@3.1.1/htm.module.js` exists.
- [ ] Import map paths use `/assets/vendor/...`.
- [ ] A temporary `internal/web/public/_smoke.html` (deleted at end of phase) imports `preact` + `htm`, renders to `#app`, and logs success.
- [ ] Versions match what `package.json` will pin in Task 2.

**Verification:**
- [ ] `curl http://localhost:3000/_smoke.html` renders the Preact node when served via `make run`.
- [ ] Browser console shows zero 404s and zero module-resolution errors.

**Dependencies:** None.
**Files touched:** `internal/web/public/vendor/**`, `internal/web/public/_smoke.html` (temp). **Scope: S.**

---

#### Task 2: Vitest + testing-library setup
**Description:** Add dev deps to `package.json`, create `vitest.config.js` at repo root with jsdom + alias map mirroring the import map. Add a one-line passing sanity test under `internal/web/public/js/_setup.test.js`.

**Acceptance criteria:**
- [ ] `package.json` adds: `vitest`, `@testing-library/preact`, `@testing-library/jest-dom`, `@testing-library/user-event`, `jsdom`, `msw`.
- [ ] Scripts: `"test:unit": "vitest run"`, `"test:unit:watch": "vitest"`.
- [ ] `vitest.config.js` aliases `preact`, `preact/hooks`, `@preact/signals`, `htm` to `node_modules` versions matching the vendored pins.
- [ ] One sanity test passes.

**Verification:**
- [ ] `npm install` succeeds.
- [ ] `npm run test:unit` passes with the sanity test.

**Dependencies:** Task 1 (versions must match).
**Files touched:** `package.json`, `package-lock.json`, `vitest.config.js`, `internal/web/public/js/_setup.test.js`. **Scope: S.**

---

#### Task 3: API client + token plumbing
**Description:** Create `internal/web/public/js/api/client.js` with `request()`, `getToken/setToken/clearToken`, `Authorization: Bearer` injection, 401 → clear + redirect, `ApiError {status, message}`. Co-located test covers all branches.

**Acceptance criteria:**
- [ ] `client.js` exports `request(method, path, opts)`, `ApiError`, `getToken`, `setToken`, `clearToken`.
- [ ] Token read from / written to `localStorage.token`.
- [ ] On 401 (when `auth: true`), token is cleared and `location.assign('/login.html')` invoked.
- [ ] Non-2xx throws `ApiError` with `status` + parsed `message`.

**Verification:**
- [ ] `npm run test:unit -- client` — all branches covered (200, 401, 4xx, 5xx, no-token, with-token).
- [ ] `make lint` clean (no ESLint configured yet, this is a no-op for now).

**Dependencies:** Task 2.
**Files touched:** `internal/web/public/js/api/client.js`, `internal/web/public/js/api/client.test.js`. **Scope: S.**

---

### Checkpoint: Foundation
- [ ] `make build` still passes (legacy UI untouched).
- [ ] `npm run test:unit` passes.
- [ ] Browser smoke page renders Preact via vendored modules.
- [ ] **Pause for human review** before continuing.

---

### Phase 2: Embed FS + Go integration

#### Task 4: Move `public/` → `internal/web/public/`; create `web.go` with embed
**Description:** Move the (still legacy) `public/` directory to `internal/web/public/` and create `internal/web/web.go` exposing `FS() fs.FS` from `//go:embed all:public`.

**Acceptance criteria:**
- [ ] Directory moved (`git mv` to preserve history).
- [ ] `internal/web/web.go` compiles; `web.FS()` returns a non-nil `fs.FS` rooted at `public/`.
- [ ] Vendored modules from Task 1 are inside the new path.
- [ ] No reference to old `public/` path remains in Go code, Makefile, Dockerfile, Playwright config, or `.golangci.yml`.

**Verification:**
- [ ] `go build ./...` succeeds.
- [ ] `grep -rn '"public"' --include='*.go' --include=Makefile --include=Dockerfile` returns only `internal/web/web.go`.

**Dependencies:** Task 1 (vendor lives under the moved tree).
**Files touched:** `internal/web/web.go` (new); rename `public/` → `internal/web/public/`; updates to `Makefile`, `Dockerfile`, `playwright.config.js`, `.golangci.yml` if any reference the old path. **Scope: M.**

---

#### Task 5: Thread `fs.FS` through router; add `PUBLIC_DIR` escape hatch
**Description:** Modify `internal/infra/http/api/router.go` to accept `staticFS fs.FS` and use `http.FileServer(http.FS(staticFS))` for `NotFound`. In `internal/app.go`, pass `web.FS()` by default, or `os.DirFS(env)` when `PUBLIC_DIR` is set.

**Acceptance criteria:**
- [ ] `router.New(...)` (or its options struct) accepts `staticFS fs.FS`.
- [ ] `internal/app.go` reads `PUBLIC_DIR` env; falls back to `web.FS()`.
- [ ] Existing tests still pass (router tests may need a fixture FS).

**Verification:**
- [ ] `make build && ./lite-reader` serves `/index.html` (still legacy content) on :3000 with **no `public/` folder present** in cwd.
- [ ] `PUBLIC_DIR=./internal/web/public ./lite-reader` serves the same content from disk.
- [ ] `make test` passes.

**Dependencies:** Task 4.
**Files touched:** `internal/infra/http/api/router.go`, `internal/infra/http/api/router_test.go`, `internal/app.go`, `cmd/main.go` (if it constructs the router directly), `cmd/testserver/main.go` (verify, likely no change). **Scope: M.**

---

### Checkpoint: Embedded binary
- [ ] `./lite-reader` (single binary, no `public/` folder beside it) serves the legacy UI correctly end-to-end.
- [ ] `make test-ui` passes against the legacy UI from the embedded FS (sanity check — Playwright shouldn't care).
- [ ] **Pause for human review.**

---

### Phase 3: Auth vertical slice

#### Task 6: `api/auth.js` + tests
**Description:** Implement `login`, `signup`, `logout` against `POST /login`, `POST /signup`. `login` stores returned token via `client.setToken`; `logout` calls `client.clearToken` + redirects.

**Acceptance criteria:**
- [ ] `auth.js` exports `login({email,password})`, `signup({email,password})`, `logout()`.
- [ ] Successful login persists token; failed login throws `ApiError`.
- [ ] All branches covered by `auth.test.js` using `vi.fn()` for `fetch`.

**Verification:** `npm run test:unit -- auth`.
**Dependencies:** Task 3.
**Files touched:** `internal/web/public/js/api/auth.js` + test. **Scope: S.**

---

#### Task 7: `LoginPage` + new `login.html`
**Description:** Replace legacy `login.html` with a minimal HTML shell + `js/login.js` entry that mounts `<LoginPage/>`. `LoginPage` renders the form, calls `api/auth.login`, redirects to `/index.html` on success, displays inline error otherwise. Every input/button has `data-testid`.

**Acceptance criteria:**
- [ ] `login.html` includes the import map and loads `js/login.js`.
- [ ] `LoginPage.js` renders `email`, `password`, submit, error region — all with `data-testid`.
- [ ] Submitting valid creds → redirect to `/`.
- [ ] Submitting invalid creds → inline error visible, no redirect.
- [ ] `LoginPage.test.js` covers happy path + error path with mocked `fetch`.

**Verification:**
- [ ] `npm run test:unit -- LoginPage`.
- [ ] Manual: `make run` → visit `/login.html` → log in with seeded creds → reach `/index.html` (still legacy app for now, but token is in localStorage).

**Dependencies:** Tasks 5, 6.
**Files touched:** `internal/web/public/login.html`, `internal/web/public/js/login.js`, `internal/web/public/js/pages/LoginPage.js` + test, `internal/web/public/js/util/html.js`. **Scope: M.**

---

#### Task 8: `SignupPage` + new `signup.html`
**Description:** Same shape as Task 7 for signup. Submits to `api/auth.signup`, then auto-logs-in or redirects to login (match current behavior — confirm by reading `public/js/signup.js` from git history).

**Acceptance criteria:**
- [ ] `signup.html` + `js/signup.js` + `pages/SignupPage.js` + test.
- [ ] Validation parity with current behavior (email format, password min length).
- [ ] Test covers happy path + error path.

**Verification:** `npm run test:unit -- SignupPage`; manual signup creates a user.
**Dependencies:** Task 7 (shared `util/html.js` and CSS).
**Files touched:** `internal/web/public/signup.html`, `internal/web/public/js/signup.js`, `internal/web/public/js/pages/SignupPage.js` + test, `internal/web/public/css/auth.css`. **Scope: M.**

---

### Checkpoint: Auth slice
- [ ] Auth e2e specs (`tests/ui/auth.spec.js`) updated to use new `data-testid` selectors and pass via `make test-ui`.
- [ ] `make test-all` passes.
- [ ] **Pause for human review.**

---

### Phase 4: App shell

#### Task 9: `state.js`, `router.js`, `util/html.js`, `util/time.js`
**Description:** Create global signals (`token`, `feeds`, `selection`, `items`, `currentItem`, `unreadCount`, `starredCount`, `toast`), a ~30-line hash router exporting `currentRoute` signal + `navigate(path)`, and small util modules.

**Acceptance criteria:**
- [ ] `state.js` exports each signal as a named export.
- [ ] `router.js` parses `#/`, `#/starred`, `#/feed/:id`, `#/feed/:id/item/:itemId` into `{name, params}`.
- [ ] Tests cover route parsing edge cases.

**Verification:** `npm run test:unit -- router state`.
**Dependencies:** Task 2.
**Files touched:** `internal/web/public/js/state.js` + test, `js/router.js` + test, `js/util/{html,time,dom}.js` + tests for time. **Scope: M.**

---

#### Task 10: `AppPage` shell with 3-pane CSS grid + `Resizer`
**Description:** Replace legacy `index.html` body with `<div id="app"></div>` + import map + `js/main.js` mounting `<AppPage/>`. `AppPage` renders three empty `<section>`s in a CSS grid driven by `--sidebar-w`, `--list-w`. `Resizer` updates the variables on drag, persists to localStorage. On boot, if no token → redirect to `/login.html`.

**Acceptance criteria:**
- [ ] Empty 3-pane shell renders, panes labeled with `data-testid="pane-sidebar|pane-list|pane-detail"`.
- [ ] Resizer drag adjusts widths smoothly; reload restores last sizes.
- [ ] No token → redirect.
- [ ] Tests: `AppPage.test.js` (auth gate), `Resizer.test.js` (drag math + localStorage).

**Verification:**
- [ ] `make run` → log in → see empty 3-pane shell.
- [ ] `npm run test:unit -- AppPage Resizer`.

**Dependencies:** Tasks 7, 9.
**Files touched:** `internal/web/public/index.html`, `js/main.js`, `js/pages/AppPage.js` + test, `js/components/Resizer.js` + test, `css/app.css`. **Scope: M.**

---

### Checkpoint: Shell live
- [ ] Login → empty shell visible. Smoke spec (`tests/ui/smoke.spec.js`) updated and passes.
- [ ] **Pause for human review.**

---

### Phase 5: Feeds vertical slice

#### Task 11: `api/feeds.js` + tests
**Acceptance criteria:** exports `list`, `add`, `remove`, `fetchNew`, `markRead`, `markUnread`, `items(id, opts)`. All methods covered.
**Verification:** `npm run test:unit -- feeds`.
**Dependencies:** Task 3. **Files touched:** `js/api/feeds.js` + test. **Scope: S.**

---

#### Task 12: `Sidebar` + `FeedList` + `FeedItem`
**Description:** Sidebar wraps logout + (placeholders for SmartFolders + AddFeedForm) + `FeedList`. `FeedList` reads `state.feeds` and renders `FeedItem` per row. `FeedItem` shows title, unread count, delete button. Click selects feed → `navigate('#/feed/:id')`. Delete → confirm → `feeds.remove` → refresh.

**Acceptance criteria:**
- [ ] Three components + co-located tests.
- [ ] `data-testid`: `feed-list`, `feed-item`, `feed-item-delete`, `feed-item-title`, `feed-item-unread-count`.
- [ ] Tests assert: render with mock signal, click selects, delete confirms then calls API.

**Verification:** `npm run test:unit -- Sidebar FeedList FeedItem`.
**Dependencies:** Tasks 10, 11.
**Files touched:** `js/components/Sidebar.js`, `FeedList.js`, `FeedItem.js`, plus tests. **Scope: M.**

---

#### Task 13: `AddFeedForm`
**Description:** URL input + submit; calls `api/feeds.add`, refreshes `state.feeds`, clears input. Shows inline validation/error.
**Acceptance criteria:** valid URL adds; invalid URL shows error; pending state disables submit. Test covers all three.
**Verification:** `npm run test:unit -- AddFeedForm`.
**Dependencies:** Task 12.
**Files touched:** `js/components/AddFeedForm.js` + test. **Scope: S.**

---

#### Task 14: `SmartFolders`
**Description:** "Unread" + "Starred" entries in the sidebar showing live counts from `state.unreadCount`/`starredCount`. Click → `navigate('#/')` or `#/starred`. Counts refreshed via `api/items.unreadCount/starredCount` on mount and after relevant actions.
**Acceptance criteria:** counts render, refresh on signal change, navigation works. Test covers count display + click.
**Verification:** `npm run test:unit -- SmartFolders`.
**Dependencies:** Task 12.
**Files touched:** `js/components/SmartFolders.js` + test. **Scope: S.**

---

### Checkpoint: Feeds slice
- [ ] `tests/ui/feeds.spec.js` updated POMs and passes.
- [ ] Manual: add a mock feed (`http://localhost:3001/feeds/tech-news.xml`) via `make run-test-server`, see it in sidebar with unread count.
- [ ] **Pause for human review.**

---

### Phase 6: Items vertical slice

#### Task 15: `api/items.js` + tests
**Acceptance criteria:** exports `unread`, `starred`, `unreadCount`, `starredCount`, `update(id, patch)`. Test covers each.
**Verification:** `npm run test:unit -- items`.
**Dependencies:** Task 3. **Files touched:** `js/api/items.js` + test. **Scope: S.**

---

#### Task 16: `Toolbar`
**Description:** Middle-pane header with scope title (feed name / "Unread" / "Starred"), refresh button, mark-all-read button. Buttons trigger `api/feeds.fetchNew` or `api/feeds.markRead` based on scope.
**Acceptance criteria:** title reflects current selection signal; buttons disabled when no scope; tests cover all branches.
**Verification:** `npm run test:unit -- Toolbar`.
**Dependencies:** Task 11, 15.
**Files touched:** `js/components/Toolbar.js` + test. **Scope: S.**

---

#### Task 17: `ItemList` + `ItemRow`
**Description:** `ItemList` reads scope from `state.selection`, fetches via the right `api/items.*` or `api/feeds.items`, renders `ItemRow` per row. `ItemRow` shows title, source, relative time, read/unread + star toggles. Click → set `state.currentItem` + `navigate(...)`.
**Acceptance criteria:**
- [ ] `data-testid`: `item-list`, `item-row`, `item-row-star`, `item-row-toggle-read`, `item-row-title`.
- [ ] Selecting a row marks it read (via `items.update`).
- [ ] Star toggle round-trips via `items.update`.
- [ ] Tests cover scope switching and toggle interactions.

**Verification:** `npm run test:unit -- ItemList ItemRow`.
**Dependencies:** Tasks 14, 16.
**Files touched:** `js/components/ItemList.js`, `ItemRow.js`, plus tests. **Scope: M.**

---

#### Task 18: `ItemDetail`
**Description:** Right pane reads `state.currentItem`, renders title + body via `<iframe sandbox srcdoc="...">` (browser-native isolation; backend does no sanitization — see `internal/app/feed/service.go:91`). The `sandbox` attribute is empty (no `allow-*` flags) so scripts, forms, top-navigation, and same-origin access are all blocked. RTL detection preserved (port logic from `public/js/items.js:41` in git history) and applied via the iframe's inline body style.

**Acceptance criteria:**
- [ ] Renders title + source link in the parent document; body in a sandboxed iframe.
- [ ] iframe has `sandbox=""` (no flags) and content via `srcdoc`.
- [ ] No XSS: test injects `<script>window.__pwn=1</script>` into mocked item content; after render, `window.__pwn` is `undefined`.
- [ ] Empty-state when no item selected.
- [ ] iframe height auto-sizes to content (postMessage from iframe → parent, or `ResizeObserver` once loaded).

**Verification:** `npm run test:unit -- ItemDetail`.
**Dependencies:** Task 17.
**Files touched:** `js/components/ItemDetail.js` + test. No `util/dom.js` sanitize helper needed. **Scope: S.**

---

### Checkpoint: Items slice
- [ ] `tests/ui/items.spec.js` POMs updated and passes.
- [ ] Manual smoke: signup → login → add mock feed → see items → mark read → star → switch to Starred folder.
- [ ] **Pause for human review.**

---

### Phase 7: Polish + cleanup

#### Task 19: `ErrorBoundary`, `Spinner`, `ConfirmDialog`
**Description:** Wrap `<AppPage/>` in `<ErrorBoundary/>`. Use `Spinner` during async ops. Replace any `window.confirm` with `ConfirmDialog`.
**Acceptance criteria:** thrown error inside any descendant renders fallback UI; spinner shows on `state.loading`; confirm dialog tested for keyboard + click.
**Verification:** `npm run test:unit -- ErrorBoundary Spinner ConfirmDialog`.
**Dependencies:** Task 18.
**Files touched:** three components + tests; small wiring in `AppPage.js`. **Scope: S.**

---

#### Task 20: MSW integration tests
**Description:** Set up MSW handlers covering every endpoint in `internal/infra/http/api/router.go`. Write three integration tests under `js/tests/integration/`:
- `login-flow.test.js` — signup → login → dashboard.
- `feeds-flow.test.js` — add feed → list → delete.
- `items-flow.test.js` — select feed → mark read → star → switch to Starred → see item.

**Acceptance criteria:** all three pass; MSW handlers exhaustive (any unhandled request → loud error).
**Verification:** `npm run test:unit -- integration`.
**Dependencies:** Task 19.
**Files touched:** `js/tests/msw/{handlers,server}.js`, `js/tests/integration/*.test.js`. **Scope: M.**

---

#### Task 21: Rewrite Playwright POMs to `data-testid`
**Description:** Update `tests/ui/pages/{LoginPage,SignupPage,MainPage}.js` to select via `[data-testid="..."]`. No spec file changes other than minor flow tweaks if testids replace text-based interactions. Document the convention in `CLAUDE.md`.
**Acceptance criteria:**
- [ ] All three POMs use `data-testid` exclusively.
- [ ] `make test-ui` passes against `make run-test-server` UI.
- [ ] `CLAUDE.md` adds a short "Test selectors" section.

**Verification:** `make test-ui` green.
**Dependencies:** Task 20.
**Files touched:** `tests/ui/pages/*.js`, possibly `tests/ui/utils/helpers.js`, `CLAUDE.md`. **Scope: M.**

---

#### Task 22: Delete legacy code; finalize Makefile + docs
**Description:** Remove obsolete files. Add `test-unit` Makefile target; update `pre-commit` and `test-all`.
**Acceptance criteria:**
- [ ] Deleted under `internal/web/public/`: old `css/` (legacy stylesheets), old `images/`, `app-next/`, and old `js/{jquery*,auth,items,feeds,main,login,signup,utils,stage,spin.min,moment.min}.js`.
- [ ] No `<script src="...jquery...">` references anywhere.
- [ ] `Makefile`: `test-unit:` target; `pre-commit: gomod update-mocks lint test test-unit`; `test-all: test test-unit test-ui`.
- [ ] `CLAUDE.md` reflects new paths and commands.

**Verification:**
- [ ] `make pre-commit` passes.
- [ ] `make test-all` passes.
- [ ] `git grep -i jquery` → no source matches (only changelog/comments at most).

**Dependencies:** Task 21.
**Files touched:** deletions across `internal/web/public/`; `Makefile`, `CLAUDE.md`. **Scope: M.**

---

### Checkpoint: Complete
- [ ] All acceptance criteria from `SPEC.md §2` met.
- [ ] `make build && ./lite-reader` runs the full new UI from a single binary in an empty directory.
- [ ] `make test-all` passes.
- [ ] PR description summarizes scope and links `SPEC.md` + this plan.
- [ ] Ready for review.

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| Vendored ESM versions drift from npm-resolved versions used by Vitest | Med — silent test/runtime divergence | Pin both in lockstep; add a `make verify-vendor` sha256 check (follow-up). |
| `jQuery Layout` resize behavior is hard to replicate cleanly | Med | Replace with CSS Grid + `Resizer.js` (Task 10). Accept honest downgrade: no animated collapse. |
| Embed path move breaks Dockerfile / CI pipelines | Med | Audit all references in Task 4; verify CI green before merging the move. |
| MSW + jsdom adds ~30MB to `node_modules` | Low | Acceptable cost for integration coverage. |
| Token in localStorage = XSS risk (unchanged) | Low | Out of scope per SPEC; flagged as future hardening. |
| Playwright POM rewrite churns flaky selectors | Med | Adopt `data-testid` from the first new component; rewrite POMs in one task (21) after components stabilize. |
| Long PR — single PR with 22 tasks | High | Use phase checkpoints as natural review pauses; squash within phases if reviewer prefers. Alternative: split at Phase 2/3 boundary into two PRs. |

---

## Resolved Questions

- **Two PRs** ✓ Split at Phase 2 / 3 boundary. PR-A: Phases 1+2 (foundation + embed FS, legacy UI still served). PR-B: Phases 3–7 (Preact rewrite + cleanup).
- **Vendor location** ✓ `internal/web/public/assets/vendor/`, served at `/assets/vendor/`.
- **HTML sanitization** ✓ Backend does no sanitization (raw `gofeed` Content stored to DB). Use `<iframe sandbox srcdoc>` with no `allow-*` flags in `ItemDetail`. No DOMPurify dep needed.

---

## Verification Summary (before declaring done)

- [ ] Every task has acceptance criteria — ✓
- [ ] Every task has a verification step — ✓
- [ ] Task dependencies are identified — ✓
- [ ] No task touches more than ~5 files — ✓ (largest is Task 4 directory move; tracked as M)
- [ ] Checkpoints exist between phases — ✓ (6 checkpoints)
- [ ] Human reviewed and approved — pending
