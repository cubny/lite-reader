import { useState } from 'preact/hooks';
import { html } from '../util/html.js';
import { add as addFeed, fetchNew, list as listFeeds } from '../api/feeds.js';
import { feeds } from '../state.js';

function isValidUrl(s) {
  try {
    const u = new URL(s);
    return u.protocol === 'http:' || u.protocol === 'https:';
  } catch {
    return false;
  }
}

export function AddFeedForm() {
  const [url, setUrl] = useState('');
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);

  async function onSubmit(e) {
    e.preventDefault();
    if (!isValidUrl(url.trim())) {
      setError('Please enter a valid URL');
      return;
    }
    setError('');
    setPending(true);
    try {
      const created = await addFeed(url.trim());
      if (created && created.id) {
        try { await fetchNew(created.id); } catch { /* ignore fetch error */ }
      }
      const fresh = await listFeeds();
      feeds.value = fresh || [];
      setUrl('');
    } catch (err) {
      setError(err.message || 'Failed to add feed');
    } finally {
      setPending(false);
    }
  }

  return html`
    <form class="add-feed-form" novalidate onSubmit=${onSubmit} data-testid="add-feed-form">
      <input
        type="url"
        placeholder="http://example.com/feed.xml"
        data-testid="add-feed-url"
        value=${url}
        onInput=${(e) => setUrl(e.target.value)}
      />
      <button type="submit" data-testid="add-feed-submit" disabled=${pending}>Add</button>
      ${error && html`<div class="add-feed-form-error" data-testid="add-feed-error">${error}</div>`}
    </form>
  `;
}
