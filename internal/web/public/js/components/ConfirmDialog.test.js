import { describe, it, expect, vi, afterEach } from 'vitest';
import { render, screen, cleanup, fireEvent } from '@testing-library/preact';
import { html } from '../util/html.js';
import { ConfirmDialog } from './ConfirmDialog.js';

describe('ConfirmDialog', () => {
  afterEach(() => cleanup());

  it('renders message and action buttons', () => {
    render(html`<${ConfirmDialog} message="Are you sure?" onConfirm=${() => {}} onCancel=${() => {}} />`);
    expect(screen.getByTestId('confirm-dialog')).toBeInTheDocument();
    expect(screen.getByTestId('confirm-message').textContent).toBe('Are you sure?');
    expect(screen.getByTestId('confirm-yes')).toBeInTheDocument();
    expect(screen.getByTestId('confirm-no')).toBeInTheDocument();
  });

  it('click yes calls onConfirm', () => {
    const onConfirm = vi.fn();
    render(html`<${ConfirmDialog} message="x" onConfirm=${onConfirm} onCancel=${() => {}} />`);
    fireEvent.click(screen.getByTestId('confirm-yes'));
    expect(onConfirm).toHaveBeenCalled();
  });

  it('click no calls onCancel', () => {
    const onCancel = vi.fn();
    render(html`<${ConfirmDialog} message="x" onConfirm=${() => {}} onCancel=${onCancel} />`);
    fireEvent.click(screen.getByTestId('confirm-no'));
    expect(onCancel).toHaveBeenCalled();
  });

  it('Escape key cancels', () => {
    const onCancel = vi.fn();
    render(html`<${ConfirmDialog} message="x" onConfirm=${() => {}} onCancel=${onCancel} />`);
    fireEvent.keyDown(window, { key: 'Escape' });
    expect(onCancel).toHaveBeenCalled();
  });

  it('Enter key confirms', () => {
    const onConfirm = vi.fn();
    render(html`<${ConfirmDialog} message="x" onConfirm=${onConfirm} onCancel=${() => {}} />`);
    fireEvent.keyDown(window, { key: 'Enter' });
    expect(onConfirm).toHaveBeenCalled();
  });
});
