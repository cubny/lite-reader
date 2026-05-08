import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, fireEvent, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

const fetchNewMock = vi.fn();
const markReadMock = vi.fn();
const feedItemsMock = vi.fn();
vi.mock('../api/feeds.js', () => ({
  fetchNew: (...a) => fetchNewMock(...a),
  markRead: (...a) => markReadMock(...a),
  items: (...a) => feedItemsMock(...a),
}));
vi.mock('../api/items.js', () => ({
  unread: () => Promise.resolve([]),
  starred: () => Promise.resolve([]),
  unreadCount: () => Promise.resolve({ count: 0 }),
  starredCount: () => Promise.resolve({ count: 0 }),
}));

const { Toolbar } = await import('./Toolbar.js');
const { selection, feeds } = await import('../state.js');

describe('Toolbar', () => {
  beforeEach(() => {
    fetchNewMock.mockReset();
    markReadMock.mockReset();
    feedItemsMock.mockReset();
    feeds.value = [{ id: 7, title: 'TechCrunch' }];
    selection.value = { kind: 'unread' };
  });
  afterEach(() => cleanup());

  it('title reflects selection: unread', () => {
    selection.value = { kind: 'unread' };
    render(html`<${Toolbar} />`);
    expect(screen.getByTestId('toolbar-title').textContent).toBe('Unread');
  });

  it('title reflects selection: starred', () => {
    selection.value = { kind: 'starred' };
    render(html`<${Toolbar} />`);
    expect(screen.getByTestId('toolbar-title').textContent).toBe('Starred');
  });

  it('title reflects selection: feed name', () => {
    selection.value = { kind: 'feed', id: 7 };
    render(html`<${Toolbar} />`);
    expect(screen.getByTestId('toolbar-title').textContent).toBe('TechCrunch');
  });

  it('refresh+mark-read disabled on non-feed scope', () => {
    selection.value = { kind: 'unread' };
    render(html`<${Toolbar} />`);
    expect(screen.getByTestId('toolbar-refresh').disabled).toBe(true);
    expect(screen.getByTestId('toolbar-mark-read').disabled).toBe(true);
  });

  it('refresh on feed scope calls fetchNew', async () => {
    selection.value = { kind: 'feed', id: 7 };
    fetchNewMock.mockResolvedValue([]);
    feedItemsMock.mockResolvedValue([]);
    render(html`<${Toolbar} />`);
    fireEvent.click(screen.getByTestId('toolbar-refresh'));
    await waitFor(() => expect(fetchNewMock).toHaveBeenCalledWith(7));
  });

  it('mark-all on feed scope calls markRead', async () => {
    selection.value = { kind: 'feed', id: 7 };
    markReadMock.mockResolvedValue(null);
    feedItemsMock.mockResolvedValue([]);
    render(html`<${Toolbar} />`);
    fireEvent.click(screen.getByTestId('toolbar-mark-read'));
    await waitFor(() => expect(markReadMock).toHaveBeenCalledWith(7));
  });
});
