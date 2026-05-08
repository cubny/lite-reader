import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, cleanup, fireEvent } from '@testing-library/preact';
import { html } from '../util/html.js';
import { FeedItem } from './FeedItem.js';
import { selection } from '../state.js';

describe('FeedItem', () => {
  beforeEach(() => {
    selection.value = { kind: 'unread' };
    vi.stubGlobal('location', { hash: '#/', assign: vi.fn() });
  });
  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  const feed = { id: 7, title: 'Hacker News', url: 'http://x', unread_count: 3 };

  it('renders title, unread count, delete with data-testid', () => {
    render(html`<${FeedItem} feed=${feed} onDelete=${() => {}} />`);
    expect(screen.getByTestId('feed-item')).toBeInTheDocument();
    expect(screen.getByTestId('feed-item-title').textContent).toBe('Hacker News');
    expect(screen.getByTestId('feed-item-unread-count').textContent).toBe('3');
    expect(screen.getByTestId('feed-item-delete')).toBeInTheDocument();
  });

  it('click selects the feed', () => {
    render(html`<${FeedItem} feed=${feed} onDelete=${() => {}} />`);
    fireEvent.click(screen.getByTestId('feed-item'));
    expect(selection.value).toEqual({ kind: 'feed', id: 7 });
  });

  it('delete button calls onDelete (does not select)', () => {
    const onDelete = vi.fn();
    render(html`<${FeedItem} feed=${feed} onDelete=${onDelete} />`);
    fireEvent.click(screen.getByTestId('feed-item-delete'));
    expect(onDelete).toHaveBeenCalledWith(feed);
    expect(selection.value).toEqual({ kind: 'unread' });
  });

  it('falls back to url when title missing', () => {
    render(html`<${FeedItem} feed=${{ id: 1, url: 'http://only-url' }} onDelete=${() => {}} />`);
    expect(screen.getByTestId('feed-item-title').textContent).toBe('http://only-url');
  });
});
