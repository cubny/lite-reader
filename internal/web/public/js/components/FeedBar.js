import { useState } from 'preact/hooks';
import { html } from '../util/html.js';
import { selection, feeds, items } from '../state.js';
import {
  fetchNew as fetchFeed,
  markRead,
  markUnread,
  remove as removeFeed,
  items as feedItems,
  list as listFeeds,
} from '../api/feeds.js';
import { unread as unreadItems, starred as starredItems } from '../api/items.js';
import { logout as authLogout } from '../api/auth.js';
import { refreshCounts } from './SmartFolders.js';
import { ConfirmDialog } from './ConfirmDialog.js';

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

export function FeedBar() {
  const [pending, setPending] = useState(false);
  const [confirming, setConfirming] = useState(false);
  const sel = selection.value;
  const title = titleFor(sel, feeds.value);
  const isFeedScope = sel.kind === 'feed';
  const isSmartScope = sel.kind === 'unread' || sel.kind === 'starred';
  const barClass = isFeedScope ? 'has-feed' : (isSmartScope ? 'has-smart' : '');

  async function refresh() {
    if (!isFeedScope || pending) return;
    setPending(true);
    try {
      await fetchFeed(sel.id);
      await reloadCurrent(sel);
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  async function readAll() {
    if (!isFeedScope || pending) return;
    setPending(true);
    try {
      await markRead(sel.id);
      await reloadCurrent(sel);
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  async function unreadAll() {
    if (!isFeedScope || pending) return;
    setPending(true);
    try {
      await markUnread(sel.id);
      await reloadCurrent(sel);
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  function askRemove() {
    if (!isFeedScope) return;
    setConfirming(true);
  }

  async function confirmRemove() {
    setConfirming(false);
    if (!isFeedScope) return;
    setPending(true);
    try {
      await removeFeed(sel.id);
      const fresh = await listFeeds();
      feeds.value = fresh || [];
      items.value = [];
      selection.value = { kind: 'unread' };
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  return html`
    <div id="feedbar" class=${barClass}>
      <div id="title" data-testid="toolbar-title">${title}</div>
      <div id="actions">
        <div class="action update" data-testid="toolbar-refresh" onClick=${refresh}>
          <i class="icon-repeat"></i> Update
        </div>
        <div id="mark-read-all" class="action markread" data-testid="toolbar-mark-read" onClick=${readAll}>
          <i class="icon-circle-blank"></i> Read All
        </div>
        <div id="mark-unread-all" class="action markread" data-testid="toolbar-mark-unread" onClick=${unreadAll}>
          <i class="icon-circle"></i> Unread All
        </div>
        <div class="action logout" id="logout" data-testid="toolbar-logout" onClick=${authLogout}>
          <i class="icon-signout"></i> Logout
        </div>
        <div class="action remove" data-testid="toolbar-remove" onClick=${askRemove}>
          <i class="icon-trash"></i> Remove
        </div>
      </div>
      ${confirming && html`
        <${ConfirmDialog}
          message=${`Delete "${title}"?`}
          onConfirm=${confirmRemove}
          onCancel=${() => setConfirming(false)}
        />
      `}
    </div>
  `;
}
