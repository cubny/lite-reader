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

### Testing

- `make test` - Run all unit tests with coverage
- Coverage reports are generated in `reports/` directory
- Aim for high test coverage (current: ~36.3% total, with 100% in feed and item app logic)

### Code Quality

- `make lint` - Run golangci-lint for code quality checks
- `make gci` - Format imports using gci
- `make pre-commit` - Run all essential checks (gomod, update-mocks, lint, test)

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

- Use testify for assertions: `assert` package
- Mock external dependencies using uber/mock
- Table-driven tests are preferred for multiple test cases
- Test files should be named `*_test.go`
- Place tests in the same package as the code being tested

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

Ensure all checks pass before merging PRs.

## Additional Resources

- [Project README](../README.md)
- [Code of Conduct](../CODE_OF_CONDUCT.md)
- [License](../LICENSE)
- [TODO](../TODO.md)
