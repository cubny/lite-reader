/**
 * Simple smoke test to verify the testing infrastructure works
 */

import { test, expect } from '@playwright/test';

test('should load the application', async ({ page }) => {
  await page.goto('/login.html');
  
  // Check if the page loads
  await expect(page).toHaveTitle(/Login - Lite Reader/);
  
  // Check if form elements exist
  const emailInput = page.locator('input[name="email"]');
  const passwordInput = page.locator('input[name="password"]');
  
  await expect(emailInput).toBeVisible();
  await expect(passwordInput).toBeVisible();
});
