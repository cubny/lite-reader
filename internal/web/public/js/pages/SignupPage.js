import { useState, useEffect } from 'preact/hooks';

import { html } from '../util/html.js';
import { signup, checkNeedsSetup } from '../api/auth.js';

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export function SignupPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkNeedsSetup().then((data) => {
      if (data.needs_setup) {
        location.assign('/setup.html');
        return;
      }
      if (!data.allow_signup) {
        location.assign('/login.html');
        return;
      }
      setLoading(false);
    }).catch(() => {
      setLoading(false);
    });
  }, []);

  function validate() {
    if (!EMAIL_RE.test(email.trim())) return 'Please enter a valid email address';
    if (!password || password.length < 6) return 'Password must be at least 6 characters long';
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
      await signup({ email: email.trim(), password });
      sessionStorage.setItem('signupSuccess', 'true');
      location.assign('/login.html');
    } catch (err) {
      setError(err.message || 'Error creating account.');
      setPending(false);
    }
  }

  if (loading) {
    return html`<div class="login-container"><p style="text-align:center;">Loading...</p></div>`;
  }

  return html`
    <div class="login-container">
      <h1>Sign Up for Lite Reader</h1>
      <form class="login-form" novalidate onSubmit=${onSubmit} data-testid="signup-form">
        <div class="form-group">
          <label for="email">Email:</label>
          <input
            type="email"
            id="email"
            name="email"
            class="form-control"
            data-testid="signup-email"
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
            data-testid="signup-password"
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
            data-testid="signup-confirm-password"
            value=${confirm}
            onInput=${(e) => setConfirm(e.target.value)}
            required
          />
        </div>
        <button
          type="submit"
          class="btn-primary"
          data-testid="signup-submit"
          disabled=${pending}
        >Sign Up</button>
        <div class="error-message" data-testid="signup-error" role="alert">${error}</div>
      </form>
      <div style="text-align: center; margin-top: 20px;">
        <a href="/login.html" style="color: #000000; text-decoration: underline; text-underline-offset: 3px;">Already have an account? Log in</a>
      </div>
    </div>
  `;
}
