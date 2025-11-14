/**
 * Page Object Model for Main application page (feed reader)
 */

export class MainPage {
  constructor(page) {
    this.page = page;
    
    // Feed management
    this.addFeedButton = page.locator('.add.btn, a.add');
    this.feedUrlInput = page.locator('#urlToAdd');
    this.feedsList = page.locator('#feeds ul');
    this.unreadFeed = page.locator('#unread');
    this.starredFeed = page.locator('#starred');
    
    // Feed actions
    this.updateButton = page.locator('.action.update');
    this.markReadAllButton = page.locator('#mark-read-all');
    this.markUnreadAllButton = page.locator('#mark-unread-all');
    this.removeButton = page.locator('.action.remove');
    this.logoutButton = page.locator('#logout');
    
    // Items
    this.itemsList = page.locator('#items');
    this.items = page.locator('#items li');
    this.message = page.locator('#msg');
  }

  async goto() {
    await this.page.goto('/');
  }

  async isLoggedIn() {
    // Check if we're on the main page and not redirected to login
    const url = this.page.url();
    return !url.includes('login.html');
  }

  async addFeed(feedUrl) {
    // Click the add feed button to show input
    await this.addFeedButton.click();
    
    // Wait for input to be visible
    await this.feedUrlInput.waitFor({ state: 'visible', timeout: 3000 });
    
    // Fill in the URL
    await this.feedUrlInput.fill(feedUrl);
    
    // Wait a moment for the input to be filled
    await this.page.waitForTimeout(300);
    
    // Click the button again to submit (it changes behavior after first click)
    await this.addFeedButton.click();
    
    // Wait for the feed to be added (button changes back to purple)
    await this.page.waitForTimeout(2000);
  }

  async clickFeed(feedTitle) {
    const feedItem = this.page.locator(`.feed:has-text("${feedTitle}")`);
    await feedItem.click();
  }

  async clickUnread() {
    await this.unreadFeed.click();
  }

  async clickStarred() {
    await this.starredFeed.click();
  }

  async getUnreadCount() {
    const countElement = this.unreadFeed.locator('.count span');
    const text = await countElement.textContent();
    return parseInt(text || '0', 10);
  }

  async getStarredCount() {
    const countElement = this.starredFeed.locator('.count span');
    const text = await countElement.textContent();
    return parseInt(text || '0', 10);
  }

  async getItemsCount() {
    return await this.items.count();
  }

  async markItemRead(itemIndex = 0) {
    const item = this.items.nth(itemIndex);
    const readButton = item.locator('a[name="read"]');
    await readButton.click();
  }

  async markItemStarred(itemIndex = 0) {
    const item = this.items.nth(itemIndex);
    const starButton = item.locator('a[name="starred"]');
    await starButton.click();
  }

  async updateFeed() {
    // Use JavaScript to click the update button
    await this.page.evaluate(() => {
      const btn = document.querySelector('.update');
      if (btn) btn.click();
    });
    // Wait for update to complete
    await this.page.waitForTimeout(1000);
  }

  async markAllRead() {
    // Use JavaScript to click the button directly
    await this.page.evaluate(() => {
      const btn = document.querySelector('#mark-read-all');
      if (btn) btn.click();
    });
    await this.page.waitForTimeout(500);
  }

  async markAllUnread() {
    // Use JavaScript to click the button directly  
    await this.page.evaluate(() => {
      const btn = document.querySelector('#mark-unread-all');
      if (btn) btn.click();
    });
    await this.page.waitForTimeout(500);
  }

  async removeFeed() {
    // Set up dialog handler before clicking
    this.page.on('dialog', dialog => dialog.accept());
    
    // Use JavaScript to click the remove button
    await this.page.evaluate(() => {
      const btn = document.querySelector('.remove');
      if (btn) btn.click();
    });
    await this.page.waitForTimeout(500);
  }

  async logout() {
    // The logout button is in the actions bar which is only shown when a feed is selected
    // First, click on "Unread" to ensure actions are visible
    await this.unreadFeed.click();
    await this.page.waitForTimeout(2000);
    
    // Use JavaScript to click the logout button directly, bypassing visibility checks
    // This is more reliable as the button might be in DOM but CSS hidden
    await this.page.evaluate(() => {
      const logoutBtn = document.querySelector('#logout, .logout');
      if (logoutBtn) {
        logoutBtn.click();
      } else {
        throw new Error('Logout button not found in DOM');
      }
    });
    
    // Wait for redirect to login page
    await this.page.waitForURL(/login/, { timeout: 5000 });
  }

  async waitForItems(count, timeout = 5000) {
    await this.page.waitForFunction(
      (expectedCount) => {
        const items = document.querySelectorAll('#items li');
        return items.length >= expectedCount;
      },
      count,  // Pass count as the argument to the function
      { timeout }
    ).catch(() => {
      // Ignore timeout - test will check count afterwards
    });
  }

  async getItemTitle(itemIndex = 0) {
    const item = this.items.nth(itemIndex);
    const title = item.locator('.title span').first();
    return await title.textContent();
  }
}
