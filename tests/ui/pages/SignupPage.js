/**
 * Page Object Model for Signup page (data-testid only).
 */
export class SignupPage {
  constructor(page) {
    this.page = page;
    this.emailInput = page.getByTestId('signup-email');
    this.passwordInput = page.getByTestId('signup-password');
    this.confirmPasswordInput = page.getByTestId('signup-confirm-password');
    this.signupButton = page.getByTestId('signup-submit');
    this.errorMessage = page.getByTestId('signup-error');
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
    await this.page.waitForURL(/login\.html/, { timeout: 5000 });
  }

  async waitForError() {
    await this.errorMessage.waitFor({ state: 'visible', timeout: 3000 });
    return this.errorMessage.textContent();
  }
}
