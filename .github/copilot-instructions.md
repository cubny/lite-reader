# GitHub Copilot Instructions for Lite Reader

## Project Overview

Lite Reader is a lightweight RSS feed aggregator written in Go. It allows users to read feeds on their own machine with a simple and minimal application. The application supports multi-user functionality, feed management, and item management.

## Technology Stack

- **Language:** Go 1.23.0
- **Database:** SQLite (using modernc.org/sqlite)
- **HTTP Router:** julienschmidt/httprouter
- **Testing:** testify/assert
- **Mocking:** uber/mock (go.uber.org/mock)
- **Build System:** Makefile
- **Dependency Management:** Go modules with vendoring

## Architecture

The project follows a clean architecture pattern:

- `cmd/` - Application entry points
- `internal/app/` - Application business logic layer
  - `auth/` - Authentication and user management
  - `feed/` - Feed management logic
  - `item/` - Item management logic
- `internal/infra/` - Infrastructure layer
  - `http/api/` - HTTP API handlers and middleware
  - `sqlite/` - SQLite repository implementations
  - `job/` - Background job processing
- `internal/config/` - Configuration management
- `internal/mocks/` - Mock implementations for testing
- `public/` - Static assets
- `vendor/` - Vendored dependencies

## Development Workflow

### Running the Application

- `make run` - Run the app without Docker
- `make docker-run` - Run the app with Docker Compose
- `make build` - Build the binary
- `make run-test-server` - Run the app with mock feed server for testing (uses test database)

### Testing

- `make test` - Run all unit tests with coverage
- `make test-ui` - Run UI tests with Playwright (headless)
- `make test-ui-headed` - Run UI tests with visible browser (for debugging)
- `make test-all` - Run both unit and UI tests
- Coverage reports are generated in `reports/` directory
- Aim for high test coverage (current: ~36.3% total, with 100% in feed and item app logic)

#### UI Testing Setup

Before running UI tests for the first time:
- `make test-ui-setup` - Install Playwright and dependencies (one-time setup)

### Code Quality

- `make lint` - Run golangci-lint for code quality checks
- `make gci` - Format imports using gci
- `make pre-commit` - Run all essential checks (gomod, update-mocks, lint, test)
- **Note**: `make pre-commit` runs unit tests but not UI tests. Run `make test-ui` separately before finalizing PRs

### Dependency Management

- `make gomod` - Tidy go.mod and update vendor directory
- Always run `make gomod` after adding or updating dependencies
- Use vendored dependencies (`-mod=vendor` flag)

### Mocks

- `make update-mocks` - Regenerate mocks using mockgen
- Mock definitions are in `internal/mocks/definition.go`
- Generated mocks should not be edited manually

## Coding Standards

### General Guidelines

1. **Follow Go conventions:** Use gofmt, follow effective Go practices
2. **Vendor dependencies:** Always use `-mod=vendor` flag for builds and tests
3. **Test coverage:** Write tests for new features, maintain or improve coverage
4. **Clean architecture:** Keep business logic in `app/` layer, infrastructure in `infra/`
5. **Error handling:** Use proper error wrapping and return descriptive errors
6. **Logging:** Use logrus for structured logging

### File Organization

- Each domain (auth, feed, item) should have:
  - `domain.go` - Domain models and interfaces
  - `service.go` - Service implementation
  - `commands.go` - Command/request structures
  - `dependencies.go` - Dependency injection setup
  - `service_test.go` - Unit tests

### Testing Guidelines

#### Unit Tests
- Use testify for assertions: `assert` package
- Mock external dependencies using uber/mock
- Table-driven tests are preferred for multiple test cases
- Test files should be named `*_test.go`
- Place tests in the same package as the code being tested

#### UI Tests
- UI tests are located in `tests/ui/` directory
- Use Page Object Model pattern (see `tests/ui/pages/`)
- Tests use Playwright with Chromium browser
- **IMPORTANT**: Always use the backend testserver for feed testing instead of fetching feeds from the internet
  - Mock feeds are served from `http://localhost:3001/feeds/`
  - Available mock feeds are defined in `tests/ui/utils/helpers.js` as `MOCK_FEEDS`
  - Never use real external feed URLs in tests to avoid network blockages
- Test database is separate: `data/test-agg.db`
- Tests run sequentially to avoid database conflicts

### Naming Conventions

- Interfaces should be defined where they are used (consumer side)
- Repository interfaces: `*Repository` (e.g., `FeedRepository`)
- Service interfaces: `*Service` (e.g., `FeedService`)
- Commands/Requests: `*Command` (e.g., `CreateFeedCommand`)
- Use descriptive variable names, avoid single-letter names except for common cases (i, j for loops)

### Database Migrations

- Migrations are in `internal/infra/sqlite/migrations/`
- Use goose for migration management
- Migration files should be numbered sequentially

## Common Tasks

### Adding a New Feature

1. Define domain models and interfaces in `internal/app/<domain>/domain.go`
2. Implement service logic in `internal/app/<domain>/service.go`
3. Add repository implementation in `internal/infra/sqlite/<domain>/`
4. Create HTTP handlers in `internal/infra/http/api/`
5. Write unit tests in `*_test.go` files
6. Update mocks if needed: `make update-mocks`
7. Run pre-commit checks: `make pre-commit`
8. **Run UI tests: `make test-ui`** (if feature affects UI)
9. Ensure all tests pass before completing PR

### Testing Feed-Related Features

**CRITICAL**: When testing features that involve RSS/Atom feeds:
1. **Always use the mock feed server** instead of real internet feeds
2. Start the test server: `make run-test-server`
3. Use mock feed URLs from the testserver:
   - `http://localhost:3001/feeds/tech-news.xml` - RSS 2.0 feed with 3 tech articles
   - `http://localhost:3001/feeds/science-blog.xml` - Atom 1.0 feed with 3 science articles
   - `http://localhost:3001/feeds/empty.xml` - Empty RSS feed for edge cases
4. Mock feeds are located in `internal/testserver/fixtures/`
5. Add new mock feeds in fixtures directory and reference them in tests

**Why use mock feeds?**
- Avoids network blockages in Copilot environment
- Ensures consistent, reliable test results
- Faster test execution
- No dependency on external services

### Adding Dependencies

1. Add to `go.mod` using `go get <package>`
2. Run `make gomod` to update vendor directory
3. Commit both `go.mod`, `go.sum`, and `vendor/` changes

### Fixing Linting Issues

1. Run `make lint` to see issues
2. Run `make gci` to fix import formatting
3. Address other issues as reported by golangci-lint

## Important Notes

- **CGO:** The project uses `CGO_ENABLED=0` for builds to create static binaries
- **Multi-user support:** The application supports multiple users with individual feed subscriptions
- **Migration support:** Legacy Lite Reader data can be migrated from `agg.db`
- **Default port:** Application runs on port 3000 by default
- **Data persistence:** Database file is stored in `data/agg.db`
- **Test database:** UI tests use separate database at `data/test-agg.db`

## Mock Feed System

The project includes a mock feed server for reliable, offline testing:

### Mock Feed Server
- **Location**: `internal/testserver/`
- **Port**: 3001 (when started via `make run-test-server`)
- **Purpose**: Serves mock RSS and Atom feeds locally for testing
- **Benefits**: No internet required, consistent results, faster tests

### Available Mock Feeds
- **tech-news.xml**: RSS 2.0 feed with 3 technology articles
  - URL: `http://localhost:3001/feeds/tech-news.xml`
- **science-blog.xml**: Atom 1.0 feed with 3 science articles
  - URL: `http://localhost:3001/feeds/science-blog.xml`
- **empty.xml**: Empty RSS feed for edge case testing
  - URL: `http://localhost:3001/feeds/empty.xml`

### Adding New Mock Feeds
1. Create a new XML file in `internal/testserver/fixtures/`
2. Follow RSS 2.0 or Atom 1.0 format
3. Add the feed URL to `MOCK_FEEDS` in `tests/ui/utils/helpers.js`
4. Use the feed in your tests

### Using Mock Feeds in Tests
```javascript
import { MOCK_FEEDS } from './utils/helpers.js';

// In your test
await mainPage.addFeed(MOCK_FEEDS.techNews);
```

## Security Considerations

- User passwords are hashed (using bcrypt via golang.org/x/crypto)
- Always validate and sanitize user inputs
- Use prepared statements for database queries (already handled by SQLite driver)
- Be cautious when parsing external RSS feeds

## CI/CD

The project uses GitHub Actions for:
- Running tests (`tests.yaml`)
- Pull request checks (`pr.yml`)
- Release automation (`release.yml`)

**Before Completing a PR:**
1. Ensure all unit tests pass: `make test`
2. **Ensure all UI tests pass: `make test-ui`**
3. Ensure linting passes: `make lint`
4. Verify all CI checks pass

All checks must pass before merging PRs.

## Additional Resources

- [Project README](../README.md)
- [Code of Conduct](../CODE_OF_CONDUCT.md)
- [License](../LICENSE)
- [TODO](../TODO.md)
- [Testing Guide](../TEST.md) - Comprehensive UI and unit testing documentation
- [UI Testing Summary](../UI-TESTING-SUMMARY.md) - UI testing implementation details

## Pull Request Workflow

### Before Completing Any PR

**Required Checks** (must all pass):
1. ✅ Run unit tests: `make test`
2. ✅ Run UI tests: `make test-ui` 
3. ✅ Run linting: `make lint`
4. ✅ Verify all GitHub Actions checks pass
5. ✅ Review test coverage reports

**For Feed-Related Changes:**
- Always test with mock feeds from the testserver
- Never use real internet URLs in tests
- Verify feeds load correctly in UI tests
- Test both RSS 2.0 and Atom 1.0 formats if applicable

**For UI Changes:**
- Run `make test-ui-headed` to visually verify changes
- Ensure all UI tests still pass
- Add new UI tests if introducing new features
- Follow Page Object Model pattern for new tests

### Testing Checklist Template

When working on a PR, use this checklist:

```
- [ ] Unit tests pass (`make test`)
- [ ] UI tests pass (`make test-ui`)
- [ ] Linting passes (`make lint`)
- [ ] Used mock feeds (not internet) for any feed testing
- [ ] Added tests for new features
- [ ] Updated documentation if needed
- [ ] All CI checks pass
```
