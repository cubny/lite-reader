# Testing Guide for Lite Reader

This document provides comprehensive information about testing the Lite Reader application, including both unit tests and automated UI tests.

## Table of Contents

- [Overview](#overview)
- [Unit Tests](#unit-tests)
- [UI Tests](#ui-tests)
  - [Setup](#setup)
  - [Running Tests](#running-tests)
  - [Writing Tests](#writing-tests)
  - [Page Object Models](#page-object-models)
- [Mock Feed System](#mock-feed-system)
  - [How It Works](#how-it-works)
  - [Available Mock Feeds](#available-mock-feeds)
  - [Adding New Mock Feeds](#adding-new-mock-feeds)
- [CI/CD Integration](#cicd-integration)
- [Troubleshooting](#troubleshooting)

## Overview

Lite Reader has two types of tests:

1. **Unit Tests**: Go-based tests for backend logic (services, repositories, API handlers)
2. **UI Tests**: Playwright-based end-to-end tests for the web interface

Both test types are designed to run without internet access, ensuring reliability and speed.

## Unit Tests

### Running Unit Tests

```bash
# Run all unit tests with coverage
make test

# View coverage report
make coverage-report
```

### Unit Test Structure

Unit tests are located alongside the code they test:
- `internal/app/*/service_test.go` - Service layer tests
- `internal/infra/http/api/*_test.go` - API handler tests
- `internal/infra/job/*_test.go` - Background job tests

## UI Tests

### Setup

#### Prerequisites

- Node.js 18+ and npm
- Go 1.23+
- Make

#### First-Time Setup

1. Install Node.js dependencies and Playwright:

```bash
make test-ui-setup
```

This command will:
- Install npm packages (Playwright test framework)
- Download Playwright browsers (Chromium)

### Running Tests

```bash
# Run all UI tests (headless mode)
make test-ui

# Run UI tests with visible browser (helpful for debugging)
make test-ui-headed

# Run specific test file
npm run test:ui tests/ui/auth.spec.js

# Run tests in debug mode (step through tests)
npm run test:ui:debug

# Run both unit and UI tests
make test-all
```

### Test Execution Details

- **Duration**: UI tests typically complete in 2-3 minutes
- **Parallelization**: Tests run sequentially (workers: 1) to avoid database conflicts
- **Browser**: Tests use Chromium in headless mode by default
- **Database**: Each test run uses a separate test database (`data/test-agg.db`)
- **Retries**: In CI, tests retry up to 2 times on failure

### Test Reports

After running UI tests, reports are available at:
- HTML Report: `reports/playwright/index.html`
- Screenshots: `test-results/` (on failure)
- Videos: `test-results/` (on failure)
- Traces: `test-results/` (on failure)

### Test Structure

UI tests are organized as follows:

```
tests/ui/
├── auth.spec.js          # Authentication flow tests
├── feeds.spec.js         # Feed management tests
├── items.spec.js         # Item management tests
├── pages/                # Page Object Models
│   ├── LoginPage.js
│   ├── SignupPage.js
│   └── MainPage.js
└── utils/                # Test utilities
    └── helpers.js        # Helper functions
```

### Writing Tests

#### Basic Test Structure

```javascript
import { test, expect } from '@playwright/test';
import { MainPage } from './pages/MainPage.js';
import { generateUsername, MOCK_FEEDS } from './utils/helpers.js';

test.describe('My Feature', () => {
  test.beforeEach(async ({ page }) => {
    // Setup code (e.g., create user, login)
  });

  test('should do something', async ({ page }) => {
    const mainPage = new MainPage(page);
    
    // Perform actions
    await mainPage.addFeed(MOCK_FEEDS.techNews);
    
    // Assert expectations
    expect(await mainPage.getItemsCount()).toBeGreaterThan(0);
  });
});
```

#### Best Practices

1. **Use Page Objects**: Always use page object models instead of direct selectors
2. **Generate Unique Data**: Use helper functions for usernames/emails to avoid conflicts
3. **Use Mock Feeds**: Always use `MOCK_FEEDS` constants for feed URLs
4. **Add Waits**: Use appropriate waits (`waitForTimeout`, `waitForSelector`) for async operations
5. **Clean State**: Each test should be independent and not rely on previous test state

#### Available Helpers

```javascript
import { 
  generateUsername,     // Generate unique username
  generateEmail,        // Generate unique email
  generatePassword,     // Generate test password
  MOCK_FEEDS,          // Mock feed URLs
} from './utils/helpers.js';
```

### Page Object Models

#### LoginPage

```javascript
const loginPage = new LoginPage(page);
await loginPage.goto();
await loginPage.login(username, password);
```

#### SignupPage

```javascript
const signupPage = new SignupPage(page);
await signupPage.goto();
await signupPage.signup(username, email, password);
await signupPage.waitForSuccess();
```

#### MainPage

```javascript
const mainPage = new MainPage(page);
await mainPage.goto();
await mainPage.addFeed(MOCK_FEEDS.techNews);
await mainPage.clickFeed('Tech News');
await mainPage.markItemRead(0);
await mainPage.logout();
```

See the page object files in `tests/ui/pages/` for all available methods.

## Mock Feed System

### How It Works

The mock feed system provides offline RSS/Atom feeds for testing without requiring internet access.

**Architecture**:
1. **Mock Feed Server**: A Go HTTP server (`internal/testserver/feedserver.go`) that serves XML feeds
2. **Feed Fixtures**: Static XML files in `internal/testserver/fixtures/`
3. **Automatic Startup**: The server starts automatically when running UI tests (via `run-test-server`)

**Ports**:
- Main app: `localhost:3000`
- Mock feed server: `localhost:3001`

### Available Mock Feeds

| Feed Name | URL | Format | Items | Description |
|-----------|-----|--------|-------|-------------|
| Tech News | `http://localhost:3001/feeds/tech-news.xml` | RSS 2.0 | 3 | Technology articles |
| Science Blog | `http://localhost:3001/feeds/science-blog.xml` | Atom 1.0 | 3 | Science articles |
| Empty Feed | `http://localhost:3001/feeds/empty.xml` | RSS 2.0 | 0 | Feed with no items |

Access feeds in tests using:

```javascript
import { MOCK_FEEDS } from './utils/helpers.js';

await mainPage.addFeed(MOCK_FEEDS.techNews);
await mainPage.addFeed(MOCK_FEEDS.scienceBlog);
```

### Adding New Mock Feeds

To add a new mock feed:

1. **Create XML file** in `internal/testserver/fixtures/`:

```bash
# Create a new RSS 2.0 feed
cat > internal/testserver/fixtures/my-feed.xml << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>My Feed</title>
    <link>http://localhost:3001/feeds/my-feed</link>
    <description>My test feed</description>
    <item>
      <title>Item 1</title>
      <link>http://example.com/1</link>
      <description>First item</description>
      <pubDate>Mon, 14 Nov 2025 10:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>
EOF
```

2. **Add to helpers.js**:

```javascript
export const MOCK_FEEDS = {
  techNews: 'http://localhost:3001/feeds/tech-news.xml',
  scienceBlog: 'http://localhost:3001/feeds/science-blog.xml',
  empty: 'http://localhost:3001/feeds/empty.xml',
  myFeed: 'http://localhost:3001/feeds/my-feed.xml',  // Add this
};
```

3. **Use in tests**:

```javascript
await mainPage.addFeed(MOCK_FEEDS.myFeed);
```

**Feed Format Guidelines**:

- Use valid RSS 2.0 or Atom 1.0 format
- Include required fields: `title`, `link`, `description`
- For items, include: `title`, `link`, `pubDate` (RSS) or `published` (Atom)
- Use `localhost:3001` URLs to avoid external dependencies
- Add realistic content for better test coverage

## CI/CD Integration

### GitHub Actions

UI tests are integrated into the CI pipeline. See `.github/workflows/` for configuration.

**Running Tests in CI**:

Tests run automatically on:
- Pull requests
- Pushes to main branch
- Manual workflow dispatch

**CI Environment**:
- Uses `ubuntu-latest` runner
- Installs Go, Node.js, and Playwright
- Sets `CI=true` environment variable
- Generates and uploads test reports as artifacts

### Adding to CI Workflow

If not already configured, add this to your GitHub Actions workflow:

```yaml
- name: Setup Node.js
  uses: actions/setup-node@v4
  with:
    node-version: '18'

- name: Install Playwright
  run: make test-ui-setup

- name: Run UI Tests
  run: make test-ui
  env:
    CI: true

- name: Upload test reports
  if: always()
  uses: actions/upload-artifact@v4
  with:
    name: playwright-report
    path: reports/playwright/
```

## Troubleshooting

### Common Issues

#### Issue: Playwright browsers not installed

**Error**: `Executable doesn't exist at /path/to/chromium`

**Solution**:
```bash
make test-ui-setup
```

#### Issue: Port already in use

**Error**: `listen tcp :3000: bind: address already in use`

**Solution**:
```bash
# Find and kill process using port 3000
lsof -ti:3000 | xargs kill -9

# Or use a different port
HTTP_PORT=3002 make run-test-server
```

#### Issue: Test database locked

**Error**: `database is locked`

**Solution**:
```bash
# Stop all running processes
pkill -f testserver

# Remove test database
rm data/test-agg.db

# Run tests again
make test-ui
```

#### Issue: Tests timing out

**Error**: `Test timeout of 30000ms exceeded`

**Solution**:
- Increase timeout in `playwright.config.js`
- Check if mock feed server is running
- Verify network connectivity between app and mock server

#### Issue: Flaky tests

**Symptoms**: Tests pass sometimes, fail other times

**Solutions**:
1. Add explicit waits: `await page.waitForTimeout(1000)`
2. Wait for elements: `await element.waitFor({ state: 'visible' })`
3. Use retry logic in CI (already configured)
4. Check for race conditions

#### Issue: Screenshots/videos not captured

**Solution**: Check that the `test-results/` directory exists and has write permissions:
```bash
mkdir -p test-results
chmod 755 test-results
```

### Debug Mode

Run tests in debug mode to step through them:

```bash
npm run test:ui:debug
```

This opens the Playwright Inspector where you can:
- Step through test actions
- Inspect page state
- View console logs
- Modify selectors

### Viewing Test Reports

```bash
# Open HTML report
npx playwright show-report reports/playwright

# Or manually open in browser
open reports/playwright/index.html  # macOS
xdg-open reports/playwright/index.html  # Linux
```

### Verbose Logging

Enable verbose logging:

```bash
# Run with debug output
DEBUG=pw:api npm run test:ui

# Or modify playwright.config.js
use: {
  trace: 'on',  // Always capture traces
  video: 'on',  // Always capture video
}
```

### Getting Help

If you encounter issues not covered here:

1. Check the [Playwright documentation](https://playwright.dev/docs/intro)
2. Review test output and error messages
3. Check CI logs for additional details
4. Open an issue on GitHub with:
   - Error message
   - Steps to reproduce
   - Test logs
   - Screenshots/videos if available

## Performance Tips

### Speed Up Tests

1. **Skip browser download in CI**: Browsers are cached automatically
2. **Run specific tests**: `npm run test:ui tests/ui/auth.spec.js`
3. **Disable video recording** for faster execution (edit `playwright.config.js`)
4. **Increase workers** (only if tests are truly independent)

### Test Organization

- Group related tests in `describe` blocks
- Use `beforeEach` for common setup
- Keep tests focused and independent
- Use descriptive test names

## Additional Resources

- [Playwright Documentation](https://playwright.dev)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [Page Object Model Pattern](https://playwright.dev/docs/pom)
- [Debugging Tests](https://playwright.dev/docs/debug)
