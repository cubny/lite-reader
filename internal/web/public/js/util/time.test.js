import { describe, it, expect } from 'vitest';
import { relativeTime } from './time.js';

const now = new Date('2026-05-08T12:00:00Z').getTime();

describe('relativeTime', () => {
  it('< 1 min → just now', () => {
    expect(relativeTime(now - 30 * 1000, now)).toBe('just now');
  });
  it('minutes', () => {
    expect(relativeTime(now - 60 * 1000, now)).toBe('1 minute ago');
    expect(relativeTime(now - 5 * 60 * 1000, now)).toBe('5 minutes ago');
  });
  it('hours', () => {
    expect(relativeTime(now - 3600 * 1000, now)).toBe('1 hour ago');
    expect(relativeTime(now - 4 * 3600 * 1000, now)).toBe('4 hours ago');
  });
  it('days/weeks/months/years', () => {
    expect(relativeTime(now - 86400 * 1000, now)).toBe('1 day ago');
    expect(relativeTime(now - 7 * 86400 * 1000, now)).toBe('1 week ago');
    expect(relativeTime(now - 31 * 86400 * 1000, now)).toBe('1 month ago');
    expect(relativeTime(now - 366 * 86400 * 1000, now)).toBe('1 year ago');
  });
  it('invalid → empty', () => {
    expect(relativeTime('not a date', now)).toBe('');
  });
});
