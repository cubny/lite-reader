import { useEffect } from 'preact/hooks';

import { html } from '../util/html.js';
import { getToken } from '../api/client.js';
import { Sidebar } from '../components/Sidebar.js';
import { FeedBar } from '../components/FeedBar.js';
import { ItemList } from '../components/ItemList.js';
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
      <aside class="ui-layout-west" id="feeds" data-testid="pane-sidebar">
        <${Sidebar} />
      </aside>
      <div id="content" class="ui-layout-center">
        <div class="ui-layout-north">
          <${FeedBar} />
        </div>
        <div class="ui-layout-center">
          <${ItemList} />
        </div>
      </div>
    <//>
  `;
}
