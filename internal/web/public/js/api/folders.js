import { request } from './client.js';

export function list() {
  return request('GET', '/folders');
}

export function create(name) {
  return request('POST', '/folders', { body: { name } });
}

export function rename(id, name) {
  return request('PATCH', `/folders/${id}`, { body: { name } });
}

export function reorder(id, position) {
  return request('PATCH', `/folders/${id}`, { body: { position } });
}

export function remove(id) {
  return request('DELETE', `/folders/${id}`);
}

export function items(id) {
  return request('GET', `/folders/${id}/items`);
}

export function markRead(id) {
  return request('POST', `/folders/${id}/read`);
}
