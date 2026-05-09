import { clearToken, request, setToken } from './client.js';

export async function login({ email, password }) {
  const response = await request('POST', '/login', {
    body: { email, password },
    auth: false,
  });
  if (response && typeof response.access_token === 'string') {
    setToken(response.access_token);
  }
  return response;
}

export async function signup({ email, password }) {
  return request('POST', '/signup', {
    body: { email, password },
    auth: false,
  });
}

export async function setupAdmin({ email, password, confirm_password, allow_signup }) {
  return request('POST', '/setup', {
    body: { email, password, confirm_password, allow_signup },
    auth: false,
  });
}

export async function checkNeedsSetup() {
  return request('GET', '/setup/status', { auth: false });
}

export async function getSettings() {
  return request('GET', '/settings');
}

export async function updateSettings({ allow_signup }) {
  return request('PUT', '/settings', {
    body: { allow_signup },
  });
}

export function logout() {
  clearToken();
  location.assign('/login.html');
}
