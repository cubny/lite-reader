const MINUTE = 60;
const HOUR = 3600;
const DAY = 86400;
const WEEK = 604800;
const MONTH = 2592000;
const YEAR = 31536000;

export function relativeTime(date, now = Date.now()) {
  const t = date instanceof Date ? date.getTime() : new Date(date).getTime();
  if (Number.isNaN(t)) return '';
  const diff = Math.max(0, Math.round((now - t) / 1000));
  if (diff < MINUTE) return 'just now';
  if (diff < HOUR) return plural(Math.floor(diff / MINUTE), 'minute');
  if (diff < DAY) return plural(Math.floor(diff / HOUR), 'hour');
  if (diff < WEEK) return plural(Math.floor(diff / DAY), 'day');
  if (diff < MONTH) return plural(Math.floor(diff / WEEK), 'week');
  if (diff < YEAR) return plural(Math.floor(diff / MONTH), 'month');
  return plural(Math.floor(diff / YEAR), 'year');
}

function plural(n, unit) {
  return `${n} ${unit}${n === 1 ? '' : 's'} ago`;
}
