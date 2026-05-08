import { describe, it, expect, beforeAll, afterAll, beforeEach, afterEach, vi } from 'vitest';
import { server } from '../msw/server.js';
import { reset } from '../msw/handlers.js';
import { signup, login, logout } from '../../api/auth.js';
import { getToken } from '../../api/client.js';

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterAll(() => server.close());

describe('integration: login flow', () => {
  beforeEach(() => {
    reset();
    localStorage.clear();
    vi.stubGlobal('location', { assign: vi.fn(), origin: 'http://localhost:3000', pathname: '/', href: 'http://localhost:3000/' });
  });
  afterEach(() => {
    server.resetHandlers();
    vi.unstubAllGlobals();
  });

  it('signup → login persists token; logout clears it', async () => {
    await signup({ email: 'a@b.com', password: 'secret1' });
    await login({ email: 'a@b.com', password: 'secret1' });
    expect(getToken()).toMatch(/^tok-/);

    logout();
    expect(getToken()).toBeNull();
    expect(location.assign).toHaveBeenCalledWith('/login.html');
  });

  it('login with wrong password rejects without persisting', async () => {
    await signup({ email: 'a@b.com', password: 'secret1' });
    await expect(login({ email: 'a@b.com', password: 'wrong' })).rejects.toThrow(/invalid/);
    expect(getToken()).toBeNull();
  });
});
