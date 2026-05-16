import { useEffect, useRef } from 'preact/hooks';
import { html } from '../util/html.js';
import { feeds, folders } from '../state.js';
import { list as listFeeds } from '../api/feeds.js';
import { list as listFolders } from '../api/folders.js';
import { FeedList } from './FeedList.js';
import { AddFeedForm } from './AddFeedForm.js';
import { AddFolderButton } from './AddFolderButton.js';
import { SmartFolders } from './SmartFolders.js';
import { attachFeedDnd } from '../dnd/feedDnd.js';

export function Sidebar() {
  const rootRef = useRef(null);

  useEffect(() => {
    Promise.all([
      listFeeds().then((list) => { feeds.value = list || []; }).catch(() => {}),
      listFolders().then((list) => { folders.value = list || []; }).catch(() => {}),
    ]);
  }, []);

  // Attach SortableJS to the rendered lists. Re-attach when feeds or folders
  // change so newly-rendered <ul>s become drop targets.
  useEffect(() => {
    if (!rootRef.current) return undefined;
    let detach = () => {};
    try {
      detach = attachFeedDnd(rootRef.current);
    } catch {
      // SortableJS is best-effort; failure shouldn't break the app
    }
    return () => {
      try { detach && detach(); } catch { /* ignore */ }
    };
  }, [feeds.value, folders.value]);

  return html`
    <div data-testid="sidebar" ref=${rootRef}>
      <a class="lr-brand" href="#/" data-testid="sidebar-brand">
        <span class="lr-brand-name">Lite Reader</span>
      </a>
      <div class="sidebar-scroll">
        <ul class="smart-folders">
          <${SmartFolders} />
        </ul>
        <div class="lr-divider"></div>
        <${FeedList} />
        <${AddFolderButton} />
      </div>
      <div class="sidebar-footer">
        <div id="toolbar">
          <${AddFeedForm} />
        </div>
      </div>
    </div>
  `;
}
