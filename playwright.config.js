import { defineConfig, devices } from '@playwright/test';

/**
 * Playwright configuration for Lite Reader UI tests
 */
export default defineConfig({
  testDir: './tests/ui',
  globalSetup: './tests/ui/global-setup.js',
  
  // Maximum time one test can run
  timeout: 30 * 1000,
  
  // Test execution settings
  fullyParallel: false, // Run tests sequentially to avoid database conflicts
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1, // Single worker to avoid race conditions
  
  // Reporter configuration
  reporter: [
    ['html', { outputFolder: 'reports/playwright' }],
    ['list'],
    process.env.CI ? ['github'] : ['list']
  ],
  
  // Shared settings for all tests
  use: {
    // Base URL for the application
    baseURL: process.env.BASE_URL || 'http://localhost:3000',
    
    // Screenshot and video settings
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    trace: 'retain-on-failure',
    
    // Browser settings
    headless: true,
    viewport: { width: 1280, height: 720 },
  },
  
  // Configure projects for different browsers
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  
  // Web server configuration - start app before tests
  webServer: {
    command: 'make run-test-server',
    url: 'http://localhost:3000/health',
    reuseExistingServer: !process.env.CI,
    timeout: 30 * 1000,
    stdout: 'pipe',
    stderr: 'pipe',
  },
});
