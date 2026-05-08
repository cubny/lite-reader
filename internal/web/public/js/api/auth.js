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

export function logout() {
  clearToken();
  location.assign('/login.html');
}
