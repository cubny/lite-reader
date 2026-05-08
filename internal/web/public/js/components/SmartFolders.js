import { useEffect } from 'preact/hooks';
import { html } from '../util/html.js';
import { unreadCount, starredCount, selection } from '../state.js';
import { unreadCount as fetchUnread, starredCount as fetchStarred } from '../api/items.js';
import { navigate } from '../router.js';

export async function refreshCounts() {
  try {
    const u = await fetchUnread();
    if (u && typeof u.count === 'number') unreadCount.value = u.count;
    const s = await fetchStarred();
    if (s && typeof s.count === 'number') starredCount.value = s.count;
  } catch {
    // silent — toolbar/handlers surface errors
  }
}

export function SmartFolders() {
  useEffect(() => { refreshCounts(); }, []);

  function selectUnread() {
    selection.value = { kind: 'unread' };
    navigate('#/');
  }
  function selectStarred() {
    selection.value = { kind: 'starred' };
    navigate('#/starred');
  }

  const sel = selection.value;
  const u = unreadCount.value;
  const s = starredCount.value;
  return html`
    <li
      id="unread"
      class=${`feed${sel.kind === 'unread' ? ' selected' : ''}`}
      data-testid="smart-folder-unread"
      onClick=${selectUnread}
    >
      <div class="count">${u > 0 ? html`<span data-testid="smart-folder-unread-count">${u}</span>` : html`<span data-testid="smart-folder-unread-count" style="display:none">${u}</span>`}</div>
      <i class="icon-circle"></i>
      <div class="feedtitle">Unread</div>
    </li>
    <li
      id="starred"
      class=${`feed${sel.kind === 'starred' ? ' selected' : ''}`}
      data-testid="smart-folder-starred"
      onClick=${selectStarred}
    >
      <div class="count">${s > 0 ? html`<span data-testid="smart-folder-starred-count">${s}</span>` : html`<span data-testid="smart-folder-starred-count" style="display:none">${s}</span>`}</div>
      <i class="icon-star"></i>
      <div class="feedtitle">Starred</div>
    </li>
  `;
}
