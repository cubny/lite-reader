import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, fireEvent, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

const addMock = vi.fn();
const listMock = vi.fn();
const fetchNewMock = vi.fn();
vi.mock('../api/feeds.js', () => ({
  add: (...a) => addMock(...a),
  list: (...a) => listMock(...a),
  fetchNew: (...a) => fetchNewMock(...a),
}));

const { AddFeedForm } = await import('./AddFeedForm.js');
const { feeds } = await import('../state.js');

function setInputValue(input, value) {
  input.value = value;
}

describe('AddFeedForm', () => {
  beforeEach(() => {
    addMock.mockReset();
    listMock.mockReset();
    fetchNewMock.mockReset();
    fetchNewMock.mockResolvedValue(null);
    feeds.value = [];
  });
  afterEach(() => cleanup());

  it('reveal-on-click → enter URL → submit', async () => {
    addMock.mockResolvedValue({ id: 1 });
    listMock.mockResolvedValue([{ id: 1, title: 'X' }]);
    render(html`<${AddFeedForm} />`);

    // First click: reveal input
    fireEvent.click(screen.getByTestId('add-feed-submit'));
    const input = screen.getByTestId('add-feed-url');
    setInputValue(input, 'http://example.com/rss.xml');
    fireEvent.keyDown(input, { key: 'Enter' });

    await waitFor(() => expect(addMock).toHaveBeenCalledWith('http://example.com/rss.xml'));
    await waitFor(() => expect(feeds.value).toEqual([{ id: 1, title: 'X' }]));
  });

  it('invalid URL → shows error, does NOT call add', async () => {
    render(html`<${AddFeedForm} />`);
    fireEvent.click(screen.getByTestId('add-feed-submit'));
    const input = screen.getByTestId('add-feed-url');
    setInputValue(input, 'not a valid url with spaces');
    fireEvent.keyDown(input, { key: 'Enter' });

    await waitFor(() => expect(screen.getByTestId('add-feed-error')).toBeInTheDocument());
    expect(addMock).not.toHaveBeenCalled();
  });

  it('shows pending spinner icon while add is in flight', async () => {
    let resolveAdd;
    addMock.mockReturnValue(new Promise((r) => { resolveAdd = r; }));
    listMock.mockResolvedValue([]);
    render(html`<${AddFeedForm} />`);

    fireEvent.click(screen.getByTestId('add-feed-submit'));
    const input = screen.getByTestId('add-feed-url');
    setInputValue(input, 'http://x.com/r.xml');
    fireEvent.keyDown(input, { key: 'Enter' });

    await waitFor(() => {
      const btn = screen.getByTestId('add-feed-submit');
      expect(btn.querySelector('i.icon-spin')).toBeTruthy();
    });
    resolveAdd({ id: 1 });
  });
});
