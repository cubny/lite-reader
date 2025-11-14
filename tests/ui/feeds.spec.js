/**
 * Feed management tests
 */

import { test, expect } from '@playwright/test';
import { SignupPage } from './pages/SignupPage.js';
import { LoginPage } from './pages/LoginPage.js';
import { MainPage } from './pages/MainPage.js';
import { generateUsername, generateEmail, generatePassword, MOCK_FEEDS } from './utils/helpers.js';

test.describe('Feed Management', () => {
  let username;
  let email;
  let password;

  test.beforeEach(async ({ page }) => {
    // Create a new user and login for each test
    username = generateUsername();
    email = generateEmail(username);
    password = generatePassword();

    const signupPage = new SignupPage(page);
    await signupPage.goto();
    await signupPage.signup(username, email, password);
    await signupPage.waitForSuccess();

    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(username, password);
    await page.waitForURL(/index|^\/$/, { timeout: 5000 });
  });

  test('should add RSS feed successfully', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    await mainPage.addFeed(MOCK_FEEDS.techNews);
    
    // Wait for feed to appear in the list
    await page.waitForTimeout(2000);
    
    // Verify feed was added by checking if we can click on it
    const feedItem = page.locator('.feed:has-text("Tech News")');
    await expect(feedItem).toBeVisible({ timeout: 5000 });
  });

  test('should add Atom feed successfully', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    await mainPage.addFeed(MOCK_FEEDS.scienceBlog);
    
    // Wait for feed to appear
    await page.waitForTimeout(2000);
    
    // Verify feed was added
    const feedItem = page.locator('.feed:has-text("Science Blog")');
    await expect(feedItem).toBeVisible({ timeout: 5000 });
  });

  test('should display feed items after adding feed', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    await mainPage.addFeed(MOCK_FEEDS.techNews);
    
    // Click on the feed to view its items
    await page.waitForTimeout(2000);
    await mainPage.clickFeed('Tech News');
    
    // Wait for items to load
    await mainPage.waitForItems(1, 10000);
    
    // Verify items are displayed
    const itemsCount = await mainPage.getItemsCount();
    expect(itemsCount).toBeGreaterThan(0);
  });

  test('should fetch and update feed items', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    await mainPage.addFeed(MOCK_FEEDS.techNews);
    await page.waitForTimeout(2000);
    await mainPage.clickFeed('Tech News');
    
    // Update the feed
    await mainPage.updateFeed();
    
    // Items should still be present
    await mainPage.waitForItems(1, 10000);
    const itemsCount = await mainPage.getItemsCount();
    expect(itemsCount).toBeGreaterThan(0);
  });

  test('should handle empty feed', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    await mainPage.addFeed(MOCK_FEEDS.empty);
    await page.waitForTimeout(2000);
    
    // Click on the empty feed
    await mainPage.clickFeed('Empty Feed');
    
    // Wait a bit for items to load (or not)
    await page.waitForTimeout(2000);
    
    // Should have no items
    const itemsCount = await mainPage.getItemsCount();
    expect(itemsCount).toBe(0);
  });

  test('should remove feed', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    await mainPage.addFeed(MOCK_FEEDS.techNews);
    await page.waitForTimeout(2000);
    
    // Click on the feed
    await mainPage.clickFeed('Tech News');
    await page.waitForTimeout(1000);
    
    // Remove the feed
    await mainPage.removeFeed();
    
    // Wait a bit for removal to complete
    await page.waitForTimeout(1000);
    
    // Feed should no longer be visible
    const feedItem = page.locator('.feed:has-text("Tech News")');
    await expect(feedItem).not.toBeVisible({ timeout: 3000 });
  });
});
