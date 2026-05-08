import { useEffect, useState } from 'preact/hooks';
import { html } from '../util/html.js';
import { selection, items, currentItem } from '../state.js';
import { unread as unreadItems, starred as starredItems } from '../api/items.js';
import { items as feedItems } from '../api/feeds.js';
import { ItemRow } from './ItemRow.js';

async function loadFor(sel) {
  if (!sel) return [];
  if (sel.kind === 'unread') return (await unreadItems()) || [];
  if (sel.kind === 'starred') return (await starredItems()) || [];
  if (sel.kind === 'feed') return (await feedItems(sel.id)) || [];
  return [];
}

export function ItemList() {
  const sel = selection.value;
  const [expanded, setExpanded] = useState(null);

  useEffect(() => {
    let cancelled = false;
    loadFor(sel).then((list) => {
      if (!cancelled) {
        items.value = list;
        setExpanded(null);
        currentItem.value = null;
      }
    }).catch(() => {});
    return () => { cancelled = true; };
  }, [sel.kind, sel.id]);

  function patch(updated) {
    items.value = items.value.map((it) => (it.id === updated.id ? updated : it));
    if (currentItem.value && currentItem.value.id === updated.id) {
      currentItem.value = updated;
    }
  }

  function toggleExpand(id) {
    setExpanded((prev) => (prev === id ? null : id));
    const next = items.value.find((it) => it.id === id);
    currentItem.value = next || null;
  }

  if (!items.value || items.value.length === 0) {
    return html`<ul id="items" data-testid="item-list"></ul>`;
  }

  return html`
    <ul id="items" data-testid="item-list">
      ${items.value.map((it) => html`
        <${ItemRow}
          key=${it.id}
          item=${it}
          isSelected=${expanded === it.id}
          onToggle=${() => toggleExpand(it.id)}
          onChanged=${patch}
        />
      `)}
    </ul>
  `;
}
