import { useMemo } from 'preact/hooks';
import { html } from '../util/html.js';
import { currentItem } from '../state.js';
import { detectDir } from '../util/dom.js';

function buildSrcdoc(content, dir) {
  const safeDir = dir === 'rtl' ? 'rtl' : 'ltr';
  return `<!doctype html><html dir="${safeDir}"><head><meta charset="utf-8"><base target="_blank"><style>body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;color:#2a2a28;line-height:1.5;padding:0;margin:0;}img,video,iframe{max-width:100%;height:auto;}a{color:#6c8a46;}</style></head><body>${content || ''}</body></html>`;
}

export function ItemDetail() {
  const item = currentItem.value;

  const srcdoc = useMemo(() => {
    if (!item) return '';
    const dir = detectDir((item.title || '') + ' ' + (item.desc || ''));
    return buildSrcdoc(item.desc || '', dir);
  }, [item && item.id]);

  if (!item) {
    return html`<div class="item-detail-empty" data-testid="item-detail-empty">Select an item</div>`;
  }

  return html`
    <article class="item-detail" data-testid="item-detail">
      <header class="item-detail-header">
        <h1 data-testid="item-detail-title">${item.title || '(no title)'}</h1>
        ${item.link && html`
          <a
            class="item-detail-source"
            data-testid="item-detail-link"
            href=${item.link}
            target="_blank"
            rel="noopener noreferrer"
          >${item.link}</a>
        `}
      </header>
      <iframe
        sandbox=""
        srcdoc=${srcdoc}
        data-testid="item-detail-iframe"
        title=${item.title || 'Item content'}
      ></iframe>
    </article>
  `;
}
