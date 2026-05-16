import { useEffect, useRef, useState } from 'preact/hooks';
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
  const [open, setOpen] = useState(false);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState('');
  const inputRef = useRef(null);
  const containerRef = useRef(null);

  useEffect(() => {
    if (!open) return undefined;
    function onDocMouseDown(e) {
      const root = containerRef.current;
      if (root && !root.contains(e.target)) {
        setOpen(false);
        setError('');
        if (inputRef.current) inputRef.current.value = '';
      }
    }
    document.addEventListener('mousedown', onDocMouseDown);
    return () => document.removeEventListener('mousedown', onDocMouseDown);
  }, [open]);

  function reveal(e) {
    e.preventDefault();
    e.stopPropagation();
    setError('');
    setOpen(true);
    setTimeout(() => inputRef.current && inputRef.current.focus(), 0);
  }

  function reset() {
    setOpen(false);
    setError('');
    if (inputRef.current) inputRef.current.value = '';
  }

  async function submit() {
    const raw = (inputRef.current && inputRef.current.value || '').trim();
    if (!raw) {
      reset();
      return;
    }
    const url = raw.startsWith('http') ? raw : `http://${raw}`;
    if (!isValidUrl(url)) {
      setError('Please enter a valid URL');
      return;
    }
    setPending(true);
    try {
      const created = await addFeed(url);
      if (created && created.id) {
        try { await fetchNew(created.id); } catch { /* ignore */ }
      }
      const fresh = await listFeeds();
      feeds.value = fresh || [];
      reset();
    } catch (err) {
      setError(err.message || 'Failed to add feed');
    } finally {
      setPending(false);
    }
  }

  function onButtonClick(e) {
    e.preventDefault();
    e.stopPropagation();
    if (!open) {
      reveal(e);
      return;
    }
    submit();
  }

  function onKey(e) {
    if (e.key === 'Enter') {
      e.preventDefault();
      submit();
    } else if (e.key === 'Escape') {
      reset();
    }
  }

  const btnClass = open ? 'add btn btn-green' : 'add btn btn-purple';
  const submitIcon = pending ? 'icon-spin icon-spinner' : (open ? 'icon-ok' : 'icon-plus');

  return html`
    <div id="addfeed" ref=${containerRef} data-testid="add-feed-form">
      ${open && html`<i class="addfeed-lead icon-plus" aria-hidden="true"></i>`}
      <input
        ref=${inputRef}
        type="text"
        id="urlToAdd"
        data-testid="add-feed-url"
        placeholder="https://"
        style=${open ? '' : 'display: none'}
        onKeyDown=${onKey}
        disabled=${pending}
      />
      <a class=${btnClass} href="#" onClick=${onButtonClick} data-testid="add-feed-submit">
        <i class=${submitIcon}></i> <span>${open ? '' : 'Add Feed'}</span>
      </a>
      ${error && html`<div class="add-feed-form-error" data-testid="add-feed-error">${error}</div>`}
    </div>
  `;
}
