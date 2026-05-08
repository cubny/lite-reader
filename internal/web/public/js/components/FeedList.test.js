import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { render, screen, cleanup } from '@testing-library/preact';
import { html } from '../util/html.js';
import { FeedList } from './FeedList.js';
import { feeds, selection } from '../state.js';

describe('FeedList', () => {
  beforeEach(() => {
    feeds.value = [];
    selection.value = { kind: 'unread' };
  });
  afterEach(() => cleanup());

  it('empty list renders nothing', () => {
    const { container } = render(html`<ul>${html`<${FeedList} />`}</ul>`);
    expect(container.querySelectorAll('li.feed').length).toBe(0);
  });

  it('renders one row per feed', () => {
    feeds.value = [
      { id: 1, title: 'A', unread_count: 2 },
      { id: 2, title: 'B', unread_count: 0 },
    ];
    render(html`<ul>${html`<${FeedList} />`}</ul>`);
    expect(screen.getAllByTestId('feed-item')).toHaveLength(2);
  });
});
