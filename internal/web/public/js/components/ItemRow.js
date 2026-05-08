import { useMemo } from 'preact/hooks';
import { html } from '../util/html.js';
import { update as updateItem } from '../api/items.js';
import { relativeTime } from '../util/time.js';
import { detectDir } from '../util/dom.js';

export function ItemRow({ item, isSelected, onToggle, onChanged }) {
  const dir = useMemo(
    () => detectDir((item.title || '') + ' ' + (item.desc || '')),
    [item.id, item.title, item.desc],
  );

  async function onClickTitle() {
    onToggle && onToggle();
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
    e.preventDefault();
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
    e.preventDefault();
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
    item.is_new ? 'new' : '',
    isSelected ? 'selected' : '',
    dir === 'rtl' ? 'rtl' : '',
  ].filter(Boolean).join(' ');

  const starIcon = item.starred ? 'icon-star' : 'icon-star-empty';
  const readIcon = item.is_new ? 'icon-circle' : 'icon-circle-blank';
  const ts = item.timestamp ? relativeTime(item.timestamp) : '';

  return html`
    <li
      id=${`item-${item.id}`}
      class=${cls}
      data-testid="item-row"
      data-item-id=${item.id}
    >
      <div class="title" onClick=${onClickTitle} data-testid="item-row-title">
        <a
          name="starred"
          class="item-action item-star"
          data-testid="item-row-star"
          aria-label=${item.starred ? 'Unstar' : 'Star'}
          aria-pressed=${item.starred ? 'true' : 'false'}
          onClick=${toggleStar}
        ><i class=${starIcon}></i></a>
        <a
          name="read"
          class="item-action item-read"
          data-testid="item-row-toggle-read"
          aria-label=${item.is_new ? 'Mark read' : 'Mark unread'}
          onClick=${toggleRead}
        ><i class=${readIcon}></i></a>
        <span>${item.title || '(no title)'}</span>
        ${item.link && html`
          <a
            href=${item.link}
            target="_blank"
            rel="noopener noreferrer"
            class="item-action item-link"
            data-testid="item-row-link"
            onClick=${(e) => e.stopPropagation()}
          ><i class="icon-external-link"></i> link</a>
        `}
        <span class="timestamp" data-testid="item-row-time">${ts}</span>
      </div>
      <div class="dir">${dir}</div>
      <div class="desc" dangerouslySetInnerHTML=${{ __html: item.desc || '' }}></div>
    </li>
  `;
}
