/**
 * Article scraper tests — "Load full article" button on item rows.
 */

import { test, expect } from '@playwright/test';
import { SignupPage } from './pages/SignupPage.js';
import { LoginPage } from './pages/LoginPage.js';
import { MainPage } from './pages/MainPage.js';
import { generateEmail, generatePassword, MOCK_FEEDS } from './utils/helpers.js';

test.describe('Article Scraper', () => {
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
    await mainPage.addFeed(MOCK_FEEDS.techNews);
    await page.waitForTimeout(2000);
  });

  test('loads full article into the item body on click', async ({ page }) => {
    const mainPage = new MainPage(page);

    await mainPage.clickFeed('Tech News');
    await mainPage.waitForItems(1, 10000);

    // Expand the first item so its body and inline actions become visible.
    await mainPage.items.nth(0).getByTestId('item-row-title').click();
    await expect(mainPage.itemBody(0)).toBeVisible();

    // The fixture article with a known unique sentence in the full body.
    const fullArticleSentence = 'Edge runtimes change the latency calculus';

    // Body must NOT contain the full sentence before scraping.
    await expect(mainPage.itemBody(0)).not.toContainText(fullArticleSentence);

    await mainPage.loadFullArticle(0);

    // After scraping, the body should include extracted content. Allow up to
    // 15s for the synchronous fetch + readability extraction.
    await expect(mainPage.itemBody(0)).toContainText(fullArticleSentence, { timeout: 15000 });
  });
});
