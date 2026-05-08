# SPEC — Lite Reader Frontend Modernization

## 1. Objective

Replace the legacy jQuery 1.x + jQuery Layout frontend under `public/` with a modern, build-free Preact-based static client served from an embedded filesystem in the Go binary. The backend remains a pure JSON API.

**Target users (in priority order):**
1. **Non-technical self-hosters** — must download one binary and run it. No install steps, no asset folders.
2. **Developers contributing to the project** — clean component boundaries, fast unit tests, obvious file layout.
3. **Coding agents** (Claude Code, Copilot, etc.) — predictable file structure, explicit imports, `data-testid` selectors, no hidden DSLs.

**Out of scope:**
- Auth hardening (httpOnly cookies, CSRF, refresh tokens) — keep current localStorage + Bearer model.
- UI/UX redesign or new features — pure migration, visual parity with current jQuery UI.
- Mobile / responsive improvements — desktop-only, matching today.
- Internationalization or formal a11y audit — preserve current semantics, no new framework.

## 2. Acceptance Criteria

The change is "done" when **all** of the following hold:

- [ ] `make build` produces a single static binary; `./lite-reader` serves the full UI on :3000 with no `public/` folder present.
- [ ] All existing Playwright specs pass: `auth.spec.js`, `feeds.spec.js`, `items.spec.js`, `smoke.spec.js` (POMs rewritten to use `data-testid`).
- [ ] Every file under `internal/web/public/js/components/` and `internal/web/public/js/api/` has at least one Vitest unit test.
- [ ] At least one MSW-backed integration test covers each main flow: login → feeds load → select feed → mark read → star → logout.
- [ ] Visual + interaction parity with current jQuery UI: same 3-pane layout, same actions, same keyboard behavior.
- [ ] `make test-all` passes (Go unit + Vitest + Playwright).
- [ ] `make lint` passes; `make pre-commit` runs Go + JS test suites.
- [ ] Legacy code deleted: `public/app-next/`, `public/js/jquery*`, `public/js/{auth,items,feeds,main,login,signup,utils,stage,spin.min,moment.min}.js`, old `public/css/`, `public/images/`.

## 3. Commands

| Command | Purpose |
|---|---|
| `make run` | Run via `go run`; serves embedded UI on :3000. |
| `make run-test-server` | Run app + mock feed server (port 3001) against `data/test-agg.db`. |
| `make build` | Produce static `lite-reader` binary with embedded UI. |
| `PUBLIC_DIR=./internal/web/public make run` | Dev escape hatch — serve UI from disk, no rebuild needed on JS changes. |
| `make test` | Go unit tests (race, coverage). |
| `make test-unit` | **NEW** — Vitest unit + integration (MSW) tests. |
| `make test-ui` | Playwright e2e (boots testserver). |
| `make test-all` | `test` + `test-unit` + `test-ui`. |
| `make pre-commit` | `gomod` + `update-mocks` + `lint` + `test` + `test-unit`. (Run `make test-ui` separately for UI-touching PRs.) |
| `npm run test:unit:watch` | Vitest watch mode for local TDD. |

## 4. Project Structure

```
internal/web/
  web.go                                # //go:embed all:public, exposes fs.FS
  public/
    index.html                          # SPA shell
    login.html, signup.html             # static auth pages
    favicon.ico
    css/
      app.css                           # 3-pane layout, theme
      auth.css
    js/
      main.js                           # entry: mounts <AppPage/>
      login.js                          # entry: mounts <LoginPage/>
      signup.js                         # entry: mounts <SignupPage/>
      router.js                         # ~30-line hash router
      state.js                          # global signals
      util/
        html.js                         # export const html = htm.bind(h)
        time.js                         # relative-time formatter
        dom.js
      api/
        client.js                       # fetch wrapper, token, 401 handler
        client.test.js
        auth.js                         # login, signup, logout
        auth.test.js
        feeds.js                        # list, add, remove, fetchNew, markRead, ...
        feeds.test.js
        items.js                        # unread, starred, counts, update
        items.test.js
      pages/
        LoginPage.js   + .test.js
        SignupPage.js  + .test.js
        AppPage.js     + .test.js
      components/
        Sidebar.js, FeedList.js, FeedItem.js, AddFeedForm.js,
        SmartFolders.js, ItemList.js, ItemRow.js, ItemDetail.js,
        Toolbar.js, Resizer.js, Spinner.js, ConfirmDialog.js,
        ErrorBoundary.js
        # each .js has a co-located .test.js
      tests/
        integration/
          login-flow.test.js
          feeds-flow.test.js
          items-flow.test.js
        msw/
          handlers.js                   # MSW request handlers
          server.js                     # node setup
    vendor/
      preact@10.24.3/{preact,hooks,signals}.module.js
      htm@3.1.1/htm.module.js

tests/ui/                                # existing Playwright suite (POMs rewritten)
  pages/{LoginPage,SignupPage,MainPage}.js
  *.spec.js

vitest.config.js                         # NEW — at repo root
package.json                             # NEW deps + scripts
```

**Naming conventions:**
- One component or module per file. File name = exported symbol (`FeedItem.js` exports `FeedItem`).
- Tests co-located: `Foo.js` ↔ `Foo.test.js`.
- Integration tests live under `internal/web/public/js/tests/integration/`.

## 5. Code Style

**JavaScript:**
- Plain ES2022 modules. No TypeScript. No JSX (use `htm` template literals).
- Always import via bare specifiers resolved by the import map: `import { h } from 'preact'`, `import { html } from '../util/html.js'`.
- Use `import { signal, computed, effect } from '@preact/signals'` for shared state. Use `useState`/`useEffect` only for component-local state.
- Components are named-exported functions returning `html\`...\``. No default exports.
- API modules export pure async functions; never read DOM or `window` directly except in `client.js` (token + redirect).
- Every interactive element has `data-testid="..."`. Test selectors **never** rely on CSS classes, tag structure, or visible text.
- No `console.log` in committed code. Errors flow through `ErrorBoundary` or `state.toast`.

**Go:**
- New `internal/web/web.go` — single file, `//go:embed all:public`, exposes `FS() fs.FS`.
- `internal/infra/http/api/router.go` — `New(...)` accepts `staticFS fs.FS`; uses `http.FileServer(http.FS(staticFS))`.
- Dev escape hatch in `internal/app.go`: if `os.Getenv("PUBLIC_DIR")` is set, use `os.DirFS(env)` instead of embedded FS.

**CSS:**
- Plain CSS, no preprocessor. Custom properties (CSS variables) for layout dimensions and theme colors.
- Pane widths controlled by `--sidebar-w`, `--list-w` updated by `Resizer.js` and persisted to localStorage.

## 6. Testing Strategy

**Three layers, all required:**

### Unit (Vitest + @testing-library/preact + jsdom)
- Co-located `*.test.js`. Run via `make test-unit`.
- Cover every `components/*.js`, every `pages/*.js`, every `api/*.js`.
- Vitest config aliases bare specifiers (`preact`, `htm`, `@preact/signals`) to `node_modules` versions matching the vendored ones.
- Assertions via `@testing-library/preact` + `@testing-library/jest-dom`. Selectors via `getByTestId` only.
- `api/client.test.js` must cover: token attached, 401 clears token + redirects, non-2xx throws `ApiError`.

### Integration (Vitest + MSW)
- Lives under `internal/web/public/js/tests/integration/`.
- MSW handlers in `tests/msw/handlers.js` mock every endpoint listed in `internal/infra/http/api/router.go`.
- Flows covered: login → feeds load → select feed → mark item read → star → logout. Add feed. Delete feed. Switch between Unread / Starred / per-feed scopes.

### E2E (Playwright, existing suite)
- `tests/ui/` specs unchanged in intent.
- POMs (`tests/ui/pages/{LoginPage,SignupPage,MainPage}.js`) rewritten to use `data-testid` selectors.
- Run against `make run-test-server` with `TEST_DB_PATH=data/test-agg.db`.
- Mock feeds only — never real internet URLs (per current `CLAUDE.md`).

**CI gate:** `make test-all` must pass. `make pre-commit` includes `test-unit` but not `test-ui` (UI tests run separately due to time cost).

## 7. Boundaries

### NEVER
- **Add a frontend build step.** No Vite, webpack, Rollup, esbuild, Babel, or transpiler. Vendored ESM ships byte-for-byte to the browser.
- **Introduce React, Vue, Svelte, Next.js, or any non-Preact framework.** Stack is locked: Preact + htm + @preact/signals.
- **Render HTML server-side from Go for application pages.** No HTMX, no `html/template` for app routes. Backend is JSON-only. (Static `.html` files served from embedded FS are fine — they are static assets, not server-rendered templates.)
- **Use CSS classes, tag structure, or visible text as test selectors.** Always `data-testid`.
- **Edit generated mock files** under `internal/mocks/` by hand (existing rule, reaffirmed).

### ALWAYS
- **Keep `internal/web/public/vendor/` versions pinned** and matched to `package.json` dev versions used by Vitest.
- **Tag every interactive element with `data-testid`** in new components.
- **Co-locate tests** with the file under test (`Foo.test.js` next to `Foo.js`).
- **Embed assets via `//go:embed`** so the binary is self-contained.
- **Honor `PUBLIC_DIR` env var** for dev hot-reload.

### ASK FIRST
- **Adding any new frontend runtime dependency** (even tiny). Each dep must be vendored, pinned, and justified — and confirmed with the user before adoption.
- **Changing the Preact/htm/signals versions.**
- **Moving or renaming the `internal/web/public/` directory.**
- **Diverging from visual/UX parity** (any redesign, new buttons, layout changes).
- **Extending scope** into auth hardening, mobile, i18n, or a11y work.

---

**Plan reference:** `/Users/cubny/.claude/plans/ideate-how-we-can-eager-pudding.md` — contains step-by-step migration order and per-file change list.
