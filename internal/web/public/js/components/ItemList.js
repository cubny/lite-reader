import { useEffect, useMemo, useState } from 'preact/hooks';
import { html } from '../util/html.js';
import { selection, items, currentItem, feeds } from '../state.js';
import { unread as unreadItems, starred as starredItems } from '../api/items.js';
import { items as feedItems } from '../api/feeds.js';
import { items as folderItems } from '../api/folders.js';
import { ItemRow } from './ItemRow.js';

async function loadFor(sel) {
  if (!sel) return [];
  if (sel.kind === 'unread') return (await unreadItems()) || [];
  if (sel.kind === 'starred') return (await starredItems()) || [];
  if (sel.kind === 'feed') return (await feedItems(sel.id)) || [];
  if (sel.kind === 'folder') return (await folderItems(sel.id)) || [];
  return [];
}

export function ItemList() {
  const sel = selection.value;
  const [expanded, setExpanded] = useState(null);
  const feedById = useMemo(() => {
    const map = new Map();
    for (const f of feeds.value || []) map.set(f.id, f.title || f.url);
    return map;
  }, [feeds.value]);

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

  // Legacy: space advances to the next item when the current one is fully
  // visible (i.e. you've scrolled it). main.js bound this on document.
  useEffect(() => {
    function onKey(e) {
      if (e.code !== 'Space' && e.key !== ' ') return;
      const tag = (e.target && e.target.tagName) || '';
      if (tag === 'INPUT' || tag === 'TEXTAREA' || (e.target && e.target.isContentEditable)) return;
      if (expanded == null) return;
      const list = items.value || [];
      const idx = list.findIndex((it) => it.id === expanded);
      if (idx < 0 || idx >= list.length - 1) return;

      const currentLi = document.getElementById(`item-${expanded}`);
      const scroller = document.querySelector('#content > .ui-layout-center');
      if (currentLi && scroller) {
        const liRect = currentLi.getBoundingClientRect();
        const scRect = scroller.getBoundingClientRect();
        // Only advance once the bottom of the current item is on-screen.
        if (liRect.bottom > scRect.bottom) return;
      }

      e.preventDefault();
      const nextItem = list[idx + 1];
      toggleExpand(nextItem.id);
      requestAnimationFrame(() => {
        const el = document.getElementById(`item-${nextItem.id}`);
        if (el && el.scrollIntoView) el.scrollIntoView({ behavior: 'smooth', block: 'start' });
      });
    }
    document.addEventListener('keydown', onKey);
    return () => document.removeEventListener('keydown', onKey);
  }, [expanded]);

  if (!items.value || items.value.length === 0) {
    return html`
      <ul id="items" data-testid="item-list"></ul>
      <div class="lr-empty-callout" data-testid="item-list-empty">
        <div class="lr-callout-eyebrow">— Quiet inbox</div>
        <h2 class="lr-callout-title">No items yet</h2>
        <p class="lr-callout-sub">
          Add a feed in the sidebar to start collecting briefs. New items will appear here, numbered, in reverse chronological order.
        </p>
      </div>
    `;
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
          feedTitle=${feedById.get(it.feed_id) || ''}
        />
      `)}
    </ul>
  `;
}
