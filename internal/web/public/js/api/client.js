const TOKEN_KEY = 'token';
const LOGIN_PATH = '/login.html';

export class ApiError extends Error {
  constructor(status, message) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

export async function request(method, path, opts = {}) {
  const { body, auth = true } = opts;
  const headers = { Accept: 'application/json' };

  if (body !== undefined) {
    headers['Content-Type'] = 'application/json';
  }

  if (auth) {
    const token = getToken();
    if (token) headers.Authorization = `Bearer ${token}`;
  }

  const init = { method, headers };
  if (body !== undefined) init.body = JSON.stringify(body);

  const response = await fetch(path, init);

  if (response.status === 401 && auth) {
    clearToken();
    if (location.pathname !== LOGIN_PATH) {
      location.assign(LOGIN_PATH);
    }
  }

  if (!response.ok) {
    const message = await readErrorMessage(response);
    throw new ApiError(response.status, message);
  }

  return readJsonOrNull(response);
}

async function readJsonOrNull(response) {
  if (response.status === 204) return null;
  const text = await response.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}

async function readErrorMessage(response) {
  try {
    const data = await response.json();
    if (data && typeof data.message === 'string') return data.message;
  } catch {
    // fall through
  }
  return response.statusText || `HTTP ${response.status}`;
}
