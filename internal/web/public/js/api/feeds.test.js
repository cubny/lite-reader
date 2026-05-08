import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import * as feeds from './feeds.js';

function mockFetchOnce({ status = 200, body = null } = {}) {
  globalThis.fetch = vi.fn().mockResolvedValue({
    status,
    ok: status >= 200 && status < 300,
    json: vi.fn().mockResolvedValue(body),
    text: vi.fn().mockResolvedValue(body == null ? '' : JSON.stringify(body)),
  });
}

describe('api/feeds', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.stubGlobal('location', { assign: vi.fn(), pathname: '/', href: '/' });
  });
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it('list → GET /feeds', async () => {
    mockFetchOnce({ body: [{ id: 1 }] });
    const res = await feeds.list();
    expect(res).toEqual([{ id: 1 }]);
    expect(fetch.mock.calls[0][0]).toBe('/feeds');
    expect(fetch.mock.calls[0][1].method).toBe('GET');
  });

  it('add → POST /feeds with {url}', async () => {
    mockFetchOnce({ body: { id: 1, url: 'http://x' } });
    await feeds.add('http://x');
    const [path, init] = fetch.mock.calls[0];
    expect(path).toBe('/feeds');
    expect(init.method).toBe('POST');
    expect(JSON.parse(init.body)).toEqual({ url: 'http://x' });
  });

  it('remove → DELETE /feeds/:id', async () => {
    mockFetchOnce({ status: 204 });
    await feeds.remove(7);
    expect(fetch.mock.calls[0][0]).toBe('/feeds/7');
    expect(fetch.mock.calls[0][1].method).toBe('DELETE');
  });

  it('fetchNew → PUT /feeds/:id/fetch', async () => {
    mockFetchOnce({ body: [] });
    await feeds.fetchNew(7);
    expect(fetch.mock.calls[0][0]).toBe('/feeds/7/fetch');
    expect(fetch.mock.calls[0][1].method).toBe('PUT');
  });

  it('markRead → POST /feeds/:id/read', async () => {
    mockFetchOnce({ status: 200 });
    await feeds.markRead(7);
    expect(fetch.mock.calls[0][0]).toBe('/feeds/7/read');
    expect(fetch.mock.calls[0][1].method).toBe('POST');
  });

  it('markUnread → POST /feeds/:id/unread', async () => {
    mockFetchOnce({ status: 200 });
    await feeds.markUnread(7);
    expect(fetch.mock.calls[0][0]).toBe('/feeds/7/unread');
    expect(fetch.mock.calls[0][1].method).toBe('POST');
  });

  it('items(id) → GET /feeds/:id/items', async () => {
    mockFetchOnce({ body: [{ id: 100 }] });
    const res = await feeds.items(7);
    expect(res).toEqual([{ id: 100 }]);
    expect(fetch.mock.calls[0][0]).toBe('/feeds/7/items');
  });
});
