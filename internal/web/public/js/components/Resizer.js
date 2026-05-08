import { useEffect, useRef } from 'preact/hooks';
import { html } from '../util/html.js';

const KEY_PREFIX = 'resizer:';

function clamp(n, lo, hi) {
  return Math.max(lo, Math.min(hi, n));
}

export function Resizer({ id, cssVar, min = 120, max = 600, root = null, testId }) {
  const handleRef = useRef(null);

  useEffect(() => {
    const r = root || document.documentElement;
    const stored = Number(localStorage.getItem(KEY_PREFIX + id));
    if (Number.isFinite(stored) && stored > 0) {
      r.style.setProperty(cssVar, `${clamp(stored, min, max)}px`);
    }
  }, [id, cssVar, min, max, root]);

  function onPointerDown(e) {
    e.preventDefault();
    const r = root || document.documentElement;
    const startX = e.clientX;
    const current = parseFloat(getComputedStyle(r).getPropertyValue(cssVar)) || min;

    function onMove(ev) {
      const next = clamp(current + (ev.clientX - startX), min, max);
      r.style.setProperty(cssVar, `${next}px`);
    }
    function onUp() {
      document.removeEventListener('mousemove', onMove);
      document.removeEventListener('mouseup', onUp);
      const raw = parseFloat(getComputedStyle(r).getPropertyValue(cssVar))
        || parseFloat(r.style.getPropertyValue(cssVar))
        || min;
      localStorage.setItem(KEY_PREFIX + id, String(raw));
    }
    document.addEventListener('mousemove', onMove);
    document.addEventListener('mouseup', onUp);
  }

  return html`
    <div
      ref=${handleRef}
      class="resizer"
      data-testid=${testId || `resizer-${id}`}
      role="separator"
      aria-orientation="vertical"
      onMouseDown=${onPointerDown}
    ></div>
  `;
}
