import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, waitFor } from '@testing-library/preact';
import userEvent from '@testing-library/user-event';
import { html } from '../util/html.js';

const addMock = vi.fn();
const listMock = vi.fn();
vi.mock('../api/feeds.js', () => ({
  add: (...a) => addMock(...a),
  list: (...a) => listMock(...a),
}));

const { AddFeedForm } = await import('./AddFeedForm.js');
const { feeds } = await import('../state.js');

describe('AddFeedForm', () => {
  beforeEach(() => {
    addMock.mockReset();
    listMock.mockReset();
    feeds.value = [];
  });
  afterEach(() => cleanup());

  it('valid URL → calls add and refreshes list', async () => {
    addMock.mockResolvedValue({ id: 1 });
    listMock.mockResolvedValue([{ id: 1, title: 'X' }]);
    const user = userEvent.setup();
    render(html`<${AddFeedForm} />`);
    await user.type(screen.getByTestId('add-feed-url'), 'http://example.com/rss.xml');
    await user.click(screen.getByTestId('add-feed-submit'));
    await waitFor(() => expect(addMock).toHaveBeenCalledWith('http://example.com/rss.xml'));
    await waitFor(() => expect(feeds.value).toEqual([{ id: 1, title: 'X' }]));
  });

  it('invalid URL → shows error, does NOT call add', async () => {
    const user = userEvent.setup();
    render(html`<${AddFeedForm} />`);
    await user.type(screen.getByTestId('add-feed-url'), 'not-a-url');
    await user.click(screen.getByTestId('add-feed-submit'));
    await waitFor(() => expect(screen.getByTestId('add-feed-error')).toBeInTheDocument());
    expect(addMock).not.toHaveBeenCalled();
  });

  it('disables submit while pending', async () => {
    let resolve;
    addMock.mockReturnValue(new Promise((r) => { resolve = r; }));
    const user = userEvent.setup();
    render(html`<${AddFeedForm} />`);
    await user.type(screen.getByTestId('add-feed-url'), 'http://x.com/r.xml');
    await user.click(screen.getByTestId('add-feed-submit'));
    await waitFor(() => expect(screen.getByTestId('add-feed-submit').disabled).toBe(true));
    resolve({ id: 1 });
    listMock.mockResolvedValue([]);
  });
});
