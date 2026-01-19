app_name := "gosrvdir"
bin_dir := "bin"
bin_file := bin_dir / app_name

# Auto-version: tag if on tag, otherwise tag-hash or dev-hash
app_version := `git describe --tags --exact-match 2>/dev/null || { tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev"); hash=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); echo "${tag}-${hash}"; }`

# Detect current shell from $SHELL
detected_shell := `basename "$SHELL" 2>/dev/null || echo "unknown"`

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

# Deploy binary and completions locally (interactive shell selection)
[group('deploy')]
deploy: build deploy-bin deploy-completion

# Copy binary to ~/.local/bin
[group('deploy')]
deploy-bin:
    @mkdir -p ~/.local/bin
    @cp {{bin_file}} ~/.local/bin/{{app_name}}
    @echo "Installed {{app_name}} to ~/.local/bin/"

# Generate and install shell completions (interactive)
[group('deploy')]
deploy-completion:
    #!/usr/bin/env sh
    detected="{{detected_shell}}"

    # Map detected shell to supported shells
    case "$detected" in
        fish|bash|zsh|pwsh) default="$detected" ;;
        *) default="fish" ;;
    esac

    echo "Detected shell: $detected"
    echo ""
    echo "Install completions for which shell?"
    echo "  [f] fish"
    echo "  [b] bash"
    echo "  [z] zsh"
    echo "  [p] pwsh (PowerShell)"
    echo "  [n] none (skip)"
    echo ""
    printf "Choice [$default]: "
    read choice

    # Use default if empty
    choice="${choice:-$default}"

    case "$choice" in
        f|fish)
            mkdir -p ~/.config/fish/completions
            {{bin_file}} completion fish > ~/.config/fish/completions/{{app_name}}.fish
            echo "Installed fish completions to ~/.config/fish/completions/"
            ;;
        b|bash)
            mkdir -p ~/.local/share/bash-completion/completions
            {{bin_file}} completion bash > ~/.local/share/bash-completion/completions/{{app_name}}
            echo "Installed bash completions to ~/.local/share/bash-completion/completions/"
            ;;
        z|zsh)
            mkdir -p ~/.local/share/zsh/site-functions
            {{bin_file}} completion zsh > ~/.local/share/zsh/site-functions/_{{app_name}}
            echo "Installed zsh completions to ~/.local/share/zsh/site-functions/"
            ;;
        p|pwsh)
            mkdir -p ~/.config/powershell
            {{bin_file}} completion pwsh > ~/.config/powershell/{{app_name}}.ps1
            echo "Installed pwsh completions to ~/.config/powershell/"
            echo "Add to your profile: . ~/.config/powershell/{{app_name}}.ps1"
            ;;
        n|none)
            echo "Skipped shell completions"
            ;;
        *)
            echo "Unknown choice: $choice"
            exit 1
            ;;
    esac

# Install completions for a specific shell (non-interactive)
[group('deploy')]
deploy-completion-for shell:
    #!/usr/bin/env sh
    case "{{shell}}" in
        fish)
            mkdir -p ~/.config/fish/completions
            {{bin_file}} completion fish > ~/.config/fish/completions/{{app_name}}.fish
            echo "Installed fish completions"
            ;;
        bash)
            mkdir -p ~/.local/share/bash-completion/completions
            {{bin_file}} completion bash > ~/.local/share/bash-completion/completions/{{app_name}}
            echo "Installed bash completions"
            ;;
        zsh)
            mkdir -p ~/.local/share/zsh/site-functions
            {{bin_file}} completion zsh > ~/.local/share/zsh/site-functions/_{{app_name}}
            echo "Installed zsh completions"
            ;;
        pwsh)
            mkdir -p ~/.config/powershell
            {{bin_file}} completion pwsh > ~/.config/powershell/{{app_name}}.ps1
            echo "Installed pwsh completions"
            echo "Add to your profile: . ~/.config/powershell/{{app_name}}.ps1"
            ;;
        *)
            echo "Unknown shell: {{shell}}. Available: fish, bash, zsh, pwsh"
            exit 1
            ;;
    esac

# ============================================================
# Clean
# ============================================================

# Remove build artifacts
[group('clean')]
clean:
    @rm -rf {{bin_dir}}
