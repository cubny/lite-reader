import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

vi.mock('../components/Sidebar.js', () => ({
  Sidebar: () => html`<div data-testid="mock-sidebar"></div>`,
}));
vi.mock('../components/FeedBar.js', () => ({
  FeedBar: () => html`<div data-testid="mock-feedbar"></div>`,
}));
vi.mock('../components/ItemList.js', () => ({
  ItemList: () => html`<div data-testid="mock-itemlist"></div>`,
}));
vi.mock('../components/ErrorBoundary.js', () => ({
  ErrorBoundary: ({ children }) => html`<div data-testid="mock-eb">${children}</div>`,
}));

const { AppPage } = await import('./AppPage.js');

describe('AppPage', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.stubGlobal('location', { assign: vi.fn(), pathname: '/', href: '/', hash: '#/' });
  });
  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('with token: renders sidebar + feedbar + item list', () => {
    localStorage.setItem('token', 'tok');
    render(html`<${AppPage} />`);
    expect(screen.getByTestId('pane-sidebar')).toBeInTheDocument();
    expect(screen.getByTestId('mock-sidebar')).toBeInTheDocument();
    expect(screen.getByTestId('mock-feedbar')).toBeInTheDocument();
    expect(screen.getByTestId('mock-itemlist')).toBeInTheDocument();
  });

  it('without token: redirects to /login.html', async () => {
    render(html`<${AppPage} />`);
    await waitFor(() => {
      expect(location.assign).toHaveBeenCalledWith('/login.html');
    });
  });
});
