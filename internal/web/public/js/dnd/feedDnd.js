import Sortable from 'sortablejs';
import { feeds, folders, folderCollapsed, persistCollapsed } from '../state.js';
import { move as moveFeed, reorder as reorderFeed, list as listFeeds } from '../api/feeds.js';
import { reorder as reorderFolder, list as listFolders } from '../api/folders.js';

const POSITION_GAP = 1024;
const HOVER_EXPAND_DELAY = 400;

function parseFolderId(el) {
  if (!el) return null;
  const v = el.getAttribute('data-folder-id');
  if (v == null || v === '') return null;
  const n = Number(v);
  return Number.isFinite(n) ? n : null;
}

function readFeedId(el) {
  const v = el && el.getAttribute('data-feed-id');
  if (!v) return null;
  const n = Number(v);
  return Number.isFinite(n) ? n : null;
}

async function refreshAll() {
  const [freshFeeds, freshFolders] = await Promise.all([listFeeds(), listFolders()]);
  feeds.value = freshFeeds || [];
  folders.value = freshFolders || [];
}

// Rewrite positions of every feed in `listEl` based on its current DOM order.
// Sparse integers (POSITION_GAP apart) leave room for future insertions
// without renumbering. Sent in parallel; failures still fall back to a
// listFeeds refresh.
async function persistListPositions(listEl) {
  const items = Array.from(listEl.querySelectorAll(':scope > li[data-feed-id]'));
  await Promise.all(items.map((li, idx) => {
    const feedId = readFeedId(li);
    if (feedId == null) return Promise.resolve();
    return reorderFeed(feedId, (idx + 1) * POSITION_GAP);
  }));
}

// Auto-expand a folder while a feed is being dragged over its header. Mimics
// the Finder/Explorer pattern so the user never has to expand a folder first
// just to drop into it.
function attachHoverExpand(root) {
  const headers = root.querySelectorAll('.folder-wrap');
  let timer = null;

  function clearTimer() {
    if (timer) {
      clearTimeout(timer);
      timer = null;
    }
  }

  const handlers = [];
  for (const wrap of headers) {
    const folderId = parseFolderId(wrap);
    if (folderId == null) continue;
    const onEnter = () => {
      clearTimer();
      const map = folderCollapsed.value;
      if (!map[folderId]) return; // already expanded
      timer = setTimeout(() => {
        const next = { ...folderCollapsed.value };
        delete next[folderId];
        folderCollapsed.value = next;
        persistCollapsed(next);
      }, HOVER_EXPAND_DELAY);
    };
    const onLeave = () => clearTimer();
    wrap.addEventListener('dragenter', onEnter);
    wrap.addEventListener('dragleave', onLeave);
    handlers.push(() => {
      wrap.removeEventListener('dragenter', onEnter);
      wrap.removeEventListener('dragleave', onLeave);
    });
  }
  return () => {
    clearTimer();
    for (const off of handlers) off();
  };
}

export function attachFeedDnd(root) {
  const instances = [];
  const lists = root.querySelectorAll('ul[data-folder-id]');
  for (const list of lists) {
    instances.push(Sortable.create(list, {
      group: 'feeds',
      animation: 150,
      delay: 200,
      delayOnTouchOnly: true,
      touchStartThreshold: 5,
      emptyInsertThreshold: 24,
      ghostClass: 'feed-drag-ghost',
      dragClass: 'feed-drag',
      filter: '.folder-chevron, .folder-menu-btn',
      preventOnFilter: false,
      onEnd: async (evt) => {
        const feedId = readFeedId(evt.item);
        if (feedId == null) return;
        const fromFolderId = parseFolderId(evt.from);
        const toFolderId = parseFolderId(evt.to);
        const moved = fromFolderId !== toFolderId;
        const reordered = evt.newIndex !== evt.oldIndex;
        if (!moved && !reordered) return;
        try {
          if (moved) {
            await moveFeed(feedId, toFolderId);
          }
          // Renumber positions of the destination list so the new order
          // sticks across reloads. When a feed crossed lists, also renumber
          // the source list to keep things tidy.
          await persistListPositions(evt.to);
          if (moved && evt.from !== evt.to) {
            await persistListPositions(evt.from);
          }
        } finally {
          await refreshAll();
        }
      },
    }));
  }

  const folderListEl = root.querySelector('.folder-list');
  if (folderListEl) {
    instances.push(Sortable.create(folderListEl, {
      group: 'folders',
      animation: 150,
      delay: 200,
      delayOnTouchOnly: true,
      handle: '.folder-icon',
      ghostClass: 'feed-drag-ghost',
      dragClass: 'feed-drag',
      onEnd: async (evt) => {
        const li = evt.item;
        const folderId = li && Number(li.getAttribute('data-folder-id'));
        if (!Number.isFinite(folderId)) return;
        const newPosition = (evt.newIndex || 0) * POSITION_GAP + POSITION_GAP;
        try {
          await reorderFolder(folderId, newPosition);
        } finally {
          const fresh = await listFolders();
          folders.value = fresh || [];
        }
      },
    }));
  }

  const detachHover = attachHoverExpand(root);

  return () => {
    detachHover();
    for (const inst of instances) {
      try { inst.destroy(); } catch { /* ignore */ }
    }
  };
}
