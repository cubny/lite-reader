import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

vi.mock('../components/Sidebar.js', () => ({
  Sidebar: () => html`<div data-testid="mock-sidebar"></div>`,
}));
vi.mock('../components/Toolbar.js', () => ({
  Toolbar: () => html`<div data-testid="mock-toolbar"></div>`,
}));
vi.mock('../components/ItemList.js', () => ({
  ItemList: () => html`<div data-testid="mock-itemlist"></div>`,
}));
vi.mock('../components/ItemDetail.js', () => ({
  ItemDetail: () => html`<div data-testid="mock-itemdetail"></div>`,
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

  it('with token: renders 3-pane shell with all data-testid', () => {
    localStorage.setItem('token', 'tok');
    render(html`<${AppPage} />`);
    expect(screen.getByTestId('pane-sidebar')).toBeInTheDocument();
    expect(screen.getByTestId('pane-list')).toBeInTheDocument();
    expect(screen.getByTestId('pane-detail')).toBeInTheDocument();
  });

  it('without token: redirects to /login.html', async () => {
    render(html`<${AppPage} />`);
    await waitFor(() => {
      expect(location.assign).toHaveBeenCalledWith('/login.html');
    });
  });
});
