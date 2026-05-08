import { html } from '../util/html.js';
import { selection } from '../state.js';
import { navigate } from '../router.js';

export function FeedItem({ feed, onDelete, isSelected }) {
  function select() {
    selection.value = { kind: 'feed', id: feed.id };
    navigate(`#/feed/${feed.id}`);
  }

  function onDeleteClick(e) {
    e.stopPropagation();
    onDelete(feed);
  }

  return html`
    <li
      class=${`feed-item${isSelected ? ' is-selected' : ''}`}
      data-testid="feed-item"
      data-feed-id=${feed.id}
      onClick=${select}
    >
      <span class="feed-item-title" data-testid="feed-item-title">${feed.title || feed.url}</span>
      <span class="feed-item-unread-count" data-testid="feed-item-unread-count">${feed.unread_count || 0}</span>
      <button
        type="button"
        class="feed-item-delete"
        data-testid="feed-item-delete"
        aria-label="Delete feed"
        onClick=${onDeleteClick}
      >×</button>
    </li>
  `;
}
