/**
 * Page Object Model for the two-pane app shell (data-testid only).
 */
export class MainPage {
  constructor(page) {
    this.page = page;
    this.sidebar = page.getByTestId('sidebar');
    this.feedItems = page.getByTestId('feed-item');
    this.unreadFolder = page.getByTestId('smart-folder-unread');
    this.starredFolder = page.getByTestId('smart-folder-starred');
    this.unreadCountEl = page.getByTestId('smart-folder-unread-count');
    this.starredCountEl = page.getByTestId('smart-folder-starred-count');

    this.addFeedUrl = page.getByTestId('add-feed-url');
    this.addFeedSubmit = page.getByTestId('add-feed-submit');

    this.toolbarTitle = page.getByTestId('toolbar-title');
    this.toolbarRefresh = page.getByTestId('toolbar-refresh');
    this.toolbarMarkRead = page.getByTestId('toolbar-mark-read');
    this.toolbarMarkUnread = page.getByTestId('toolbar-mark-unread');
    this.toolbarRemove = page.getByTestId('toolbar-remove');
    this.toolbarLogout = page.getByTestId('toolbar-logout');

    this.itemList = page.getByTestId('item-list');
    this.items = page.getByTestId('item-row');

    this.confirmDialog = page.getByTestId('confirm-dialog');
    this.confirmYes = page.getByTestId('confirm-yes');
    this.confirmNo = page.getByTestId('confirm-no');
  }

  async goto() {
    await this.page.goto('/');
  }

  async isLoggedIn() {
    return !this.page.url().includes('login.html');
  }

  async addFeed(feedUrl) {
    const before = await this.feedItems.count();
    // First click reveals the input, then we type, then submit via Enter.
    await this.addFeedSubmit.click();
    await this.addFeedUrl.fill(feedUrl);
    await this.addFeedUrl.press('Enter');
    await this.page.waitForFunction(
      (n) => document.querySelectorAll('[data-testid="feed-item"]').length > n,
      before,
      { timeout: 10000 },
    ).catch(() => {});
  }

  feedByText(text) {
    return this.page.getByTestId('feed-item').filter({ hasText: text });
  }

  async clickFeed(feedTitleSubstring) {
    const feed = this.feedByText(feedTitleSubstring).first();
    await feed.waitFor({ state: 'visible', timeout: 5000 });
    await feed.click();
    await this.page.waitForTimeout(500);
  }

  async clickUnread() {
    await this.unreadFolder.click();
    await this.page.waitForTimeout(300);
  }

  async clickStarred() {
    await this.starredFolder.click();
    await this.page.waitForTimeout(300);
  }

  async getUnreadCount() {
    const visible = await this.unreadCountEl.isVisible().catch(() => false);
    if (!visible) return 0;
    const t = (await this.unreadCountEl.textContent()) || '0';
    return parseInt(t, 10) || 0;
  }

  async getStarredCount() {
    const visible = await this.starredCountEl.isVisible().catch(() => false);
    if (!visible) return 0;
    const t = (await this.starredCountEl.textContent()) || '0';
    return parseInt(t, 10) || 0;
  }

  async getItemsCount() {
    return await this.items.count();
  }

  async waitForItems(count, timeout = 10000) {
    await this.page.waitForFunction(
      (expected) => document.querySelectorAll('[data-testid="item-row"]').length >= expected,
      count,
      { timeout },
    ).catch(() => {});
  }

  async markItemRead(itemIndex = 0) {
    const item = this.items.nth(itemIndex);
    await item.getByTestId('item-row-toggle-read').click();
    await this.page.waitForTimeout(200);
  }

  async markItemStarred(itemIndex = 0) {
    const item = this.items.nth(itemIndex);
    await item.getByTestId('item-row-star').click();
    await this.page.waitForTimeout(200);
  }

  async updateFeed() {
    await this.toolbarRefresh.click();
    await this.page.waitForTimeout(800);
  }

  async markAllRead() {
    await this.toolbarMarkRead.click();
    await this.page.waitForTimeout(500);
  }

  async removeFeed(feedTitleSubstring) {
    if (feedTitleSubstring) {
      await this.clickFeed(feedTitleSubstring);
    }
    await this.toolbarRemove.click();
    await this.confirmDialog.waitFor({ state: 'visible', timeout: 2000 });
    await this.confirmYes.click();
    await this.page.waitForTimeout(500);
  }

  async logout() {
    await this.toolbarLogout.click();
    await this.page.waitForURL(/login\.html/, { timeout: 5000 });
  }

  async getItemTitle(itemIndex = 0) {
    const item = this.items.nth(itemIndex);
    return await item.getByTestId('item-row-title').textContent();
  }
}
