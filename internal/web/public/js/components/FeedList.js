import { html } from '../util/html.js';
import { feeds, folders, selection } from '../state.js';
import { FeedItem } from './FeedItem.js';
import { FolderRow } from './FolderRow.js';

export function FeedList() {
  const list = feeds.value || [];
  const folderList = folders.value || [];
  const sel = selection.value;

  const byFolder = new Map();
  const rootFeeds = [];
  for (const f of list) {
    if (f.folder_id == null) {
      rootFeeds.push(f);
    } else {
      if (!byFolder.has(f.folder_id)) byFolder.set(f.folder_id, []);
      byFolder.get(f.folder_id).push(f);
    }
  }
  const sortedFolders = [...folderList].sort((a, b) => (a.position - b.position) || (a.id - b.id));

  return html`
    <ul class="folder-list" data-testid="folder-list">
      ${sortedFolders.map((f) => html`
        <${FolderRow}
          key=${f.id}
          folder=${f}
          childFeeds=${byFolder.get(f.id) || []}
        />
      `)}
    </ul>
    <ul class="feed-list root-feed-list" data-testid="root-feed-list" data-folder-id="">
      ${rootFeeds.map((f) => html`
        <${FeedItem}
          key=${f.id}
          feed=${f}
          isSelected=${sel.kind === 'feed' && sel.id === f.id}
        />
      `)}
    </ul>
  `;
}
