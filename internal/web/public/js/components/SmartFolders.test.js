import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, fireEvent, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

const unreadMock = vi.fn();
const starredMock = vi.fn();
vi.mock('../api/items.js', () => ({
  unreadCount: () => unreadMock(),
  starredCount: () => starredMock(),
}));

const { SmartFolders } = await import('./SmartFolders.js');
const { unreadCount, starredCount, selection } = await import('../state.js');

describe('SmartFolders', () => {
  beforeEach(() => {
    unreadMock.mockReset();
    starredMock.mockReset();
    unreadCount.value = 0;
    starredCount.value = 0;
    selection.value = { kind: 'unread' };
    vi.stubGlobal('location', { hash: '#/', assign: vi.fn() });
  });
  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('renders unread + starred entries with counts', async () => {
    unreadMock.mockResolvedValue({ count: 12 });
    starredMock.mockResolvedValue({ count: 4 });
    render(html`<${SmartFolders} />`);
    await waitFor(() => expect(unreadCount.value).toBe(12));
    await waitFor(() => expect(starredCount.value).toBe(4));
    expect(screen.getByTestId('smart-folder-unread-count').textContent).toBe('12');
    expect(screen.getByTestId('smart-folder-starred-count').textContent).toBe('4');
  });

  it('click unread → selects unread', async () => {
    unreadMock.mockResolvedValue({ count: 0 });
    starredMock.mockResolvedValue({ count: 0 });
    render(html`<${SmartFolders} />`);
    fireEvent.click(screen.getByTestId('smart-folder-unread'));
    expect(selection.value).toEqual({ kind: 'unread' });
  });

  it('click starred → selects starred', async () => {
    unreadMock.mockResolvedValue({ count: 0 });
    starredMock.mockResolvedValue({ count: 0 });
    render(html`<${SmartFolders} />`);
    fireEvent.click(screen.getByTestId('smart-folder-starred'));
    expect(selection.value).toEqual({ kind: 'starred' });
  });
});
