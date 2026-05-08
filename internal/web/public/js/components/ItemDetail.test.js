import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { render, screen, cleanup } from '@testing-library/preact';
import { html } from '../util/html.js';
import { ItemDetail } from './ItemDetail.js';
import { currentItem } from '../state.js';

describe('ItemDetail', () => {
  beforeEach(() => { currentItem.value = null; });
  afterEach(() => cleanup());

  it('shows empty state when no item selected', () => {
    render(html`<${ItemDetail} />`);
    expect(screen.getByTestId('item-detail-empty')).toBeInTheDocument();
  });

  it('renders title and source link in parent document', () => {
    currentItem.value = { id: 1, title: 'Hello', desc: '<p>Body</p>', link: 'https://example.com/a' };
    render(html`<${ItemDetail} />`);
    expect(screen.getByTestId('item-detail-title').textContent).toBe('Hello');
    expect(screen.getByTestId('item-detail-link').getAttribute('href')).toBe('https://example.com/a');
  });

  it('iframe has empty sandbox attribute (no flags)', () => {
    currentItem.value = { id: 2, title: 'X', desc: '<p>p</p>', link: '' };
    render(html`<${ItemDetail} />`);
    const iframe = screen.getByTestId('item-detail-iframe');
    expect(iframe.tagName).toBe('IFRAME');
    expect(iframe.getAttribute('sandbox')).toBe('');
  });

  it('XSS: script in srcdoc cannot execute in parent (window.__pwn undefined)', () => {
    delete globalThis.__pwn;
    currentItem.value = {
      id: 3,
      title: 'pwn',
      desc: '<script>window.__pwn=1</script>',
      link: '',
    };
    render(html`<${ItemDetail} />`);
    expect(globalThis.__pwn).toBeUndefined();
    const iframe = screen.getByTestId('item-detail-iframe');
    expect(iframe.getAttribute('sandbox')).toBe('');
  });

  it('RTL detection sets dir="rtl" in srcdoc when title is RTL', () => {
    currentItem.value = { id: 4, title: 'مرحبا', desc: '<p>x</p>', link: '' };
    render(html`<${ItemDetail} />`);
    const iframe = screen.getByTestId('item-detail-iframe');
    expect(iframe.getAttribute('srcdoc')).toMatch(/dir="rtl"/);
  });
});
