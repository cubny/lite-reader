# Self-Hosted RSS Reader Comparison

Top self-hosted RSS readers on GitHub benchmarked against **lite-reader**. Star counts pulled live from the GitHub API on 2026-05-14.

## At a glance

| Reader | Stars | Language | Database | Binary size / footprint |
|---|---:|---|---|---|
| **FreshRSS** | 15.0k | PHP | SQLite / MySQL / PostgreSQL | Heavy (PHP runtime) |
| **Miniflux** | 9.2k | Go | PostgreSQL only | Single binary, minimal |
| **Stringer** | 4.1k | Ruby | PostgreSQL | Ruby runtime |
| **yarr** | 3.8k | Go | SQLite | Single binary, tiny |
| **CommaFeed** | 3.5k | Java | H2 / MySQL / PostgreSQL | JVM |
| **Tiny Tiny RSS** | n/a¹ | PHP | MySQL / PostgreSQL | PHP runtime |
| **lite-reader** | — | Go | SQLite (pure Go) | Single static binary, CGO-free |

¹ Tiny Tiny RSS development lives on the author's own git server; the GitHub mirror is unofficial.

## Feature matrix

| Feature | lite-reader | FreshRSS | Miniflux | yarr | CommaFeed | TT-RSS | Stringer |
|---|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| **RSS 2.0 / Atom** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **JSON Feed** | ❌ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **Multi-user** | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ |
| **Admin / signup control** | ✅ | ✅ | ✅ | n/a | ✅ | ✅ | ❌ |
| **Setup wizard** | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Folders / categories** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ (tags only) |
| **Drag-and-drop org.** | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ |
| **Starred / favorites** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Full-text search** | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Article scraper / full content** | ✅ on-demand | ✅ | ✅ | ❌ | ❌ | ✅ (plugin) | ❌ |
| **Filtering rules / scoring** | ❌ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **Themes / dark mode** | ⚠️ basic | ✅ many | ✅ light/dark | ✅ light/dark | ✅ light/dark | ✅ many | ⚠️ basic |
| **OPML import / export** | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Mobile-friendly UI** | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ⚠️ |
| **PWA / mobile API** | ❌ | ✅ (Fever, Google Reader API) | ✅ (Fever, Google Reader API) | ✅ (Fever) | ✅ (Fever, Google Reader API) | ✅ (own + Fever) | ❌ |
| **WebSub / PubSubHubbub** | ❌ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **Plugins / extensions** | ❌ | ✅ | ⚠️ integrations | ❌ | ❌ | ✅ | ❌ |
| **Webhooks / integrations** | ❌ | ⚠️ via ext | ✅ (Pocket, Wallabag, Telegram…) | ❌ | ❌ | ⚠️ via plugin | ❌ |
| **RTL / Unicode text dir** | ✅ auto | ✅ | ✅ | ⚠️ | ⚠️ | ✅ | ⚠️ |
| **Background scheduler** | ✅ 1h | ✅ cron | ✅ | ✅ | ✅ | ✅ cron | ✅ |
| **Static frontend** | ✅ Preact | ❌ server-rendered | ❌ server-rendered | ✅ Vue | ❌ React+API | ❌ | ❌ |
| **Single static binary** | ✅ | ❌ | ✅ | ✅ | ❌ (jar) | ❌ | ❌ |
| **No external DB required** | ✅ | ✅ | ❌ | ✅ | ✅ | ❌ | ❌ |
| **License** | — | AGPL-3.0 | Apache-2.0 | MIT | Apache-2.0 | GPL-3.0 | MIT |

## Positioning

- **lite-reader's sweet spot** — a single CGO-free Go binary with SQLite, multi-user + folders + drag-and-drop, modern Preact SPA, RTL-aware. Closest in spirit to **yarr** (Go + SQLite + tiny) but adds multi-user and admin controls. Closest in shape to **Miniflux** but without the PostgreSQL requirement.
- **Main gaps vs. the leaders** — no full-text search, no OPML import/export, no Google Reader / Fever API (so no mobile-app ecosystem), no article scraping, no filtering rules, no plugin system, no WebSub.
- **Where it leads** — true zero-dependency deploy (no DB server, no runtime), first-class folders with drag-and-drop, opinionated setup wizard, automatic RTL detection.

If parity work were prioritized: **OPML import/export → full-text search → Google Reader/Fever API → article scraper**. The first two are table stakes; the API unlocks the entire third-party mobile-client ecosystem (Reeder, NetNewsWire, FocusReader, etc.).

## Sources

- [FreshRSS on GitHub](https://github.com/FreshRSS/FreshRSS)
- [Miniflux on GitHub](https://github.com/miniflux/v2)
- [yarr on GitHub](https://github.com/nkanaev/yarr)
- [CommaFeed on GitHub](https://github.com/Athou/commafeed)
- [Stringer on GitHub](https://github.com/stringer-rss/stringer)
- [awesome-selfhosted — Feed Readers](https://awesome-selfhosted.net/tags/feed-readers.html)
- [selfh.st — RSS reader alternatives](https://selfh.st/alternatives/rss-readers/)
- [FreshRSS vs Miniflux 2026 — OSSAlt](https://ossalt.com/guides/freshrss-vs-miniflux-2026)
- [Open-source RSS reader comparison gist](https://gist.github.com/kevinmichaelchen/9d40fde5b8408fc0417f187359e07001/)
