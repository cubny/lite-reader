import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';

import { getToken, setToken } from './client.js';
import { ApiError } from './client.js';
import { login, signup, logout } from './auth.js';

function mockFetchOnce({ status = 200, body = null, ok } = {}) {
  const response = {
    status,
    ok: ok ?? (status >= 200 && status < 300),
    json: vi.fn().mockResolvedValue(body),
    text: vi.fn().mockResolvedValue(body == null ? '' : JSON.stringify(body)),
  };
  globalThis.fetch = vi.fn().mockResolvedValue(response);
  return response;
}

describe('api/auth', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.stubGlobal('location', {
      assign: vi.fn(),
      pathname: '/',
      href: 'http://localhost/',
    });
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  describe('login', () => {
    it('POSTs /login with email+password and auth=false', async () => {
      mockFetchOnce({
        status: 200,
        body: { access_token: 'tok-1', refresh_token: 'r', expires_in: 3600 },
      });
      await login({ email: 'a@b.com', password: 'pw' });
      const [path, init] = fetch.mock.calls[0];
      expect(path).toBe('/login');
      expect(init.method).toBe('POST');
      expect(init.headers.Authorization).toBeUndefined();
      expect(JSON.parse(init.body)).toEqual({ email: 'a@b.com', password: 'pw' });
    });

    it('persists access_token via setToken on success', async () => {
      mockFetchOnce({
        status: 200,
        body: { access_token: 'new-jwt', refresh_token: 'r', expires_in: 60 },
      });
      const result = await login({ email: 'a@b.com', password: 'pw' });
      expect(getToken()).toBe('new-jwt');
      expect(result).toMatchObject({ access_token: 'new-jwt' });
    });

    it('throws ApiError and does NOT persist token on failure', async () => {
      mockFetchOnce({ status: 400, body: { message: 'invalid email or password' } });
      await expect(
        login({ email: 'a@b.com', password: 'wrong' }),
      ).rejects.toMatchObject({
        name: 'ApiError',
        status: 400,
        message: 'invalid email or password',
      });
      expect(getToken()).toBeNull();
    });

    it('does not redirect on 401 (auth=false on login endpoint)', async () => {
      mockFetchOnce({ status: 401, body: { message: 'unauthorized' } });
      await expect(login({ email: 'a@b.com', password: 'pw' })).rejects.toBeInstanceOf(ApiError);
      expect(location.assign).not.toHaveBeenCalled();
    });
  });

  describe('signup', () => {
    it('POSTs /signup with email+password and auth=false', async () => {
      mockFetchOnce({ status: 201, body: null });
      await signup({ email: 'a@b.com', password: 'pw' });
      const [path, init] = fetch.mock.calls[0];
      expect(path).toBe('/signup');
      expect(init.method).toBe('POST');
      expect(init.headers.Authorization).toBeUndefined();
      expect(JSON.parse(init.body)).toEqual({ email: 'a@b.com', password: 'pw' });
    });

    it('resolves with null on 201', async () => {
      mockFetchOnce({ status: 201, body: null });
      const result = await signup({ email: 'a@b.com', password: 'pw' });
      expect(result).toBeNull();
    });

    it('throws ApiError on validation failure', async () => {
      mockFetchOnce({ status: 400, body: { message: 'email already in use' } });
      await expect(
        signup({ email: 'taken@x.com', password: 'pw' }),
      ).rejects.toMatchObject({
        name: 'ApiError',
        status: 400,
        message: 'email already in use',
      });
    });

    it('does not persist a token (signup does not return one)', async () => {
      mockFetchOnce({ status: 201, body: null });
      await signup({ email: 'a@b.com', password: 'pw' });
      expect(getToken()).toBeNull();
    });
  });

  describe('logout', () => {
    it('clears token and redirects to /login.html', () => {
      setToken('tok');
      logout();
      expect(getToken()).toBeNull();
      expect(location.assign).toHaveBeenCalledWith('/login.html');
    });

    it('clears token even when none was set, and redirects', () => {
      logout();
      expect(getToken()).toBeNull();
      expect(location.assign).toHaveBeenCalledWith('/login.html');
    });
  });
});
