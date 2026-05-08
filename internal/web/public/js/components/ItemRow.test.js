import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, cleanup, fireEvent, waitFor } from '@testing-library/preact';
import { html } from '../util/html.js';

const updateMock = vi.fn();
vi.mock('../api/items.js', () => ({ update: (...a) => updateMock(...a) }));

const { ItemRow } = await import('./ItemRow.js');
const { currentItem } = await import('../state.js');

const baseItem = {
  id: 42,
  title: 'Hello',
  is_new: true,
  starred: false,
  timestamp: new Date(Date.now() - 60_000).toISOString(),
};

describe('ItemRow', () => {
  beforeEach(() => { updateMock.mockReset(); currentItem.value = null; });
  afterEach(() => cleanup());

  it('renders title with data-testid', () => {
    render(html`<${ItemRow} item=${baseItem} onChanged=${() => {}} />`);
    expect(screen.getByTestId('item-row-title').textContent).toBe('Hello');
  });

  it('click row sets currentItem and marks read', async () => {
    updateMock.mockResolvedValue(null);
    const onChanged = vi.fn();
    render(html`<${ItemRow} item=${baseItem} onChanged=${onChanged} />`);
    fireEvent.click(screen.getByTestId('item-row'));
    expect(currentItem.value.id).toBe(42);
    await waitFor(() => expect(updateMock).toHaveBeenCalledWith(42, { is_new: false, starred: false }));
    expect(onChanged).toHaveBeenCalled();
  });

  it('star toggle round-trips through items.update', async () => {
    updateMock.mockResolvedValue(null);
    const onChanged = vi.fn();
    render(html`<${ItemRow} item=${baseItem} onChanged=${onChanged} />`);
    fireEvent.click(screen.getByTestId('item-row-star'));
    await waitFor(() => expect(updateMock).toHaveBeenCalledWith(42, { is_new: true, starred: true }));
    expect(onChanged.mock.calls[0][0].starred).toBe(true);
  });

  it('toggle-read button flips is_new', async () => {
    updateMock.mockResolvedValue(null);
    const read = { ...baseItem, is_new: false };
    const onChanged = vi.fn();
    render(html`<${ItemRow} item=${read} onChanged=${onChanged} />`);
    fireEvent.click(screen.getByTestId('item-row-toggle-read'));
    await waitFor(() => expect(updateMock).toHaveBeenCalledWith(42, { is_new: true, starred: false }));
  });

  it('does not call update on click when already read', () => {
    const read = { ...baseItem, is_new: false };
    render(html`<${ItemRow} item=${read} onChanged=${() => {}} />`);
    fireEvent.click(screen.getByTestId('item-row'));
    expect(updateMock).not.toHaveBeenCalled();
  });
});
