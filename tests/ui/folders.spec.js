import { test, expect } from '@playwright/test';

import { SignupPage } from './pages/SignupPage.js';
import { LoginPage } from './pages/LoginPage.js';
import { MainPage } from './pages/MainPage.js';
import { generateEmail, generatePassword, MOCK_FEEDS } from './utils/helpers.js';

let mainPage;

test.beforeEach(async ({ page }) => {
  const email = generateEmail();
  const password = generatePassword();

  const signup = new SignupPage(page);
  await signup.goto();
  await signup.signup(email, password);
  await signup.waitForSuccess();

  const login = new LoginPage(page);
  await login.goto();
  await login.login(email, password);
  await page.waitForURL('http://localhost:3000/', { timeout: 5000 });

  mainPage = new MainPage(page);
  await page.waitForLoadState('networkidle');
});

test.describe('Folders', () => {
  test('create a folder, list it in sidebar', async ({ page }) => {
    await mainPage.addFolder('Tech');
    await expect(mainPage.folderRowByName('Tech')).toBeVisible();
  });

  test('move a feed into a folder via toolbar dropdown', async ({ page }) => {
    await mainPage.addFolder('News');
    const folderId = await mainPage.getFolderId('News');
    expect(folderId).not.toBeNull();

    await mainPage.addFeed(MOCK_FEEDS.techNews);
    await mainPage.clickFeed('Tech');
    await mainPage.moveSelectedFeedToFolder(folderId);

    // After move, the feed should appear inside the folder's child list
    const folderRow = mainPage.folderRowByName('News');
    await expect(folderRow.getByTestId('feed-item')).toHaveCount(1);
  });

  test('clicking the folder shows aggregated items from feeds inside', async ({ page }) => {
    await mainPage.addFolder('Bundle');
    const folderId = await mainPage.getFolderId('Bundle');

    await mainPage.addFeed(MOCK_FEEDS.techNews);
    await mainPage.clickFeed('Tech');
    await mainPage.moveSelectedFeedToFolder(folderId);

    await mainPage.clickFolder('Bundle');
    await expect(mainPage.toolbarTitle).toHaveText('Bundle');
    await mainPage.waitForItems(1);
    expect(await mainPage.getItemsCount()).toBeGreaterThan(0);
  });

  test('collapse/expand folder hides its feeds', async ({ page }) => {
    await mainPage.addFolder('Hide');
    const folderId = await mainPage.getFolderId('Hide');
    await mainPage.addFeed(MOCK_FEEDS.techNews);
    await mainPage.clickFeed('Tech');
    await mainPage.moveSelectedFeedToFolder(folderId);

    const childList = mainPage.folderRowByName('Hide').getByTestId('folder-children').first();
    await expect(childList.getByTestId('feed-item')).toHaveCount(1);

    await mainPage.toggleFolder('Hide');
    // Collapsed list has display:none via CSS class — child feed should not be visible
    await expect(childList).toHaveClass(/collapsed/);
  });

  test('delete folder: feeds revert to root', async ({ page }) => {
    await mainPage.addFolder('Trash');
    const folderId = await mainPage.getFolderId('Trash');
    await mainPage.addFeed(MOCK_FEEDS.techNews);
    await mainPage.clickFeed('Tech');
    await mainPage.moveSelectedFeedToFolder(folderId);

    // Select the folder, then delete via toolbar
    await mainPage.clickFolder('Trash');
    await page.getByTestId('toolbar-remove').click();
    await mainPage.confirmDialog.waitFor({ state: 'visible', timeout: 2000 });
    await mainPage.confirmYes.click();
    await page.waitForTimeout(500);

    await expect(mainPage.folderRowByName('Trash')).toHaveCount(0);
    // Feed still exists at root
    await expect(mainPage.feedItems).toHaveCount(1);
  });
});
