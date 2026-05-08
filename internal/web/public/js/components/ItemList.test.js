import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

const unreadMock = vi.fn();
const starredMock = vi.fn();
const feedItemsMock = vi.fn();
vi.mock('../api/items.js', () => ({
  unread: () => unreadMock(),
  starred: () => starredMock(),
  update: vi.fn(),
}));
vi.mock('../api/feeds.js', () => ({ items: (id) => feedItemsMock(id) }));

const { ItemList } = await import('./ItemList.js');
const { selection, items, currentItem } = await import('../state.js');

const sample = { id: 1, title: 'A', is_new: true, starred: false, timestamp: new Date().toISOString() };

describe('ItemList', () => {
  beforeEach(() => {
    unreadMock.mockReset();
    starredMock.mockReset();
    feedItemsMock.mockReset();
    items.value = [];
    currentItem.value = null;
    selection.value = { kind: 'unread' };
  });
  afterEach(() => cleanup());

  it('on unread scope: calls items.unread', async () => {
    unreadMock.mockResolvedValue([sample]);
    render(html`<${ItemList} />`);
    await waitFor(() => expect(unreadMock).toHaveBeenCalled());
    await waitFor(() => expect(screen.getAllByTestId('item-row')).toHaveLength(1));
  });

  it('on starred scope: calls items.starred', async () => {
    selection.value = { kind: 'starred' };
    starredMock.mockResolvedValue([]);
    render(html`<${ItemList} />`);
    await waitFor(() => expect(starredMock).toHaveBeenCalled());
  });

  it('on feed scope: calls feeds.items(id)', async () => {
    selection.value = { kind: 'feed', id: 9 };
    feedItemsMock.mockResolvedValue([]);
    render(html`<${ItemList} />`);
    await waitFor(() => expect(feedItemsMock).toHaveBeenCalledWith(9));
  });
});
