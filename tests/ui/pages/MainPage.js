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
    // Ensure update button is visible
    await this.updateButton.waitFor({ state: 'visible', timeout: 5000 });
    await this.updateButton.click();
    // Wait for update to complete
    await this.page.waitForTimeout(1000);
  }

  async markAllRead() {
    // Ensure button is visible (should be visible when a feed is selected)
    await this.markReadAllButton.waitFor({ state: 'visible', timeout: 3000 });
    await this.markReadAllButton.click();
    await this.page.waitForTimeout(500);
  }

  async markAllUnread() {
    // Ensure button is visible (should be visible when a feed is selected)
    await this.markUnreadAllButton.waitFor({ state: 'visible', timeout: 3000 });
    await this.markUnreadAllButton.click();
    await this.page.waitForTimeout(500);
  }

  async removeFeed() {
    // Ensure remove button is visible
    await this.removeButton.waitFor({ state: 'visible', timeout: 5000 });
    
    // Set up dialog handler before clicking
    this.page.on('dialog', dialog => dialog.accept());
    
    await this.removeButton.click();
    await this.page.waitForTimeout(500);
  }

  async logout() {
    // Click on "Unread" to show the action buttons (including logout)
    // The action buttons are only visible when a feed is selected
    await this.unreadFeed.click();
    
    // Wait for the actions to become visible
    await this.page.waitForTimeout(2000);
    
    // Try to force click if the button is in DOM but not visible due to CSS
    try {
      await this.logoutButton.click({ force: true, timeout: 5000 });
    } catch (error) {
      // If force click doesn't work, try waiting for visibility
      await this.logoutButton.waitFor({ state: 'visible', timeout: 5000 });
      await this.logoutButton.click();
    }
    
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
