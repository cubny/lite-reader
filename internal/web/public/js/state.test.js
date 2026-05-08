import { describe, it, expect } from 'vitest';
import * as state from './state.js';

describe('state signals', () => {
  it('exports expected signals', () => {
    for (const name of ['token', 'feeds', 'selection', 'items', 'currentItem', 'unreadCount', 'starredCount', 'toast', 'loading']) {
      expect(state[name]).toBeDefined();
      expect(typeof state[name].value !== 'undefined').toBe(true);
    }
  });

  it('signals are mutable', () => {
    state.feeds.value = [{ id: 1 }];
    expect(state.feeds.value).toEqual([{ id: 1 }]);
    state.feeds.value = [];
  });
});
