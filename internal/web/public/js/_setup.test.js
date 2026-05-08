import { describe, it, expect } from 'vitest';
import { h } from 'preact';
import htm from 'htm';

const html = htm.bind(h);

describe('vitest harness', () => {
  it('preact + htm load and render a vnode', () => {
    const vnode = html`<div class="hi">hello</div>`;
    expect(vnode.type).toBe('div');
    expect(vnode.props.class).toBe('hi');
    expect(vnode.props.children).toBe('hello');
  });
});
