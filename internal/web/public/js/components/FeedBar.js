import { useState } from 'preact/hooks';
import { html } from '../util/html.js';
import { selection, feeds, folders, items } from '../state.js';
import {
  fetchNew as fetchFeed,
  markRead,
  markUnread,
  remove as removeFeed,
  items as feedItems,
  list as listFeeds,
  move as moveFeed,
} from '../api/feeds.js';
import {
  items as folderItemsApi,
  markRead as markFolderRead,
  remove as removeFolder,
  rename as renameFolder,
  list as listFolders,
} from '../api/folders.js';
import { unread as unreadItems, starred as starredItems } from '../api/items.js';
import { logout as authLogout } from '../api/auth.js';
import { refreshCounts } from './SmartFolders.js';
import { ConfirmDialog } from './ConfirmDialog.js';
import { navigate } from '../router.js';

function titleFor(sel, feedList, folderList) {
  if (!sel) return '';
  if (sel.kind === 'unread') return 'Unread';
  if (sel.kind === 'starred') return 'Starred';
  if (sel.kind === 'feed') {
    const f = (feedList || []).find((x) => x.id === sel.id);
    return f ? (f.title || f.url) : '';
  }
  if (sel.kind === 'folder') {
    const fo = (folderList || []).find((x) => x.id === sel.id);
    return fo ? fo.name : '';
  }
  return '';
}

async function reloadCurrent(sel) {
  if (!sel) return;
  if (sel.kind === 'unread') {
    items.value = (await unreadItems()) || [];
  } else if (sel.kind === 'starred') {
    items.value = (await starredItems()) || [];
  } else if (sel.kind === 'feed') {
    items.value = (await feedItems(sel.id)) || [];
  } else if (sel.kind === 'folder') {
    items.value = (await folderItemsApi(sel.id)) || [];
  }
}

export function FeedBar() {
  const [pending, setPending] = useState(false);
  const [confirming, setConfirming] = useState(false);
  const [confirmingFolderDelete, setConfirmingFolderDelete] = useState(false);
  const [renamingFolder, setRenamingFolder] = useState(false);
  const [pendingFolderName, setPendingFolderName] = useState('');
  const sel = selection.value;
  const folderList = folders.value || [];
  const title = titleFor(sel, feeds.value, folderList);
  const isFeedScope = sel.kind === 'feed';
  const isFolderScope = sel.kind === 'folder';
  const isSmartScope = sel.kind === 'unread' || sel.kind === 'starred';
  const barClass = isFeedScope ? 'has-feed' : (isFolderScope ? 'has-folder' : (isSmartScope ? 'has-smart' : ''));

  const currentFeed = isFeedScope ? (feeds.value || []).find((x) => x.id === sel.id) : null;
  const currentFolderId = currentFeed ? (currentFeed.folder_id == null ? '' : String(currentFeed.folder_id)) : '';

  async function refresh() {
    if (!isFeedScope || pending) return;
    setPending(true);
    try {
      await fetchFeed(sel.id);
      await reloadCurrent(sel);
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  async function readAll() {
    if (pending) return;
    if (!isFeedScope && !isFolderScope) return;
    setPending(true);
    try {
      if (isFeedScope) await markRead(sel.id);
      else await markFolderRead(sel.id);
      await reloadCurrent(sel);
      await refreshCounts();
      const fresh = await listFolders();
      folders.value = fresh || [];
    } finally {
      setPending(false);
    }
  }

  async function unreadAll() {
    if (!isFeedScope || pending) return;
    setPending(true);
    try {
      await markUnread(sel.id);
      await reloadCurrent(sel);
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  function askRemove() {
    if (!isFeedScope) return;
    setConfirming(true);
  }

  async function confirmRemove() {
    setConfirming(false);
    if (!isFeedScope) return;
    setPending(true);
    try {
      await removeFeed(sel.id);
      const fresh = await listFeeds();
      feeds.value = fresh || [];
      items.value = [];
      selection.value = { kind: 'unread' };
      navigate('#/');
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  async function onMoveChange(e) {
    if (!isFeedScope) return;
    const v = e.target.value;
    const folderId = v === '' ? null : Number(v);
    setPending(true);
    try {
      await moveFeed(sel.id, folderId);
      const fresh = await listFeeds();
      feeds.value = fresh || [];
      const freshFolders = await listFolders();
      folders.value = freshFolders || [];
    } finally {
      setPending(false);
    }
  }

  async function commitFolderRename(e) {
    e && e.preventDefault();
    const trimmed = (pendingFolderName || '').trim();
    setRenamingFolder(false);
    if (!isFolderScope || !trimmed) return;
    try {
      await renameFolder(sel.id, trimmed);
      const fresh = await listFolders();
      folders.value = fresh || [];
    } catch {
      // toast handled elsewhere
    }
  }

  async function confirmFolderDelete() {
    setConfirmingFolderDelete(false);
    if (!isFolderScope) return;
    setPending(true);
    try {
      await removeFolder(sel.id);
      const [freshFolders, freshFeeds] = await Promise.all([listFolders(), listFeeds()]);
      folders.value = freshFolders || [];
      feeds.value = freshFeeds || [];
      items.value = [];
      selection.value = { kind: 'unread' };
      navigate('#/');
      await refreshCounts();
    } finally {
      setPending(false);
    }
  }

  return html`
    <div id="feedbar" class=${barClass}>
      <div id="title" data-testid="toolbar-title">
        ${isFolderScope && renamingFolder
          ? html`<input
              type="text"
              class="folder-rename"
              data-testid="toolbar-folder-rename-input"
              value=${pendingFolderName}
              autofocus
              onInput=${(e) => setPendingFolderName(e.target.value)}
              onBlur=${commitFolderRename}
              onKeyDown=${(e) => { if (e.key === 'Enter') commitFolderRename(e); else if (e.key === 'Escape') setRenamingFolder(false); }}
            />`
          : title}
      </div>
      <div id="actions">
        ${isFeedScope && html`
          <div class="action move-to-folder" data-testid="toolbar-move-to-folder">
            <select
              data-testid="toolbar-move-select"
              value=${currentFolderId}
              onChange=${onMoveChange}
              disabled=${pending}
              title="Move feed to folder"
            >
              <option value="">No folder</option>
              ${folderList.map((f) => html`<option value=${String(f.id)}>${f.name}</option>`)}
            </select>
          </div>
        `}
        <div class="action update" data-testid="toolbar-refresh" onClick=${refresh}>
          <i class="icon-repeat"></i> Update
        </div>
        <div id="mark-read-all" class="action markread" data-testid="toolbar-mark-read" onClick=${readAll}>
          <i class="icon-circle-blank"></i> Read All
        </div>
        <div id="mark-unread-all" class="action markunread" data-testid="toolbar-mark-unread" onClick=${unreadAll}>
          <i class="icon-circle"></i> Unread All
        </div>
        ${isFolderScope && html`
          <div class="action rename-folder" data-testid="toolbar-folder-rename" onClick=${() => { setPendingFolderName(title); setRenamingFolder(true); }}>
            <i class="icon-edit"></i> Rename
          </div>
        `}
        <div class="action logout" id="logout" data-testid="toolbar-logout" onClick=${authLogout}>
          <i class="icon-signout"></i> Logout
        </div>
        <div class="action remove" data-testid="toolbar-remove" onClick=${isFolderScope ? () => setConfirmingFolderDelete(true) : askRemove}>
          <i class="icon-trash"></i> ${isFolderScope ? 'Delete' : 'Remove'}
        </div>
      </div>
      ${confirming && html`
        <${ConfirmDialog}
          message=${`Delete "${title}"?`}
          onConfirm=${confirmRemove}
          onCancel=${() => setConfirming(false)}
        />
      `}
      ${confirmingFolderDelete && html`
        <${ConfirmDialog}
          message=${`Delete folder "${title}"? Feeds inside will move out of the folder.`}
          onConfirm=${confirmFolderDelete}
          onCancel=${() => setConfirmingFolderDelete(false)}
        />
      `}
    </div>
  `;
}
