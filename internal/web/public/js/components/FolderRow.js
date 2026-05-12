import { html } from '../util/html.js';
import {
  selection,
  folderCollapsed,
  persistCollapsed,
} from '../state.js';
import { navigate } from '../router.js';
import { FeedItem } from './FeedItem.js';

export function FolderRow({ folder, childFeeds }) {
  const sel = selection.value;
  const collapsedMap = folderCollapsed.value;
  const isCollapsed = !!collapsedMap[folder.id];
  const isSelected = sel.kind === 'folder' && sel.id === folder.id;
  const hasChildren = (childFeeds || []).length > 0;
  const unread = folder.unread_count || 0;

  function toggleCollapse(e) {
    e.stopPropagation();
    if (!hasChildren) return;
    const next = { ...collapsedMap, [folder.id]: !isCollapsed };
    folderCollapsed.value = next;
    persistCollapsed(next);
  }

  function selectFolder() {
    selection.value = { kind: 'folder', id: folder.id };
    navigate(`#/folder/${folder.id}`);
  }

  const chevronIcon = isCollapsed ? 'icon-chevron-right' : 'icon-chevron-down';
  // Children list is always rendered so SortableJS keeps it as a drop target
  // even when visually collapsed.
  const collapsedClass = isCollapsed && hasChildren ? ' collapsed' : '';

  return html`
    <li class="folder-wrap" data-testid="folder-row" data-folder-id=${folder.id}>
      <div
        class=${`folder-header${isSelected ? ' selected' : ''}`}
        data-testid="folder-header"
        onClick=${selectFolder}
      >
        ${hasChildren
          ? html`<button
              type="button"
              class="folder-chevron"
              aria-label=${isCollapsed ? 'Expand folder' : 'Collapse folder'}
              data-testid="folder-toggle"
              onClick=${toggleCollapse}
            ><i class=${chevronIcon}></i></button>`
          : html`<span class="folder-chevron-spacer"></span>`}
        <i class="icon-folder-close folder-icon"></i>
        <div class="folder-title" data-testid="folder-title">${folder.name}</div>
        <div class="folder-count">${unread > 0 ? html`<span data-testid="folder-unread-count">${unread}</span>` : ''}</div>
      </div>
      <ul
        class=${`folder-children${collapsedClass}`}
        data-testid="folder-children"
        data-folder-id=${folder.id}
      >
        ${(childFeeds || []).map((f) => html`
          <${FeedItem}
            key=${f.id}
            feed=${f}
            isSelected=${sel.kind === 'feed' && sel.id === f.id}
          />
        `)}
      </ul>
    </li>
  `;
}
