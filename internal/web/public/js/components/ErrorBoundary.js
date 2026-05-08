import { Component } from 'preact';
import { html } from '../util/html.js';

export class ErrorBoundary extends Component {
  constructor(props) {
    super(props);
    this.state = { error: null };
  }

  static getDerivedStateFromError(error) {
    return { error };
  }

  componentDidCatch(error) {
    this.setState({ error });
  }

  render(props, state) {
    if (state.error) {
      return html`
        <div class="error-fallback" data-testid="error-fallback" role="alert">
          <h2>Something went wrong</h2>
          <pre>${state.error.message || String(state.error)}</pre>
        </div>
      `;
    }
    return props.children;
  }
}
