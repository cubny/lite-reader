import { useMemo, useState } from 'preact/hooks';
import { html } from '../util/html.js';
import { update as updateItem, scrape as scrapeItem } from '../api/items.js';
import { relativeTime } from '../util/time.js';
import { detectDir } from '../util/dom.js';

export function ItemRow({ item, isSelected, onToggle, onChanged }) {
  const body = item.full_content || item.desc || '';
  const dir = useMemo(
    () => detectDir((item.title || '') + ' ' + body),
    [item.id, item.title, body],
  );
  const hasFull = !!item.full_content;
  const [scraping, setScraping] = useState(false);
  const [scrapeError, setScrapeError] = useState('');

  async function loadFullArticle(e) {
    e.preventDefault();
    e.stopPropagation();
    if (scraping || hasFull) return;
    setScraping(true);
    setScrapeError('');
    try {
      const updated = await scrapeItem(item.id);
      onChanged && onChanged({ ...item, ...updated });
    } catch (err) {
      setScrapeError('Could not load full article');
    } finally {
      setScraping(false);
    }
  }

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
        <span class="title-toggles">
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
        </span>
        <span class="title-text">${item.title || '(no title)'}</span>
        <span class="title-meta">
          <span class="title-actions">
            ${item.link && html`
              <a
                href=${item.link}
                target="_blank"
                rel="noopener noreferrer"
                class="item-action item-link"
                data-testid="item-row-link"
                onClick=${(e) => e.stopPropagation()}
              ><i class="icon-external-link"></i><span>source</span></a>
            `}
            ${!hasFull && item.link && html`
              <a
                href="#"
                class="item-action item-load-full"
                data-testid="item-row-load-full"
                aria-label="Read full article inline"
                aria-busy=${scraping ? 'true' : 'false'}
                onClick=${loadFullArticle}
              >${scraping
                ? html`<i class="icon-spinner icon-spin"></i><span>loading</span>`
                : html`<i class="icon-file-text"></i><span>read here</span>`}</a>
            `}
          </span>
          <span class="timestamp" data-testid="item-row-time">${ts}</span>
        </span>
      </div>
      ${scrapeError && html`<div class="scrape-error" data-testid="item-row-scrape-error">${scrapeError}</div>`}
      <div class="desc" dir=${dir} data-testid="item-row-body" dangerouslySetInnerHTML=${{ __html: body }}></div>
    </li>
  `;
}
