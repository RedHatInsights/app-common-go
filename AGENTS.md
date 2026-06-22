# app-common-go

## Project Overview

app-common-go is a Go library that provides programmatic access to Clowder operator configuration for applications running on cloud.redhat.com infrastructure. It reads JSON configuration files produced by the Clowder operator and exposes them through typed Go interfaces, enabling applications to discover endpoints, database credentials, object storage buckets, Kafka topics, and feature flags without hard-coding environment-specific values.

## Dependencies

**Runtime:**

- Go 1.18+
- Clowder operator (deployment-time dependency — produces the JSON config this library reads)

**Development/Testing:**

- `github.com/stretchr/testify` — test assertions and mocking

**Note:** This library has minimal external dependencies by design. The types in `pkg/api/v1/` are auto-generated from JSON schemas.

## Development Commands

See [Development](README.md#development) in the README for installation and basic usage.

**Run tests (from repository root):**

```bash
ACG_CONFIG="testdata/test.json" go test -v ./...
```

**Agent-driven workflows:**

- When modifying JSON schemas, regenerate Go types by following the schema-to-code generation process (see Architecture section).
- Before committing changes to `pkg/api/v1/`, verify that test fixtures in `pkg/api/v1/testdata/test.json` remain valid against the updated schema.

## Architecture

### Module Structure

```text
app-common-go/
├── pkg/api/v1/          # Auto-generated Go types and config loader
│   ├── types.go         # Clowder config structs (DO NOT EDIT BY HAND)
│   ├── config.go        # Config loader, package globals, CA helpers
│   ├── config_test.go   # Unit tests
│   └── testdata/
│       ├── test.json    # Test fixture: sample Clowder config
│       └── nordsca.json # Minimal fixture: no RDS CA
└── .github/workflows/
    └── package.yml      # CI: runs tests
```

### Core Components

1. **Auto-generated Types (`types.go`):**
   - Defines structs for `AppConfig`, `DatabaseConfig`, `ObjectStoreConfig`, `KafkaConfig`, `FeatureFlagsConfig`, etc.
   - Generated from upstream Clowder JSON schemas — manual edits will be overwritten.

2. **Config Loader (`config.go`):**
   - `LoadConfig(filename string) (*AppConfig, error)` — reads JSON from the specified file path, or from `ACG_CONFIG` env var when used in `init()`.
   - `IsClowderEnabled() bool` — detects whether the app is running in a Clowder-managed environment.
   - Package-level globals (`LoadedConfig`, `KafkaTopics`, etc.) provide pre-indexed lookup maps.

3. **Test Fixtures (`testdata/`):**
   - Minimal valid Clowder config used by CI and local testing.
   - Must stay in sync with schema changes — if you add fields to `types.go`, update this fixture.

### Key Patterns

- **Environment-driven config path:** The library reads `ACG_CONFIG` at runtime. Callers must set this env var before invoking `LoadConfig()`.
- **Schema-first design:** Go types are derived from Clowder's canonical JSON schemas. This ensures compatibility with the operator's output format.
- **Zero-configuration defaults:** If `ACG_CONFIG` is unset or Clowder is disabled, `IsClowderEnabled()` returns false and applications should fall back to local dev defaults.

## Code Style

`golangci-lint` is configured in CI and runs on every PR. Locally, contributors should follow standard Go conventions:

- Run `go fmt` before committing.
- Use `gofmt -s` for simplification rewrites.
- Follow [Effective Go](https://go.dev/doc/effective_go) naming and structure guidelines.

## Common Mistakes

1. **Editing auto-generated files directly.**  
   The types in `pkg/api/v1/types.go` are generated from JSON schemas. Manual changes will be silently overwritten the next time schemas are regenerated. If you need new fields, update the upstream schema and regenerate.

2. **Forgetting to set `ACG_CONFIG` when running tests locally.**  
   Tests require `ACG_CONFIG="testdata/test.json"` to be set (relative to `pkg/api/v1/`). If you run `go test ./pkg/api/v1/...` without it, tests will fail with cryptic errors about missing config. The CI workflow sets this automatically, but local runs do not.

3. **Assuming Clowder config exists in all environments.**  
   Always call `IsClowderEnabled()` before accessing config fields. In local development or non-Clowder deployments, `ACG_CONFIG` may be unset, and dereferencing `LoadConfig()` results can panic. Guard all Clowder-dependent logic with an enabled check.

4. **Breaking the test fixture when adding new required fields.**  
   If you add a new required field to the schema (and thus to `types.go`), you must also update `pkg/api/v1/testdata/test.json` with a valid value. Otherwise, CI will fail on unmarshaling errors. Treat the test fixture as a contract — it must remain a valid example of the schema.

5. **Not checking CI config when documenting test commands.**  
   The canonical test invocation is defined in `.github/workflows/package.yml`. If you document a different command in CONTRIBUTING.md or inline comments, developers will waste time debugging discrepancies. Always cross-check workflow YAML before documenting test procedures.

## Testing

**Run all tests (from repository root):**

```bash
ACG_CONFIG="testdata/test.json" go test -v ./...
```

**Key test file:**

- `pkg/api/v1/config_test.go` — Tests `LoadConfig()`, `IsClowderEnabled()`, and CA file generation with test fixtures.

**Test coverage expectations:**

- CI runs tests on every PR via GitHub Actions.
- No coverage threshold is enforced, but new exported functions should include unit tests.
- Use `testify/assert` for readable assertions.

**Local test tips:**

- To test with a custom config, export `ACG_CONFIG=/path/to/custom.json` before running `go test` from within `pkg/api/v1/`.
- If tests fail with "config file not found," verify `ACG_CONFIG` is set and the path is correct relative to `pkg/api/v1/` (e.g., `testdata/test.json`).
