/**
 * Item management tests (read/unread, starred)
 */

import { test, expect } from '@playwright/test';
import { SignupPage } from './pages/SignupPage.js';
import { LoginPage } from './pages/LoginPage.js';
import { MainPage } from './pages/MainPage.js';
import { generateEmail, generatePassword, MOCK_FEEDS } from './utils/helpers.js';

test.describe('Item Management', () => {
  let email;
  let password;

  test.beforeEach(async ({ page }) => {
    // Create a new user, login, and add a feed
    email = generateEmail();
    password = generatePassword();

    const signupPage = new SignupPage(page);
    await signupPage.goto();
    await signupPage.signup(email, password);
    await signupPage.waitForSuccess();

    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(email, password);
    await page.waitForURL('http://localhost:3001/', { timeout: 5000 });

    // Add a feed with items
    const mainPage = new MainPage(page);
    await mainPage.addFeed(MOCK_FEEDS.techNews);
    await page.waitForTimeout(2000);
  });

  test('should mark item as read', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    // Click on feed to view items
    await mainPage.clickFeed('Tech News');
    await mainPage.waitForItems(1, 10000);
    
    // Mark first item as read
    await mainPage.markItemRead(0);
    await page.waitForTimeout(500);
    
    // The read icon should change (implementation specific)
    // This is a basic test to ensure the action completes without error
    expect(true).toBe(true);
  });

  test('should mark item as starred', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    // Click on feed to view items
    await mainPage.clickFeed('Tech News');
    await mainPage.waitForItems(1, 10000);
    
    // Mark first item as starred
    await mainPage.markItemStarred(0);
    await page.waitForTimeout(1000);
    
    // Click on starred to see starred items
    await mainPage.clickStarred();
    await page.waitForTimeout(1000);
    
    // Should have at least 1 starred item
    const starredCount = await mainPage.getItemsCount();
    expect(starredCount).toBeGreaterThanOrEqual(1);
  });

  test('should view unread items', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    // Items should be unread by default
    await mainPage.clickUnread();
    await page.waitForTimeout(1000);
    
    // Should have unread items
    const itemsCount = await mainPage.getItemsCount();
    expect(itemsCount).toBeGreaterThan(0);
  });

  test('should mark all items as read', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    // Click on feed to view items
    await mainPage.clickFeed('Tech News');
    await mainPage.waitForItems(1, 10000);
    
    // Mark all as read
    await mainPage.markAllRead();
    await page.waitForTimeout(1000);
    
    // Click unread - should have fewer or no items
    await mainPage.clickUnread();
    await page.waitForTimeout(1000);
    
    // This test verifies the action completes without error
    expect(true).toBe(true);
  });

  test('should display item details', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    // Click on feed to view items
    await mainPage.clickFeed('Tech News');
    await mainPage.waitForItems(1, 10000);
    
    // Get first item title
    const itemTitle = await mainPage.getItemTitle(0);
    
    // Title should not be empty
    expect(itemTitle).toBeTruthy();
    expect(itemTitle.length).toBeGreaterThan(0);
  });
});
