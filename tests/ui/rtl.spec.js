/**
 * RTL (Farsi) rendering — exercises the right-to-left CSS path:
 *   - Compact list row puts the timestamp in left-to-right reading order
 *   - Expanded article keeps eyebrow / action chips in LTR order
 *   - Title and body inherit `direction: rtl`
 */

import { test, expect } from '@playwright/test';
import { SignupPage } from './pages/SignupPage.js';
import { LoginPage } from './pages/LoginPage.js';
import { MainPage } from './pages/MainPage.js';
import { generateEmail, generatePassword, MOCK_FEEDS } from './utils/helpers.js';

test.describe('RTL feeds (Farsi)', () => {
  let email;
  let password;

  test.beforeEach(async ({ page }) => {
    email = generateEmail();
    password = generatePassword();

    const signupPage = new SignupPage(page);
    await signupPage.goto();
    await signupPage.signup(email, password);
    await signupPage.waitForSuccess();

    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(email, password);
    await page.waitForURL('http://localhost:3000/', { timeout: 5000 });

    const mainPage = new MainPage(page);
    await mainPage.addFeed(MOCK_FEEDS.farsi);
    await page.waitForTimeout(2000);
  });

  test('compact row gets the .rtl class on a Farsi item', async ({ page }) => {
    const mainPage = new MainPage(page);
    await mainPage.clickFeed('فارسی');
    await page.waitForTimeout(800);
    const firstRow = page.locator('[data-testid="item-row"]').first();
    await expect(firstRow).toHaveClass(/\brtl\b/);
  });

  test('expanded article keeps eyebrow + action chips in LTR direction', async ({ page }) => {
    const mainPage = new MainPage(page);
    await mainPage.clickFeed('فارسی');
    await page.waitForTimeout(800);
    const firstTitle = page.locator('[data-testid="item-row-title"]').first();
    await firstTitle.click();
    await page.waitForTimeout(400);

    const eyebrow = page.locator('.lr-article-eyebrow').first();
    const actions = page.locator('.lr-article-actions').first();
    const title = page.locator('[data-testid="item-row-title"]').first();
    const body = page.locator('[data-testid="item-row-body"]').first();

    await expect(eyebrow).toHaveCSS('direction', 'ltr');
    await expect(actions).toHaveCSS('direction', 'ltr');
    await expect(title).toHaveCSS('direction', 'rtl');
    await expect(body).toHaveCSS('direction', 'rtl');
  });
});
