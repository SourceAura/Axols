# AGENTS.md

## Essential Commands
- `go build` - Build the project
- `go test ./...` - Run all tests
- `go run main.go` - Run the application
- `go fmt ./...` - Format code (if using Go)

## Code Organization
- Project uses Go modules (go.mod/go.sum)
- Main entry point: `main.go`
- No significant subdirectories found

## Naming Conventions
- Go package names match directory structure
- Use snake_case for variables/functions
- Capitalize for exported identifiers

## Testing Approach
- Tests in `*_test.go` files
- Run with `go test ./...`
- No CI config found - add to .github/workflows

## Gotchas
- No CI setup - add GitHub Actions
- No linting - consider adding `golangci-lint`
- No documentation - add README sections

## Project-Specific Context
- Python requirements.txt suggests possible CI dependencies
- No existing AGENTS.md - this is the first version

## Recommendations
1. Add CI configuration
2. Implement linting
3. Expand documentation

---
Generated with Crush
