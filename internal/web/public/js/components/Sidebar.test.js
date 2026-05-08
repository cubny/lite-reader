import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

const listMock = vi.fn();
vi.mock('../api/feeds.js', () => ({
  list: (...a) => listMock(...a),
  add: vi.fn(),
  fetchNew: vi.fn(),
}));
vi.mock('../api/items.js', () => ({
  unreadCount: () => Promise.resolve({ count: 0 }),
  starredCount: () => Promise.resolve({ count: 0 }),
}));

const { Sidebar } = await import('./Sidebar.js');
const { feeds } = await import('../state.js');

describe('Sidebar', () => {
  beforeEach(() => {
    listMock.mockReset();
    feeds.value = [];
    vi.stubGlobal('location', { hash: '#/', assign: vi.fn() });
  });
  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('renders sidebar with smart folders, add feed, feeds', async () => {
    listMock.mockResolvedValue([]);
    render(html`<${Sidebar} />`);
    expect(screen.getByTestId('sidebar')).toBeInTheDocument();
    expect(screen.getByTestId('smart-folder-unread')).toBeInTheDocument();
    expect(screen.getByTestId('smart-folder-starred')).toBeInTheDocument();
    expect(screen.getByTestId('add-feed-form')).toBeInTheDocument();
  });

  it('on mount: loads feeds and populates state', async () => {
    listMock.mockResolvedValue([{ id: 1, title: 'A', unread_count: 2 }]);
    render(html`<${Sidebar} />`);
    await waitFor(() => expect(feeds.value).toHaveLength(1));
    expect(screen.getAllByTestId('feed-item')).toHaveLength(1);
  });
});
