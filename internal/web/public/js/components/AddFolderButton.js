import { useEffect, useRef, useState } from 'preact/hooks';
import { html } from '../util/html.js';
import { create as createFolder, list as listFolders } from '../api/folders.js';
import { folders } from '../state.js';

export function AddFolderButton() {
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
    const name = (inputRef.current && inputRef.current.value || '').trim();
    if (!name) { reset(); return; }
    setPending(true);
    try {
      await createFolder(name);
      const fresh = await listFolders();
      folders.value = fresh || [];
      reset();
    } catch (err) {
      setError(err.message || 'Failed to add folder');
    } finally {
      setPending(false);
    }
  }

  function onButtonClick(e) {
    e.preventDefault();
    e.stopPropagation();
    if (!open) { reveal(e); return; }
    submit();
  }

  function onKey(e) {
    if (e.key === 'Enter') { e.preventDefault(); submit(); }
    else if (e.key === 'Escape') { reset(); }
  }

  const btnClass = open ? 'add btn btn-green' : 'add btn';
  const iconClass = pending ? 'icon-spin icon-spinner' : (open ? 'icon-folder-open' : 'icon-folder-close');
  const containerClass = open ? 'open' : '';

  return html`
    <div id="addfolder" class=${containerClass} ref=${containerRef} data-testid="add-folder-form">
      <a class=${btnClass} href="#" onClick=${onButtonClick} data-testid="add-folder-submit">
        <i class=${iconClass}></i> <span>${open ? '' : 'New folder'}</span>
      </a>
      <input
        ref=${inputRef}
        type="text"
        data-testid="add-folder-name"
        style=${open ? '' : 'display: none'}
        onKeyDown=${onKey}
        disabled=${pending}
      />
      ${error && html`<div class="add-feed-form-error" data-testid="add-folder-error">${error}</div>`}
    </div>
  `;
}
