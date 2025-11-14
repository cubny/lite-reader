/**
 * Page Object Model for Login page
 */

export class LoginPage {
  constructor(page) {
    this.page = page;
    this.emailInput = page.locator('input[name="email"]');
    this.passwordInput = page.locator('input[name="password"]');
    this.loginButton = page.locator('button[type="submit"]');
    this.errorMessage = page.locator('.error, .alert-error');
  }

  async goto() {
    await this.page.goto('/login.html');
  }

  async login(email, password) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.loginButton.click();
  }

  async waitForError() {
    await this.errorMessage.waitFor({ state: 'visible', timeout: 3000 });
  }
}
