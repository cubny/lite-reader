import { describe, it, expect, beforeAll, afterAll, beforeEach, afterEach, vi } from 'vitest';
import { server } from '../msw/server.js';
import { reset } from '../msw/handlers.js';
import * as feeds from '../../api/feeds.js';
import { signup, login } from '../../api/auth.js';

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterAll(() => server.close());

describe('integration: feeds flow', () => {
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

  it('add → list → delete', async () => {
    expect(await feeds.list()).toEqual([]);
    const added = await feeds.add('http://example.com/r.xml');
    expect(added.id).toBeGreaterThan(0);

    const after = await feeds.list();
    expect(after).toHaveLength(1);

    await feeds.remove(added.id);
    expect(await feeds.list()).toHaveLength(0);
  });
});
