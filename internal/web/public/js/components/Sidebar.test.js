import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, fireEvent, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

const listMock = vi.fn();
const removeMock = vi.fn();
const logoutMock = vi.fn();
vi.mock('../api/feeds.js', () => ({
  list: (...a) => listMock(...a),
  remove: (...a) => removeMock(...a),
  add: vi.fn(),
}));
vi.mock('../api/auth.js', () => ({
  logout: (...a) => logoutMock(...a),
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
    removeMock.mockReset();
    logoutMock.mockReset();
    feeds.value = [];
    vi.stubGlobal('location', { hash: '#/', assign: vi.fn() });
  });
  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('renders sidebar with logout, smart folders, add feed, feed list', async () => {
    listMock.mockResolvedValue([]);
    render(html`<${Sidebar} />`);
    expect(screen.getByTestId('sidebar')).toBeInTheDocument();
    expect(screen.getByTestId('logout-button')).toBeInTheDocument();
    expect(screen.getByTestId('smart-folders')).toBeInTheDocument();
    expect(screen.getByTestId('add-feed-form')).toBeInTheDocument();
    expect(screen.getByTestId('feed-list')).toBeInTheDocument();
  });

  it('logout button calls auth.logout', async () => {
    listMock.mockResolvedValue([]);
    render(html`<${Sidebar} />`);
    fireEvent.click(screen.getByTestId('logout-button'));
    expect(logoutMock).toHaveBeenCalled();
  });

  it('on mount: loads feeds and populates state', async () => {
    listMock.mockResolvedValue([{ id: 1, title: 'A', unread_count: 2 }]);
    render(html`<${Sidebar} />`);
    await waitFor(() => expect(feeds.value).toHaveLength(1));
    expect(screen.getAllByTestId('feed-item')).toHaveLength(1);
  });

  it('delete feed: confirm → calls remove → refreshes', async () => {
    listMock.mockResolvedValueOnce([{ id: 1, title: 'A' }]);
    listMock.mockResolvedValueOnce([]);
    removeMock.mockResolvedValue(null);
    render(html`<${Sidebar} />`);
    await waitFor(() => expect(screen.getAllByTestId('feed-item')).toHaveLength(1));
    fireEvent.click(screen.getByTestId('feed-item-delete'));
    expect(screen.getByTestId('confirm-dialog')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('confirm-yes'));
    await waitFor(() => expect(removeMock).toHaveBeenCalledWith(1));
    await waitFor(() => expect(feeds.value).toEqual([]));
  });
});
