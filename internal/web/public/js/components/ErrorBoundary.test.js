import { describe, it, expect, afterEach, vi } from 'vitest';
import { render, screen, cleanup } from '@testing-library/preact';
import { html } from '../util/html.js';
import { ErrorBoundary } from './ErrorBoundary.js';

function Boom() {
  throw new Error('kaboom');
}

describe('ErrorBoundary', () => {
  afterEach(() => cleanup());

  it('renders children when no error', () => {
    render(html`<${ErrorBoundary}><div data-testid="child">ok</div><//>`);
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('renders fallback when child throws', () => {
    const spy = vi.spyOn(console, 'error').mockImplementation(() => {});
    render(html`<${ErrorBoundary}><${Boom} /><//>`);
    expect(screen.getByTestId('error-fallback')).toBeInTheDocument();
    expect(screen.getByTestId('error-fallback').textContent).toContain('kaboom');
    spy.mockRestore();
  });
});
