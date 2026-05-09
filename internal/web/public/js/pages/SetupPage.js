import { useState, useEffect } from 'preact/hooks';

import { html } from '../util/html.js';
import { setupAdmin, checkNeedsSetup } from '../api/auth.js';

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export function SetupPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [allowSignup, setAllowSignup] = useState(false);
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkNeedsSetup().then((data) => {
      if (!data.needs_setup) {
        location.assign('/login.html');
      } else {
        setLoading(false);
      }
    }).catch(() => {
      setLoading(false);
    });
  }, []);

  function validate() {
    if (!EMAIL_RE.test(email.trim())) return 'Please enter a valid email address';
    if (!password || password.length < 8) return 'Password must be at least 8 characters long';
    if (password !== confirm) return 'Passwords do not match';
    return '';
  }

  async function onSubmit(e) {
    e.preventDefault();
    const v = validate();
    if (v) {
      setError(v);
      return;
    }
    setError('');
    setPending(true);
    try {
      await setupAdmin({
        email: email.trim(),
        password,
        confirm_password: confirm,
        allow_signup: allowSignup,
      });
      sessionStorage.setItem('setupSuccess', 'true');
      location.assign('/login.html');
    } catch (err) {
      setError(err.message || 'Setup failed');
      setPending(false);
    }
  }

  if (loading) {
    return html`<div class="login-container"><p style="text-align:center;">Loading...</p></div>`;
  }

  return html`
    <div class="login-container">
      <h1>Welcome to Lite Reader</h1>
      <p style="text-align: center; color: #6e6a60; margin-bottom: 20px; font-size: 14px;">
        Create your admin account to get started.
      </p>
      <form class="login-form" novalidate onSubmit=${onSubmit} data-testid="setup-form">
        <div class="form-group">
          <label for="email">Email:</label>
          <input
            type="email"
            id="email"
            name="email"
            class="form-control"
            data-testid="setup-email"
            value=${email}
            onInput=${(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div class="form-group">
          <label for="password">Password:</label>
          <input
            type="password"
            id="password"
            name="password"
            class="form-control"
            data-testid="setup-password"
            value=${password}
            onInput=${(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <div class="form-group">
          <label for="confirm-password">Confirm Password:</label>
          <input
            type="password"
            id="confirm-password"
            name="confirm-password"
            class="form-control"
            data-testid="setup-confirm-password"
            value=${confirm}
            onInput=${(e) => setConfirm(e.target.value)}
            required
          />
        </div>
        <div class="form-group" style="flex-direction: row; align-items: center; gap: 8px; margin-top: 4px;">
          <input
            type="checkbox"
            id="allow-signup"
            data-testid="setup-allow-signup"
            checked=${allowSignup}
            onChange=${(e) => setAllowSignup(e.target.checked)}
            style="width: auto;"
          />
          <label for="allow-signup" style="font-size: 14px; color: #2a2a28;">
            Allow other people to create accounts
          </label>
        </div>
        <button
          type="submit"
          class="btn-primary"
          data-testid="setup-submit"
          disabled=${pending}
          style="margin-top: 10px;"
        >Create Account</button>
        <div class="error-message" data-testid="setup-error" role="alert">${error}</div>
      </form>
    </div>
  `;
}
