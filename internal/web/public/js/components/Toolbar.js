import { useState } from 'preact/hooks';
import { html } from '../util/html.js';
import { selection, feeds, items } from '../state.js';
import { fetchNew as fetchFeed, markRead, items as feedItems } from '../api/feeds.js';
import { unread as unreadItems, starred as starredItems } from '../api/items.js';
import { refreshCounts } from './SmartFolders.js';

function titleFor(sel, feedList) {
  if (!sel) return '';
  if (sel.kind === 'unread') return 'Unread';
  if (sel.kind === 'starred') return 'Starred';
  if (sel.kind === 'feed') {
    const f = (feedList || []).find((x) => x.id === sel.id);
    return f ? (f.title || f.url) : '';
  }
  return '';
}

async function reloadCurrent(sel) {
  if (!sel) return;
  if (sel.kind === 'unread') {
    items.value = (await unreadItems()) || [];
  } else if (sel.kind === 'starred') {
    items.value = (await starredItems()) || [];
  } else if (sel.kind === 'feed') {
    items.value = (await feedItems(sel.id)) || [];
  }
}

export function Toolbar() {
  const [pending, setPending] = useState(false);
  const sel = selection.value;
  const title = titleFor(sel, feeds.value);

  const isFeedScope = sel.kind === 'feed';

  async function refresh() {
    if (!isFeedScope) return;
    setPending(true);
    try {
      await fetchFeed(sel.id);
      await reloadCurrent(sel);
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  async function markAll() {
    if (!isFeedScope) return;
    setPending(true);
    try {
      await markRead(sel.id);
      await reloadCurrent(sel);
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  return html`
    <div class="toolbar" data-testid="toolbar">
      <span class="toolbar-title" data-testid="toolbar-title">${title}</span>
      <button
        type="button"
        data-testid="toolbar-refresh"
        disabled=${!isFeedScope || pending}
        onClick=${refresh}
      >Refresh</button>
      <button
        type="button"
        data-testid="toolbar-mark-read"
        disabled=${!isFeedScope || pending}
        onClick=${markAll}
      >Mark all read</button>
    </div>
  `;
}
