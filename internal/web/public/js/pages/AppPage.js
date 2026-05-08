import { useEffect } from 'preact/hooks';

import { html } from '../util/html.js';
import { getToken } from '../api/client.js';
import { Resizer } from '../components/Resizer.js';
import { Sidebar } from '../components/Sidebar.js';
import { Toolbar } from '../components/Toolbar.js';
import { ItemList } from '../components/ItemList.js';
import { ItemDetail } from '../components/ItemDetail.js';
import { ErrorBoundary } from '../components/ErrorBoundary.js';

export function AppPage() {
  useEffect(() => {
    if (!getToken()) {
      location.assign('/login.html');
    }
  }, []);

  if (!getToken()) {
    return html`<div data-testid="app-redirecting"></div>`;
  }

  return html`
    <${ErrorBoundary}>
      <div class="app-shell">
        <section class="pane pane-sidebar" data-testid="pane-sidebar">
          <${Sidebar} />
        </section>
        <${Resizer} id="sidebar" cssVar="--sidebar-w" min=${160} max=${500} />
        <section class="pane pane-list" data-testid="pane-list">
          <${Toolbar} />
          <${ItemList} />
        </section>
        <${Resizer} id="list" cssVar="--list-w" min=${240} max=${800} />
        <section class="pane pane-detail" data-testid="pane-detail">
          <${ItemDetail} />
        </section>
      </div>
    <//>
  `;
}
