/**
 * Page Object Model for Signup page
 */

export class SignupPage {
  constructor(page) {
    this.page = page;
    this.emailInput = page.locator('input[name="email"]');
    this.passwordInput = page.locator('input[name="password"]');
    this.confirmPasswordInput = page.locator('input[name="confirm-password"]');
    this.signupButton = page.locator('button[type="submit"]');
    this.errorMessage = page.locator('.error, .alert-error');
    this.successMessage = page.locator('.success, .alert-success');
  }

  async goto() {
    await this.page.goto('/signup.html');
  }

  async signup(email, password) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.confirmPasswordInput.fill(password);
    await this.signupButton.click();
  }

  async waitForSuccess() {
    // After signup, might redirect to login or show success message
    await this.page.waitForURL(/login|index/, { timeout: 5000 }).catch(() => {
      // If no redirect, check for success message
      return this.successMessage.waitFor({ state: 'visible', timeout: 3000 });
    });
  }

  async waitForError() {
    await this.errorMessage.waitFor({ state: 'visible', timeout: 3000 });
  }
}
