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

  it('renders title, unread count with data-testid', () => {
    render(html`<${FeedItem} feed=${feed} />`);
    expect(screen.getByTestId('feed-item')).toBeInTheDocument();
    expect(screen.getByTestId('feed-item-title').textContent).toBe('Hacker News');
    expect(screen.getByTestId('feed-item-unread-count').textContent).toBe('3');
  });

  it('click selects the feed', () => {
    render(html`<${FeedItem} feed=${feed} />`);
    fireEvent.click(screen.getByTestId('feed-item'));
    expect(selection.value).toEqual({ kind: 'feed', id: 7 });
  });

  it('falls back to url when title missing', () => {
    render(html`<${FeedItem} feed=${{ id: 1, url: 'http://only-url' }} />`);
    expect(screen.getByTestId('feed-item-title').textContent).toBe('http://only-url');
  });

  it('legacy DOM: <li class="feed"> with .count, <i>, .feedtitle', () => {
    const { container } = render(html`<${FeedItem} feed=${feed} />`);
    const li = container.querySelector('li.feed');
    expect(li).toBeTruthy();
    expect(li.querySelector('.count')).toBeTruthy();
    expect(li.querySelector('i')).toBeTruthy();
    expect(li.querySelector('.feedtitle')).toBeTruthy();
  });
});
