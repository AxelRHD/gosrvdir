app_name := "gosrvdir"
bin_dir := "bin"
bin_file := bin_dir / app_name

# Auto-version: tag if on tag, otherwise tag-hash or dev-hash
app_version := `git describe --tags --exact-match 2>/dev/null || { tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev"); hash=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); echo "${tag}-${hash}"; }`

[private]
default:
    @just --list --unsorted

# ============================================================
# Development
# ============================================================

# Run program directly
[group('dev')]
run *args:
    @go run ./cmd {{args}}

# Format code
[group('dev')]
fmt:
    @go fmt ./...

# Static analysis
[group('dev')]
vet:
    @go vet ./...

# ============================================================
# Build
# ============================================================

# Build binary
[group('build')]
build:
    @mkdir -p {{bin_dir}}
    @go build -ldflags "-X 'main.appVersion={{app_version}}'" -o {{bin_file}} ./cmd

# Build binaries for all platforms
[group('build')]
build-all:
    @mkdir -p {{bin_dir}}
    @GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.appVersion={{app_version}}'" -o {{bin_dir}}/{{app_name}}-linux-amd64 ./cmd
    @GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.appVersion={{app_version}}'" -o {{bin_dir}}/{{app_name}}-linux-arm64 ./cmd
    @GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.appVersion={{app_version}}'" -o {{bin_dir}}/{{app_name}}-darwin-amd64 ./cmd
    @GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.appVersion={{app_version}}'" -o {{bin_dir}}/{{app_name}}-darwin-arm64 ./cmd

# ============================================================
# Install
# ============================================================

# Install locally (go install)
[group('install')]
install:
    @go install -ldflags "-X 'main.appVersion={{app_version}}'" ./cmd

# ============================================================
# Deploy
# ============================================================

# Deploy binary locally to ~/.local/bin
[group('deploy')]
deploy: build deploy-bin

# Copy binary to ~/.local/bin
[group('deploy')]
deploy-bin:
    @mkdir -p ~/.local/bin
    @cp {{bin_file}} ~/.local/bin/{{app_name}}
    @echo "Installed {{app_name}} to ~/.local/bin/"

# ============================================================
# Clean
# ============================================================

# Remove build artifacts
[group('clean')]
clean:
    @rm -rf {{bin_dir}}
