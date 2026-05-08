import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, fireEvent, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

const fetchFeedMock = vi.fn();
const markReadMock = vi.fn();
const markUnreadMock = vi.fn();
const removeFeedMock = vi.fn();
const itemsMock = vi.fn();
const listFeedsMock = vi.fn();
const logoutMock = vi.fn();

vi.mock('../api/feeds.js', () => ({
  fetchNew: (...a) => fetchFeedMock(...a),
  markRead: (...a) => markReadMock(...a),
  markUnread: (...a) => markUnreadMock(...a),
  remove: (...a) => removeFeedMock(...a),
  items: (...a) => itemsMock(...a),
  list: (...a) => listFeedsMock(...a),
}));
vi.mock('../api/items.js', () => ({
  unread: () => Promise.resolve([]),
  starred: () => Promise.resolve([]),
  unreadCount: () => Promise.resolve({ count: 0 }),
  starredCount: () => Promise.resolve({ count: 0 }),
}));
vi.mock('../api/auth.js', () => ({
  logout: (...a) => logoutMock(...a),
}));

const { FeedBar } = await import('./FeedBar.js');
const { selection, feeds, items } = await import('../state.js');

describe('FeedBar', () => {
  beforeEach(() => {
    fetchFeedMock.mockReset();
    markReadMock.mockReset();
    markUnreadMock.mockReset();
    removeFeedMock.mockReset();
    itemsMock.mockReset();
    listFeedsMock.mockReset();
    logoutMock.mockReset();
    feeds.value = [{ id: 9, title: 'My Feed', url: 'http://x' }];
    items.value = [];
    selection.value = { kind: 'unread' };
    vi.stubGlobal('location', { hash: '#/', assign: vi.fn() });
  });
  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('renders title and all five action buttons', () => {
    render(html`<${FeedBar} />`);
    expect(screen.getByTestId('toolbar-title').textContent).toBe('Unread');
    expect(screen.getByTestId('toolbar-refresh')).toBeInTheDocument();
    expect(screen.getByTestId('toolbar-mark-read')).toBeInTheDocument();
    expect(screen.getByTestId('toolbar-mark-unread')).toBeInTheDocument();
    expect(screen.getByTestId('toolbar-logout')).toBeInTheDocument();
    expect(screen.getByTestId('toolbar-remove')).toBeInTheDocument();
  });

  it('shows feed title when scope is feed', () => {
    selection.value = { kind: 'feed', id: 9 };
    render(html`<${FeedBar} />`);
    expect(screen.getByTestId('toolbar-title').textContent).toBe('My Feed');
  });

  it('logout button always invokes auth.logout', () => {
    render(html`<${FeedBar} />`);
    fireEvent.click(screen.getByTestId('toolbar-logout'));
    expect(logoutMock).toHaveBeenCalled();
  });

  it('refresh on feed scope: calls fetchNew + reloads items', async () => {
    selection.value = { kind: 'feed', id: 9 };
    fetchFeedMock.mockResolvedValue(null);
    itemsMock.mockResolvedValue([{ id: 1, title: 'a' }]);
    render(html`<${FeedBar} />`);
    fireEvent.click(screen.getByTestId('toolbar-refresh'));
    await waitFor(() => expect(fetchFeedMock).toHaveBeenCalledWith(9));
    await waitFor(() => expect(items.value).toHaveLength(1));
  });

  it('refresh on smart-folder scope: no-ops (no API call)', async () => {
    render(html`<${FeedBar} />`);
    fireEvent.click(screen.getByTestId('toolbar-refresh'));
    await new Promise((r) => setTimeout(r, 10));
    expect(fetchFeedMock).not.toHaveBeenCalled();
  });

  it('read all on feed scope: calls markRead', async () => {
    selection.value = { kind: 'feed', id: 9 };
    markReadMock.mockResolvedValue(null);
    itemsMock.mockResolvedValue([]);
    render(html`<${FeedBar} />`);
    fireEvent.click(screen.getByTestId('toolbar-mark-read'));
    await waitFor(() => expect(markReadMock).toHaveBeenCalledWith(9));
  });

  it('unread all on feed scope: calls markUnread', async () => {
    selection.value = { kind: 'feed', id: 9 };
    markUnreadMock.mockResolvedValue(null);
    itemsMock.mockResolvedValue([]);
    render(html`<${FeedBar} />`);
    fireEvent.click(screen.getByTestId('toolbar-mark-unread'));
    await waitFor(() => expect(markUnreadMock).toHaveBeenCalledWith(9));
  });

  it('remove on feed scope: confirm dialog → calls remove + clears selection', async () => {
    selection.value = { kind: 'feed', id: 9 };
    removeFeedMock.mockResolvedValue(null);
    listFeedsMock.mockResolvedValue([]);
    render(html`<${FeedBar} />`);
    fireEvent.click(screen.getByTestId('toolbar-remove'));
    expect(screen.getByTestId('confirm-dialog')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('confirm-yes'));
    await waitFor(() => expect(removeFeedMock).toHaveBeenCalledWith(9));
    await waitFor(() => expect(selection.value).toEqual({ kind: 'unread' }));
  });
});
