# ── Variables ──────────────────────────────────────────────────────────────────
binary   := "neochess"
cmd_path := "./cmd/neochess"
out_dir  := "bin"

# ── Default: list all recipes ──────────────────────────────────────────────────
default:
    @just --list

# ── Build ──────────────────────────────────────────────────────────────────────

# Build the binary into bin/
build:
    @mkdir -p {{ out_dir }}
    go build -o {{ out_dir }}/{{ binary }} {{ cmd_path }}

# Build with the race detector enabled (useful during development)
build-race:
    @mkdir -p {{ out_dir }}
    go build -race -o {{ out_dir }}/{{ binary }}-race {{ cmd_path }}

# ── Run ────────────────────────────────────────────────────────────────────────

# Build and run the application (pass extra args with: just run -- --addr=:9090)
run *args:
    go run {{ cmd_path }} {{ args }}

# ── Test ───────────────────────────────────────────────────────────────────────

# Run all tests
test:
    go test ./...

# Run all tests with the race detector
test-race:
    go test -race ./...

# Run tests and show coverage in the terminal
test-cover:
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out

# Open coverage report in the browser
test-cover-html: test-cover
    go tool cover -html=coverage.out

# Run a specific package's tests (e.g.: just test-pkg internal/uci)
test-pkg pkg:
    go test -v -race ./{{ pkg }}/...

# ── Audit ──────────────────────────────────────────────────────────────────────

# Vet, staticcheck, and vulnerability scan
audit: vet
    go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Run go vet
vet:
    go vet ./...

# Tidy and verify the module graph
tidy:
    go mod tidy
    go mod verify

# ── Formatting & Linting ───────────────────────────────────────────────────────

# Format all Go source files
fmt:
    go fmt ./...

# Check formatting without modifying files (useful in CI)
fmt-check:
    @test -z "$(gofmt -l .)" || (echo "Unformatted files found:"; gofmt -l .; exit 1)

# ── Clean ──────────────────────────────────────────────────────────────────────

# Remove build artifacts
clean:
    rm -rf {{ out_dir }} coverage.out

# ── Release ────────────────────────────────────────────────────────────────────

# Dry-run GoReleaser locally (no publish, output in dist/)
release-snapshot:
    goreleaser release --snapshot --clean

# Cut a real release — push a tag and the CI workflow does the rest.
# Usage: just release v0.2.0
release version:
    git tag -a {{ version }} -m "Release {{ version }}"
    git push origin {{ version }}

# ── Composite helpers ──────────────────────────────────────────────────────────

# Full pre-commit check: fmt, vet, test with race detector, audit
check: fmt vet test-race audit
