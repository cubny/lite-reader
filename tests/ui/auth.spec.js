/**
 * Authentication flow tests
 */

import { test, expect } from '@playwright/test';
import { SignupPage } from './pages/SignupPage.js';
import { LoginPage } from './pages/LoginPage.js';
import { MainPage } from './pages/MainPage.js';
import { generateEmail, generatePassword } from './utils/helpers.js';

test.describe('Authentication', () => {
  let email;
  let password;

  test.beforeEach(() => {
    email = generateEmail();
    password = generatePassword();
  });

  test('should signup with valid credentials', async ({ page }) => {
    const signupPage = new SignupPage(page);
    
    await signupPage.goto();
    await signupPage.signup(email, password);
    
    // Should redirect to login page or show success
    await signupPage.waitForSuccess();
  });

  test('should login with valid credentials', async ({ page }) => {
    // First, create a user
    const signupPage = new SignupPage(page);
    await signupPage.goto();
    await signupPage.signup(email, password);
    await signupPage.waitForSuccess();

    // Now login
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(email, password);

    // Should be redirected to main page (root path)
    await page.waitForURL('http://localhost:3000/', { timeout: 5000 });
    
    const mainPage = new MainPage(page);
    expect(await mainPage.isLoggedIn()).toBe(true);
  });

  test('should not login with invalid credentials', async ({ page }) => {
    const loginPage = new LoginPage(page);
    
    await loginPage.goto();
    await loginPage.login('nonexistent@test.example', 'wrongpassword');
    
    // Should show error message and stay on login page
    await page.waitForTimeout(1000);
    const url = page.url();
    expect(url).toContain('login');
    
    // Verify error message is displayed
    const errorMessage = await page.locator('.error-message');
    await expect(errorMessage).toBeVisible();
    await expect(errorMessage).toContainText('Invalid email or password');
  });

  test('should logout successfully', async ({ page }) => {
    // First, signup and login
    const signupPage = new SignupPage(page);
    await signupPage.goto();
    await signupPage.signup(email, password);
    await signupPage.waitForSuccess();

    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(email, password);
    await page.waitForURL('http://localhost:3000/', { timeout: 5000 });

    // Now logout
    const mainPage = new MainPage(page);
    await mainPage.logout();
    
    // Should be redirected to login page
    const url = page.url();
    expect(url).toContain('login');
  });

  test('should not signup with duplicate email', async ({ page }) => {
    // First, create a user
    const signupPage = new SignupPage(page);
    await signupPage.goto();
    await signupPage.signup(email, password);
    await signupPage.waitForSuccess();

    // Try to signup again with the same email
    await signupPage.goto();
    await signupPage.signup(email, password);
    
    // Should show error message and stay on signup page
    await page.waitForTimeout(1000);
    const url = page.url();
    expect(url).toContain('signup');
    
    // Verify error message is displayed
    const errorMessage = await page.locator('.error-message');
    await expect(errorMessage).toBeVisible();
    await expect(errorMessage).toContainText('Email may already be registered');
  });
});
