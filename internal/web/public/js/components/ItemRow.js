import { html } from '../util/html.js';
import { currentItem } from '../state.js';
import { update as updateItem } from '../api/items.js';
import { relativeTime } from '../util/time.js';

export function ItemRow({ item, isSelected, onChanged }) {
  async function select() {
    currentItem.value = item;
    if (item.is_new) {
      try {
        await updateItem(item.id, { is_new: false, starred: !!item.starred });
        onChanged && onChanged({ ...item, is_new: false });
      } catch {
        // ignore
      }
    }
  }

  async function toggleStar(e) {
    e.stopPropagation();
    const next = !item.starred;
    try {
      await updateItem(item.id, { is_new: !!item.is_new, starred: next });
      onChanged && onChanged({ ...item, starred: next });
    } catch {
      // ignore
    }
  }

  async function toggleRead(e) {
    e.stopPropagation();
    const next = !item.is_new;
    try {
      await updateItem(item.id, { is_new: next, starred: !!item.starred });
      onChanged && onChanged({ ...item, is_new: next });
    } catch {
      // ignore
    }
  }

  const cls = [
    'item-row',
    isSelected ? 'is-selected' : '',
    item.is_new ? '' : 'is-read',
  ].filter(Boolean).join(' ');

  return html`
    <li class=${cls} data-testid="item-row" data-item-id=${item.id} onClick=${select}>
      <div class="item-row-title" data-testid="item-row-title">${item.title || '(no title)'}</div>
      <div class="item-row-meta">
        <span data-testid="item-row-time">${relativeTime(item.timestamp)}</span>
      </div>
      <div class="item-row-actions">
        <button
          type="button"
          class=${`item-row-star${item.starred ? ' is-on' : ''}`}
          data-testid="item-row-star"
          aria-label=${item.starred ? 'Unstar' : 'Star'}
          aria-pressed=${item.starred ? 'true' : 'false'}
          onClick=${toggleStar}
        >★</button>
        <button
          type="button"
          class="item-row-toggle-read"
          data-testid="item-row-toggle-read"
          aria-label=${item.is_new ? 'Mark read' : 'Mark unread'}
          onClick=${toggleRead}
        >${item.is_new ? '○' : '●'}</button>
      </div>
    </li>
  `;
}
