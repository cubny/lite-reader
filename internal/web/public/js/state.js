import { signal } from '@preact/signals';

export const token = signal(null);
export const feeds = signal([]);
export const selection = signal({ kind: 'unread' });
export const items = signal([]);
export const currentItem = signal(null);
export const unreadCount = signal(0);
export const starredCount = signal(0);
export const toast = signal(null);
export const loading = signal(false);
