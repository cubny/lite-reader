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
