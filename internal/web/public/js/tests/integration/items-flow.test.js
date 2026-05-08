import { describe, it, expect, beforeAll, afterAll, beforeEach, afterEach, vi } from 'vitest';
import { server } from '../msw/server.js';
import { reset } from '../msw/handlers.js';
import * as feeds from '../../api/feeds.js';
import * as items from '../../api/items.js';
import { signup, login } from '../../api/auth.js';

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterAll(() => server.close());

describe('integration: items flow', () => {
  beforeEach(async () => {
    reset();
    localStorage.clear();
    vi.stubGlobal('location', { assign: vi.fn(), origin: 'http://localhost:3000', pathname: '/', href: 'http://localhost:3000/' });
    await signup({ email: 'a@b.com', password: 'secret1' });
    await login({ email: 'a@b.com', password: 'secret1' });
  });
  afterEach(() => {
    server.resetHandlers();
    vi.unstubAllGlobals();
  });

  it('select feed → mark read → star → switch to Starred', async () => {
    const f = await feeds.add('http://example.com/r.xml');
    const list = await feeds.items(f.id);
    expect(list.length).toBeGreaterThan(0);
    const first = list[0];

    // Mark read
    await items.update(first.id, { is_new: false, starred: false });
    const unreadList = await items.unread();
    expect(unreadList.find((i) => i.id === first.id)).toBeUndefined();

    // Star
    await items.update(first.id, { is_new: false, starred: true });
    const starredList = await items.starred();
    expect(starredList.find((i) => i.id === first.id)).toBeDefined();

    // Counts reflect state
    const u = await items.unreadCount();
    const s = await items.starredCount();
    expect(s.count).toBe(1);
    expect(u.count).toBe(list.length - 1);
  });
});
