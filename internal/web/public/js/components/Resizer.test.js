import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, fireEvent } from '@testing-library/preact';
import { html } from '../util/html.js';
import { Resizer } from './Resizer.js';

describe('Resizer', () => {
  beforeEach(() => {
    localStorage.clear();
    document.documentElement.style.setProperty('--sidebar-w', '200px');
  });
  afterEach(() => cleanup());

  it('renders separator with data-testid', () => {
    render(html`<${Resizer} id="sidebar" cssVar="--sidebar-w" />`);
    expect(screen.getByTestId('resizer-sidebar')).toBeInTheDocument();
  });

  it('restores stored width on mount', () => {
    localStorage.setItem('resizer:sidebar', '300');
    render(html`<${Resizer} id="sidebar" cssVar="--sidebar-w" min=${100} max=${500} />`);
    expect(document.documentElement.style.getPropertyValue('--sidebar-w')).toBe('300px');
  });

  it('clamps stored width to [min,max]', () => {
    localStorage.setItem('resizer:sidebar', '9999');
    render(html`<${Resizer} id="sidebar" cssVar="--sidebar-w" min=${100} max=${500} />`);
    expect(document.documentElement.style.getPropertyValue('--sidebar-w')).toBe('500px');
  });

  it('drag updates css var and persists to localStorage', () => {
    render(html`<${Resizer} id="sidebar" cssVar="--sidebar-w" min=${100} max=${500} />`);
    const handle = screen.getByTestId('resizer-sidebar');

    fireEvent.mouseDown(handle, { clientX: 100 });
    fireEvent.mouseMove(document, { clientX: 250 });
    fireEvent.mouseUp(document);

    const stored = localStorage.getItem('resizer:sidebar');
    expect(stored).not.toBeNull();
    expect(Number(stored)).toBeGreaterThan(200);
  });
});
