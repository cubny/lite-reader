import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, waitFor, cleanup } from '@testing-library/preact';
import userEvent from '@testing-library/user-event';

import { html } from '../util/html.js';

const loginMock = vi.fn();
vi.mock('../api/auth.js', () => ({
  login: (...args) => loginMock(...args),
}));

const { LoginPage } = await import('./LoginPage.js');

describe('LoginPage', () => {
  beforeEach(() => {
    loginMock.mockReset();
    localStorage.clear();
    sessionStorage.clear();
    vi.stubGlobal('location', {
      assign: vi.fn(),
      pathname: '/login.html',
      href: 'http://localhost/login.html',
    });
  });

  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('renders email, password, submit and error region with data-testid', () => {
    render(html`<${LoginPage} />`);
    expect(screen.getByTestId('login-email')).toBeInTheDocument();
    expect(screen.getByTestId('login-password')).toBeInTheDocument();
    expect(screen.getByTestId('login-submit')).toBeInTheDocument();
    expect(screen.getByTestId('login-error')).toBeInTheDocument();
  });

  it('error region is empty initially', () => {
    render(html`<${LoginPage} />`);
    expect(screen.getByTestId('login-error').textContent).toBe('');
  });

  it('on success: calls auth.login and redirects to /', async () => {
    loginMock.mockResolvedValue({ access_token: 'tok', refresh_token: 'r', expires_in: 60 });
    const user = userEvent.setup();
    render(html`<${LoginPage} />`);

    await user.type(screen.getByTestId('login-email'), 'a@b.com');
    await user.type(screen.getByTestId('login-password'), 'secret');
    await user.click(screen.getByTestId('login-submit'));

    await waitFor(() => {
      expect(loginMock).toHaveBeenCalledWith({ email: 'a@b.com', password: 'secret' });
    });
    await waitFor(() => {
      expect(location.assign).toHaveBeenCalledWith('/');
    });
  });

  it('on failure: shows inline error and does NOT redirect', async () => {
    loginMock.mockRejectedValue(
      Object.assign(new Error('invalid email or password'), {
        name: 'ApiError',
        status: 400,
      }),
    );
    const user = userEvent.setup();
    render(html`<${LoginPage} />`);

    await user.type(screen.getByTestId('login-email'), 'a@b.com');
    await user.type(screen.getByTestId('login-password'), 'wrong');
    await user.click(screen.getByTestId('login-submit'));

    await waitFor(() => {
      expect(screen.getByTestId('login-error').textContent).toContain('invalid email or password');
    });
    expect(location.assign).not.toHaveBeenCalled();
  });

  it('disables submit while pending', async () => {
    let resolve;
    loginMock.mockReturnValue(new Promise((r) => { resolve = r; }));
    const user = userEvent.setup();
    render(html`<${LoginPage} />`);

    await user.type(screen.getByTestId('login-email'), 'a@b.com');
    await user.type(screen.getByTestId('login-password'), 'pw');
    await user.click(screen.getByTestId('login-submit'));

    await waitFor(() => {
      expect(screen.getByTestId('login-submit').disabled).toBe(true);
    });

    resolve({ access_token: 't' });
  });

  it('shows signup-success banner when sessionStorage flag is set, and clears flag', () => {
    sessionStorage.setItem('signupSuccess', 'true');
    render(html`<${LoginPage} />`);
    expect(screen.getByTestId('signup-success')).toBeInTheDocument();
    expect(sessionStorage.getItem('signupSuccess')).toBeNull();
  });

  it('hides signup-success banner when flag is not set', () => {
    render(html`<${LoginPage} />`);
    expect(screen.queryByTestId('signup-success')).toBeNull();
  });
});
