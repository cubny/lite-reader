import { signal } from '@preact/signals';

export const currentRoute = signal(parse(getHash()));

function getHash() {
  return (typeof location !== 'undefined' && location.hash) || '#/';
}

export function parse(hash) {
  const path = (hash || '').replace(/^#/, '') || '/';
  if (path === '/' || path === '') return { name: 'unread', params: {} };
  if (path === '/starred') return { name: 'starred', params: {} };

  const item = path.match(/^\/feed\/([^/]+)\/item\/([^/]+)$/);
  if (item) return { name: 'feed-item', params: { id: item[1], itemId: item[2] } };

  const feed = path.match(/^\/feed\/([^/]+)$/);
  if (feed) return { name: 'feed', params: { id: feed[1] } };

  const folderItem = path.match(/^\/folder\/([^/]+)\/item\/([^/]+)$/);
  if (folderItem) return { name: 'folder-item', params: { id: folderItem[1], itemId: folderItem[2] } };

  const folder = path.match(/^\/folder\/([^/]+)$/);
  if (folder) return { name: 'folder', params: { id: folder[1] } };

  return { name: 'unknown', params: { path } };
}

export function navigate(path) {
  if (typeof location !== 'undefined') {
    location.hash = path.startsWith('#') ? path : `#${path}`;
  }
}

if (typeof window !== 'undefined') {
  window.addEventListener('hashchange', () => {
    currentRoute.value = parse(getHash());
  });
}
