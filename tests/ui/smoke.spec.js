/**
 * Simple smoke test to verify the testing infrastructure works
 */

import { test, expect } from '@playwright/test';

test('should load the application', async ({ page }) => {
  await page.goto('/login.html');

  await expect(page).toHaveTitle(/Login - Lite Reader/);

  await expect(page.getByTestId('login-email')).toBeVisible();
  await expect(page.getByTestId('login-password')).toBeVisible();
  await expect(page.getByTestId('login-submit')).toBeVisible();
});
