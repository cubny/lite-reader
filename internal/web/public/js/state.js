import { signal } from '@preact/signals';

export const token = signal(null);
export const feeds = signal([]);
export const folders = signal([]);
export const folderCollapsed = signal(loadCollapsed());
export const selection = signal({ kind: 'unread' });
export const items = signal([]);
export const currentItem = signal(null);
export const unreadCount = signal(0);
export const starredCount = signal(0);
export const toast = signal(null);
export const loading = signal(false);

const COLLAPSED_KEY = 'folderCollapsed';

function loadCollapsed() {
  try {
    if (typeof localStorage === 'undefined') return {};
    const raw = localStorage.getItem(COLLAPSED_KEY);
    return raw ? JSON.parse(raw) : {};
  } catch {
    return {};
  }
}

export function persistCollapsed(map) {
  try {
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(COLLAPSED_KEY, JSON.stringify(map));
    }
  } catch {
    // ignore
  }
}
