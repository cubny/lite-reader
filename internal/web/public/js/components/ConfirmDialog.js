import { useEffect, useRef } from 'preact/hooks';
import { html } from '../util/html.js';

export function ConfirmDialog({ message, confirmLabel = 'OK', cancelLabel = 'Cancel', onConfirm, onCancel }) {
  const yesRef = useRef(null);

  useEffect(() => {
    yesRef.current && yesRef.current.focus();
    function onKey(e) {
      if (e.key === 'Escape') {
        e.preventDefault();
        onCancel();
        return;
      }
      // Don't handle Enter globally — when focus is on the Yes button the
      // browser already fires its native click on Enter, and a global Enter
      // handler would cause onConfirm to run twice (double-delete).
    }
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [onConfirm, onCancel]);

  return html`
    <div class="confirm-backdrop" data-testid="confirm-dialog" role="dialog" aria-modal="true">
      <div class="confirm-dialog">
        <div data-testid="confirm-message">${message}</div>
        <div class="confirm-actions">
          <button type="button" data-testid="confirm-no" onClick=${onCancel}>${cancelLabel}</button>
          <button type="button" ref=${yesRef} data-testid="confirm-yes" onClick=${onConfirm}>${confirmLabel}</button>
        </div>
      </div>
    </div>
  `;
}
