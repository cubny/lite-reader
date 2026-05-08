import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import * as items from './items.js';

function mockFetchOnce({ status = 200, body = null } = {}) {
  globalThis.fetch = vi.fn().mockResolvedValue({
    status,
    ok: status >= 200 && status < 300,
    json: vi.fn().mockResolvedValue(body),
    text: vi.fn().mockResolvedValue(body == null ? '' : JSON.stringify(body)),
  });
}

describe('api/items', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.stubGlobal('location', { assign: vi.fn(), pathname: '/', href: '/' });
  });
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it('unread → GET /items/unread', async () => {
    mockFetchOnce({ body: [] });
    await items.unread();
    expect(fetch.mock.calls[0][0]).toBe('/items/unread');
  });

  it('starred → GET /items/starred', async () => {
    mockFetchOnce({ body: [] });
    await items.starred();
    expect(fetch.mock.calls[0][0]).toBe('/items/starred');
  });

  it('unreadCount → GET /items/unread/count', async () => {
    mockFetchOnce({ body: { count: 12 } });
    const res = await items.unreadCount();
    expect(res).toEqual({ count: 12 });
    expect(fetch.mock.calls[0][0]).toBe('/items/unread/count');
  });

  it('starredCount → GET /items/starred/count', async () => {
    mockFetchOnce({ body: { count: 3 } });
    await items.starredCount();
    expect(fetch.mock.calls[0][0]).toBe('/items/starred/count');
  });

  it('update(id, patch) → PUT /items/:id with body', async () => {
    mockFetchOnce({ status: 200 });
    await items.update(42, { is_new: false, starred: true });
    const [path, init] = fetch.mock.calls[0];
    expect(path).toBe('/items/42');
    expect(init.method).toBe('PUT');
    expect(JSON.parse(init.body)).toEqual({ is_new: false, starred: true });
  });
});
