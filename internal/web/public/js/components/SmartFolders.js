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
  return html`
    <div class="sidebar-section" data-testid="smart-folders">
      <div
        class=${`smart-folder${sel.kind === 'unread' ? ' is-selected' : ''}`}
        data-testid="smart-folder-unread"
        onClick=${selectUnread}
      >
        <span>Unread</span>
        <span data-testid="smart-folder-unread-count">${unreadCount.value}</span>
      </div>
      <div
        class=${`smart-folder${sel.kind === 'starred' ? ' is-selected' : ''}`}
        data-testid="smart-folder-starred"
        onClick=${selectStarred}
      >
        <span>Starred</span>
        <span data-testid="smart-folder-starred-count">${starredCount.value}</span>
      </div>
    </div>
  `;
}
