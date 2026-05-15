import { request } from './client.js';

export function unread() {
  return request('GET', '/items/unread');
}

export function starred() {
  return request('GET', '/items/starred');
}

export function unreadCount() {
  return request('GET', '/items/unread/count');
}

export function starredCount() {
  return request('GET', '/items/starred/count');
}

export function update(id, patch) {
  return request('PUT', `/items/${id}`, { body: patch });
}

export function scrape(id) {
  return request('POST', `/items/${id}/scrape`);
}
