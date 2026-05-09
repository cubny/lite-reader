import { useState, useEffect } from 'preact/hooks';

import { html } from '../util/html.js';
import { login, checkNeedsSetup } from '../api/auth.js';

export function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);
  const [signupSuccess, setSignupSuccess] = useState(false);
  const [setupSuccess, setSetupSuccess] = useState(false);
  const [loading, setLoading] = useState(true);
  const [allowSignup, setAllowSignup] = useState(false);

  useEffect(() => {
    checkNeedsSetup().then((data) => {
      if (data.needs_setup) {
        location.assign('/setup.html');
        return;
      }
      setAllowSignup(data.allow_signup === true);
      if (sessionStorage.getItem('signupSuccess') === 'true') {
        setSignupSuccess(true);
        sessionStorage.removeItem('signupSuccess');
      }
      if (sessionStorage.getItem('setupSuccess') === 'true') {
        setSetupSuccess(true);
        sessionStorage.removeItem('setupSuccess');
      }
      setLoading(false);
    }).catch(() => {
      setLoading(false);
    });
  }, []);

  async function onSubmit(e) {
    e.preventDefault();
    setError('');
    setPending(true);
    try {
      await login({ email, password });
      location.assign('/');
    } catch (err) {
      setError(err.message || 'Login failed');
      setPending(false);
    }
  }

  if (loading) {
    return html`<div class="login-container"><p style="text-align:center;">Loading...</p></div>`;
  }

  return html`
    <div class="login-container">
      ${signupSuccess && html`
        <aside data-testid="signup-success" class="alert alert-success">
          Account created successfully. Please log in.
        </aside>
      `}
      ${setupSuccess && html`
        <aside data-testid="setup-success" class="alert alert-success">
          Admin account created. Please log in.
        </aside>
      `}
      <h1>Login to Lite Reader</h1>
      <form class="login-form" onSubmit=${onSubmit} data-testid="login-form">
        <div class="form-group">
          <label for="email">Email:</label>
          <input
            type="email"
            id="email"
            name="email"
            class="form-control"
            data-testid="login-email"
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
            data-testid="login-password"
            value=${password}
            onInput=${(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button
          type="submit"
          class="btn-primary"
          data-testid="login-submit"
          disabled=${pending}
        >Login</button>
        <div class="error-message" data-testid="login-error" role="alert">${error}</div>
      </form>
      ${allowSignup && html`
        <div style="text-align: center; margin-top: 20px;">
          <a href="/signup.html" style="color: #6C8A46;">Don't have an account? Sign up</a>
        </div>
      `}
    </div>
  `;
}
