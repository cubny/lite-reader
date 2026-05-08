import { html } from '../util/html.js';
import { selection } from '../state.js';
import { navigate } from '../router.js';

export function FeedItem({ feed, isSelected }) {
  function select() {
    selection.value = { kind: 'feed', id: feed.id };
    navigate(`#/feed/${feed.id}`);
  }

  const unread = feed.unread_count || 0;

  return html`
    <li
      class=${`feed${isSelected ? ' selected' : ''}`}
      id=${feed.id}
      data-testid="feed-item"
      data-feed-id=${feed.id}
      onClick=${select}
    >
      <div class="count">${unread > 0 ? html`<span data-testid="feed-item-unread-count">${unread}</span>` : html`<span data-testid="feed-item-unread-count" style="display:none">${unread}</span>`}</div>
      <i class="icon-rss"></i>
      <div class="feedtitle" data-testid="feed-item-title">${feed.title || feed.url}</div>
    </li>
  `;
}
