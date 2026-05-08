import { describe, it, expect } from 'vitest';
import { parse } from './router.js';

describe('router.parse', () => {
  it('empty hash → unread', () => {
    expect(parse('')).toEqual({ name: 'unread', params: {} });
    expect(parse('#/')).toEqual({ name: 'unread', params: {} });
    expect(parse('#')).toEqual({ name: 'unread', params: {} });
  });

  it('#/starred', () => {
    expect(parse('#/starred')).toEqual({ name: 'starred', params: {} });
  });

  it('#/feed/:id', () => {
    expect(parse('#/feed/42')).toEqual({ name: 'feed', params: { id: '42' } });
  });

  it('#/feed/:id/item/:itemId', () => {
    expect(parse('#/feed/42/item/9')).toEqual({
      name: 'feed-item',
      params: { id: '42', itemId: '9' },
    });
  });

  it('unknown route', () => {
    expect(parse('#/garbage')).toEqual({ name: 'unknown', params: { path: '/garbage' } });
  });
});
