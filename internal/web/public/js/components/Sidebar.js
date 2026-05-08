import { useEffect, useState } from 'preact/hooks';
import { html } from '../util/html.js';
import { feeds } from '../state.js';
import { list as listFeeds, remove as removeFeed } from '../api/feeds.js';
import { logout } from '../api/auth.js';
import { FeedList } from './FeedList.js';
import { AddFeedForm } from './AddFeedForm.js';
import { SmartFolders } from './SmartFolders.js';
import { ConfirmDialog } from './ConfirmDialog.js';

export function Sidebar() {
  const [confirming, setConfirming] = useState(null);

  useEffect(() => {
    listFeeds().then((list) => { feeds.value = list || []; }).catch(() => {});
  }, []);

  function askDelete(feed) {
    setConfirming(feed);
  }

  async function confirmDelete() {
    const f = confirming;
    setConfirming(null);
    if (!f) return;
    try {
      await removeFeed(f.id);
      const fresh = await listFeeds();
      feeds.value = fresh || [];
    } catch {
      // ignore
    }
  }

  return html`
    <div class="sidebar" data-testid="sidebar">
      <button
        type="button"
        class="sidebar-logout"
        data-testid="logout-button"
        onClick=${logout}
      >Log out</button>
      <${SmartFolders} />
      <${AddFeedForm} />
      <div class="sidebar-section">
        <h2>Feeds</h2>
        <${FeedList} onDelete=${askDelete} />
      </div>
      ${confirming && html`
        <${ConfirmDialog}
          message=${`Delete "${confirming.title || confirming.url}"?`}
          onConfirm=${confirmDelete}
          onCancel=${() => setConfirming(null)}
        />
      `}
    </div>
  `;
}
