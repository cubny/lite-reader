import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, waitFor, cleanup } from '@testing-library/preact';
import userEvent from '@testing-library/user-event';

import { html } from '../util/html.js';

const signupMock = vi.fn();
vi.mock('../api/auth.js', () => ({
  signup: (...args) => signupMock(...args),
}));

const { SignupPage } = await import('./SignupPage.js');

async function fillForm(user, { email = 'a@b.com', password = 'secret1', confirm = 'secret1' } = {}) {
  if (email) await user.type(screen.getByTestId('signup-email'), email);
  if (password) await user.type(screen.getByTestId('signup-password'), password);
  if (confirm) await user.type(screen.getByTestId('signup-confirm-password'), confirm);
}

describe('SignupPage', () => {
  beforeEach(() => {
    signupMock.mockReset();
    sessionStorage.clear();
    vi.stubGlobal('location', {
      assign: vi.fn(),
      pathname: '/signup.html',
      href: 'http://localhost/signup.html',
    });
  });

  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('renders email, password, confirm-password, submit, error region', () => {
    render(html`<${SignupPage} />`);
    expect(screen.getByTestId('signup-email')).toBeInTheDocument();
    expect(screen.getByTestId('signup-password')).toBeInTheDocument();
    expect(screen.getByTestId('signup-confirm-password')).toBeInTheDocument();
    expect(screen.getByTestId('signup-submit')).toBeInTheDocument();
    expect(screen.getByTestId('signup-error')).toBeInTheDocument();
  });

  it('rejects invalid email format without calling api', async () => {
    const user = userEvent.setup();
    render(html`<${SignupPage} />`);
    await fillForm(user, { email: 'not-an-email' });
    await user.click(screen.getByTestId('signup-submit'));
    await waitFor(() => {
      expect(screen.getByTestId('signup-error').textContent).toMatch(/valid email/i);
    });
    expect(signupMock).not.toHaveBeenCalled();
  });

  it('rejects passwords shorter than 6 characters', async () => {
    const user = userEvent.setup();
    render(html`<${SignupPage} />`);
    await fillForm(user, { password: 'abc', confirm: 'abc' });
    await user.click(screen.getByTestId('signup-submit'));
    await waitFor(() => {
      expect(screen.getByTestId('signup-error').textContent).toMatch(/at least 6/i);
    });
    expect(signupMock).not.toHaveBeenCalled();
  });

  it('rejects when password does not match confirm', async () => {
    const user = userEvent.setup();
    render(html`<${SignupPage} />`);
    await fillForm(user, { password: 'secret1', confirm: 'secret2' });
    await user.click(screen.getByTestId('signup-submit'));
    await waitFor(() => {
      expect(screen.getByTestId('signup-error').textContent).toMatch(/match/i);
    });
    expect(signupMock).not.toHaveBeenCalled();
  });

  it('on success: calls signup, sets sessionStorage flag, redirects to /login.html', async () => {
    signupMock.mockResolvedValue(null);
    const user = userEvent.setup();
    render(html`<${SignupPage} />`);
    await fillForm(user);
    await user.click(screen.getByTestId('signup-submit'));

    await waitFor(() => {
      expect(signupMock).toHaveBeenCalledWith({ email: 'a@b.com', password: 'secret1' });
    });
    await waitFor(() => {
      expect(sessionStorage.getItem('signupSuccess')).toBe('true');
      expect(location.assign).toHaveBeenCalledWith('/login.html');
    });
  });

  it('on api error: shows inline error, no redirect', async () => {
    signupMock.mockRejectedValue(
      Object.assign(new Error('email already in use'), { name: 'ApiError', status: 400 }),
    );
    const user = userEvent.setup();
    render(html`<${SignupPage} />`);
    await fillForm(user);
    await user.click(screen.getByTestId('signup-submit'));
    await waitFor(() => {
      expect(screen.getByTestId('signup-error').textContent).toContain('email already in use');
    });
    expect(location.assign).not.toHaveBeenCalled();
  });

  it('disables submit while pending', async () => {
    let resolve;
    signupMock.mockReturnValue(new Promise((r) => { resolve = r; }));
    const user = userEvent.setup();
    render(html`<${SignupPage} />`);
    await fillForm(user);
    await user.click(screen.getByTestId('signup-submit'));
    await waitFor(() => {
      expect(screen.getByTestId('signup-submit').disabled).toBe(true);
    });
    resolve(null);
  });
});
