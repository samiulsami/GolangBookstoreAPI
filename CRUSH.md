# CRUSH.md

## Build/Lint/Test Commands

- **Build:** `go build -o ./GoBookstoreAPI`
- **Run:** `PORT=3000 ./GoBookstoreAPI serve --port=3000` (or use `make run`)
- **Run All Tests:** `go test ./...`
- **Run Single Test File:** `go test <path_to_package>` (e.g., `go test ./handlers`)
- **Run a Specific Test:** `go test <path_to_package> -run <TestName>` (e.g., `go test ./handlers -run TestGetBook`)
- **Lint/Format:** `go fmt ./...`
- **Static Analysis:** `go vet ./...`

## Code Style Guidelines (Go)

- **Imports:** Group imports. Standard library imports first, followed by third-party imports, separated by a blank line. Use `goimports` to manage imports.
- **Formatting:** Adhere to `go fmt` standards.
- **Naming Conventions:**
    - Exported names (visible outside the package): `CamelCase` (e.g., `MyFunction`, `BookStore`).
    - Unexported names (package-private): `camelCase` (e.g., `myVariable`, `tempFile`).
    - Acronyms: Use all caps (e.g., `HTTPHandler`, `APIResponse`).
- **Error Handling:** Errors are returned as the last return value. Check errors explicitly and handle them immediately using `if err != nil`. Do not ignore errors.
- **Types:** Use explicit types. Define structs for complex data structures.
- **Comments:** Comment exported functions, types, and significant logic. Comments should explain *why* something is done, not *what* it does.
