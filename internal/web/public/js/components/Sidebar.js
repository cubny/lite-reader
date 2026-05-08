import { useEffect } from 'preact/hooks';
import { html } from '../util/html.js';
import { feeds } from '../state.js';
import { list as listFeeds } from '../api/feeds.js';
import { FeedList } from './FeedList.js';
import { AddFeedForm } from './AddFeedForm.js';
import { SmartFolders } from './SmartFolders.js';

export function Sidebar() {
  useEffect(() => {
    listFeeds().then((list) => { feeds.value = list || []; }).catch(() => {});
  }, []);

  return html`
    <div data-testid="sidebar">
      <div id="toolbar" class="ui-layout-north">
        <${AddFeedForm} />
      </div>
      <ul>
        <${SmartFolders} />
        <${FeedList} />
      </ul>
    </div>
  `;
}
