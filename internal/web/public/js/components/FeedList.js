import { html } from '../util/html.js';
import { feeds, selection } from '../state.js';
import { FeedItem } from './FeedItem.js';

export function FeedList({ onDelete }) {
  const list = feeds.value || [];
  const sel = selection.value;
  return html`
    <ul class="feed-list" data-testid="feed-list">
      ${list.map((f) => html`
        <${FeedItem}
          key=${f.id}
          feed=${f}
          onDelete=${onDelete}
          isSelected=${sel.kind === 'feed' && sel.id === f.id}
        />
      `)}
    </ul>
  `;
}
