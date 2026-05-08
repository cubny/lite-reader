import { html } from '../util/html.js';
import { feeds, selection } from '../state.js';
import { FeedItem } from './FeedItem.js';

export function FeedList() {
  const list = feeds.value || [];
  const sel = selection.value;
  return list.map((f) => html`
    <${FeedItem}
      key=${f.id}
      feed=${f}
      isSelected=${sel.kind === 'feed' && sel.id === f.id}
    />
  `);
}
