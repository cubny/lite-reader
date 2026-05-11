import { request } from './client.js';

export function list() {
  return request('GET', '/feeds');
}

export function add(url) {
  return request('POST', '/feeds', { body: { url } });
}

export function remove(id) {
  return request('DELETE', `/feeds/${id}`);
}

export function fetchNew(id) {
  return request('PUT', `/feeds/${id}/fetch`);
}

export function markRead(id) {
  return request('POST', `/feeds/${id}/read`);
}

export function markUnread(id) {
  return request('POST', `/feeds/${id}/unread`);
}

export function items(id) {
  return request('GET', `/feeds/${id}/items`);
}

export function move(id, folderId) {
  const body = folderId == null ? { unset_folder: true } : { folder_id: folderId };
  return request('PATCH', `/feeds/${id}`, { body });
}

export function reorder(id, position) {
  return request('PATCH', `/feeds/${id}`, { body: { position } });
}

export function bulkMove(feedIds, folderId) {
  const body = folderId == null
    ? { feed_ids: feedIds, unset_folder: true }
    : { feed_ids: feedIds, folder_id: folderId };
  return request('POST', '/feeds-bulk-move', { body });
}
