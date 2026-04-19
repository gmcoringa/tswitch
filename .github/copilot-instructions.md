# GitHub Copilot Instructions for tswitch

## Project Context

tswitch is a Go CLI tool that switches terraform/tofu and terragrunt versions based on semver constraints in `terragrunt.hcl` files.

## Code Style

- Use Go 1.24+ idioms.
- Always handle errors explicitly — return `error` in library functions, log + exit in `main.go`.
- Use `logrus` (aliased as `log`) for all logging.
- Follow existing import alias conventions: `lio` for `pkg/io`, `tf` for `pkg/terraform`, `tg` for `pkg/terragrunt`.
- Write Go-doc comments on exported functions.
- Tests use `testify` assertions and `golang/mock` for interface mocks.

## When Adding a New Tool Resolver

1. Create a new struct implementing `lib.Resolver` in a new or existing package under `pkg/`.
2. Implement: `ListVersions()`, `Name()`, `Implementation()`, `AddNewVersion(version, destination)`.
3. Wire it into `main.go` using `lib.CreateInstaller()`.
4. Add mock generation to `Makefile` if the new package exposes interfaces.

## Testing Guidelines

- Place tests in the same package (`*_test.go`).
- Use `test_data/` subdirectories for fixtures.
- Run `make test` to validate (includes `-race` flag).
- Run `make lint` before pushing.

## Things to Avoid

- Do not modify `mocks/*.go` directly — regenerate with `make mocks`.
- Do not hardcode OS/architecture — use `runtime.GOOS` and `runtime.GOARCH`.
- Do not add new linter exclusions without justification.
