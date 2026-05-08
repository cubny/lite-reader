import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';

import {
  ApiError,
  clearToken,
  getToken,
  request,
  setToken,
} from './client.js';

const TOKEN_KEY = 'token';

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

describe('api/client', () => {
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

  describe('token storage', () => {
    it('getToken returns null when nothing stored', () => {
      expect(getToken()).toBeNull();
    });

    it('setToken persists to localStorage and getToken reads it', () => {
      setToken('abc.def.ghi');
      expect(localStorage.getItem(TOKEN_KEY)).toBe('abc.def.ghi');
      expect(getToken()).toBe('abc.def.ghi');
    });

    it('clearToken removes the value', () => {
      setToken('xyz');
      clearToken();
      expect(localStorage.getItem(TOKEN_KEY)).toBeNull();
      expect(getToken()).toBeNull();
    });
  });

  describe('request: success', () => {
    it('returns parsed JSON on 2xx', async () => {
      mockFetchOnce({ status: 200, body: { ok: true, value: 42 } });
      const result = await request('GET', '/feeds');
      expect(result).toEqual({ ok: true, value: 42 });
    });

    it('returns null when response has no body (204)', async () => {
      const response = mockFetchOnce({ status: 204, body: null });
      response.text.mockResolvedValue('');
      const result = await request('DELETE', '/feeds/1');
      expect(result).toBeNull();
    });

    it('sends method, path, and JSON body', async () => {
      mockFetchOnce({ status: 200, body: {} });
      await request('POST', '/feeds', { body: { url: 'http://x' } });
      const [path, init] = fetch.mock.calls[0];
      expect(path).toBe('/feeds');
      expect(init.method).toBe('POST');
      expect(init.headers['Content-Type']).toBe('application/json');
      expect(init.body).toBe(JSON.stringify({ url: 'http://x' }));
    });
  });

  describe('request: auth header', () => {
    it('attaches Authorization: Bearer when token present and auth defaults true', async () => {
      setToken('tok-123');
      mockFetchOnce({ status: 200, body: {} });
      await request('GET', '/feeds');
      const init = fetch.mock.calls[0][1];
      expect(init.headers.Authorization).toBe('Bearer tok-123');
    });

    it('omits Authorization when no token stored', async () => {
      mockFetchOnce({ status: 200, body: {} });
      await request('GET', '/feeds');
      const init = fetch.mock.calls[0][1];
      expect(init.headers.Authorization).toBeUndefined();
    });

    it('omits Authorization when auth=false even if token present', async () => {
      setToken('tok-123');
      mockFetchOnce({ status: 200, body: {} });
      await request('POST', '/login', { body: {}, auth: false });
      const init = fetch.mock.calls[0][1];
      expect(init.headers.Authorization).toBeUndefined();
    });
  });

  describe('request: errors', () => {
    it('throws ApiError with status and parsed message on 4xx', async () => {
      mockFetchOnce({ status: 400, body: { message: 'bad input' } });
      await expect(request('POST', '/feeds', { body: {} })).rejects.toMatchObject({
        name: 'ApiError',
        status: 400,
        message: 'bad input',
      });
    });

    it('throws ApiError with status when body has no message', async () => {
      const response = mockFetchOnce({ status: 500, body: null });
      response.json.mockRejectedValue(new Error('not json'));
      response.text.mockResolvedValue('Internal Server Error');
      await expect(request('GET', '/feeds')).rejects.toMatchObject({
        name: 'ApiError',
        status: 500,
      });
    });

    it('on 401 with auth=true: clears token and redirects to /login.html', async () => {
      setToken('expired');
      mockFetchOnce({ status: 401, body: { message: 'unauthorized' } });
      await expect(request('GET', '/feeds')).rejects.toBeInstanceOf(ApiError);
      expect(getToken()).toBeNull();
      expect(location.assign).toHaveBeenCalledWith('/login.html');
    });

    it('on 401 with auth=false: does NOT clear token or redirect', async () => {
      setToken('still-valid');
      mockFetchOnce({ status: 401, body: { message: 'wrong creds' } });
      await expect(
        request('POST', '/login', { body: {}, auth: false }),
      ).rejects.toBeInstanceOf(ApiError);
      expect(getToken()).toBe('still-valid');
      expect(location.assign).not.toHaveBeenCalled();
    });

    it('does not redirect if already on /login.html', async () => {
      vi.stubGlobal('location', {
        assign: vi.fn(),
        pathname: '/login.html',
        href: 'http://localhost/login.html',
      });
      setToken('expired');
      mockFetchOnce({ status: 401, body: {} });
      await expect(request('GET', '/feeds')).rejects.toBeInstanceOf(ApiError);
      expect(location.assign).not.toHaveBeenCalled();
    });
  });

  describe('ApiError', () => {
    it('is an Error subclass with name "ApiError"', () => {
      const e = new ApiError(404, 'not found');
      expect(e).toBeInstanceOf(Error);
      expect(e.name).toBe('ApiError');
      expect(e.status).toBe(404);
      expect(e.message).toBe('not found');
    });
  });
});
