# AI Agent Instructions for tswitch

## Project Overview

**tswitch** is a CLI tool written in Go that automatically switches between versions of [Terraform](https://www.terraform.io/), [OpenTofu](https://opentofu.org/), and [Terragrunt](https://terragrunt.gruntwork.io/) by reading version constraints from `terragrunt.hcl` files.

- **Module path**: `github.com/gmcoringa/tswitch`
- **Go version**: 1.24+
- **License**: GNU General Public License v3

## Architecture

### Entry Point

`main.go` — Parses flags, loads configuration, reads HCL constraints, and triggers the install pipeline for both terraform/tofu and terragrunt.

### Package Layout

```
pkg/
├── configuration/  # CLI flags + YAML config file loading (Config struct)
├── db/             # Local version database (scribble JSON store)
├── hcl/            # terragrunt.hcl parser for version constraints
├── io/             # File operations (symlinks, downloads, move, permissions)
├── lib/            # Core installer logic + Resolver interface
├── log/            # Custom logrus formatter (message-only output)
├── terraform/      # Terraform & OpenTofu resolvers (version listing + download)
├── terragrunt/     # Terragrunt resolver (git tag listing + download)
└── util/           # String utilities (IsBlank/IsNotBlank)
```

### Key Interfaces

- **`lib.Installer`** — `Install(constraint string)` — orchestrates finding, downloading, caching, and symlinking a binary version.
- **`lib.Resolver`** — `ListVersions()`, `Name()`, `Implementation()`, `AddNewVersion(version, destination)` — implemented by `terraform.Terraform`, `terraform.Tofu`, and `terragrunt.Terragrunt`.
- **`db.Database`** — `Get()`, `GetCurrent()`, `SetCurrent()`, `Add()` — version tracking with JSON-file-backed storage via scribble.

### Install Flow

1. Parse `terragrunt.hcl` → extract `terraform_version_constraint` and `terragrunt_version_constraint`.
2. For each tool (terraform/tofu + terragrunt):
   - List all available versions (HTTP scraping for Terraform, git ls-remote for Tofu/Terragrunt).
   - Resolve the best matching version using semver constraints.
   - Check local cache DB → if already installed, just update the symlink.
   - Otherwise, download + extract + cache + symlink the binary.

## Development Commands

```bash
# Build
make build                  # Produces dist/tswitch

# Test
make test                   # Runs all tests with -race flag

# Lint
make lint                   # Downloads golangci-lint v2.1.6 and runs it

# Format
make fmt                    # go fmt ./...

# Generate mocks
make mocks                  # Uses mockgen v1.5.0

# Snapshot release
make goreleaser             # goreleaser release --snapshot --clean

# Tag a release
make promote                # Auto-increments patch version and pushes tag
```

## Code Conventions

### Style & Patterns

- **Error handling**: Log errors with `logrus` and `os.Exit(1)` in main flow; return errors in library code.
- **Logging**: Use `logrus` (`log` alias) everywhere. The custom formatter in `pkg/log` outputs only messages (no timestamps/levels) to stderr.
- **Package aliases**: Import aliases are consistently used:
  - `log` → `github.com/sirupsen/logrus`
  - `lio` → `github.com/gmcoringa/tswitch/pkg/io`
  - `tf` → `github.com/gmcoringa/tswitch/pkg/terraform`
  - `tg` → `github.com/gmcoringa/tswitch/pkg/terragrunt`
  - `formatter` → `github.com/gmcoringa/tswitch/pkg/log`
- **Mocks**: Generated with `golang/mock` mockgen into `mocks/` directory. Mock sources are `pkg/db/db.go` and `pkg/lib/installer.go`.
- **Comments**: Go-doc style (`// FunctionName : description`).
- **Configuration priority**: CLI flags → config YAML file → built-in defaults.

### Testing

- Tests live alongside source files as `*_test.go`.
- Test fixtures use local `test_data/` directories or inline HCL/YAML files (e.g., `pkg/hcl/terragrunt.hcl`).
- Use `github.com/stretchr/testify` for assertions.
- Use `github.com/golang/mock` for interface mocking.

### Linting

- golangci-lint v2 with a curated set of linters: `dupl`, `gocritic`, `gocyclo`, `gosec`, `govet`, `nakedret`, `prealloc`, `revive`, `staticcheck`, `unconvert`, `unused`.
- Max cyclomatic complexity: 15.
- `mocks/`, `*_test.go`, `third_party/`, and `examples/` are excluded from linting.

## CI/CD

- **GitHub Actions** workflow in `.github/workflows/push.yaml`.
- Triggers on: push to `main`, tags `v*`, and pull requests.
- Jobs: `lint` → `test` (Linux) → `test-macos` (best-effort) → `snapshot` (PR) or `release` (tag).
- Releases use **GoReleaser** (`.goreleaser.yaml`) to cross-compile for darwin/linux/windows on amd64/386/arm64.
- **Dependabot** updates Go modules and GitHub Actions weekly.

## Important Notes

- **Do not modify** files in `mocks/` directly — they are auto-generated. Run `make mocks` after changing source interfaces.
- **Do not commit** the `golangci-lint` binary or `dist/` directory (both are in `.gitignore`).
- The `Resolver` interface is the extension point for adding new tool support. Each resolver knows how to list versions and download binaries for its tool.
- Terraform versions are scraped from `releases.hashicorp.com`; Tofu and Terragrunt use git ls-remote to list tags.
- The terraform implementation choice (`terraform` vs `tofu`) is resolved at startup via config/flag and routed through `pkg/terraform/resolver.go`.
