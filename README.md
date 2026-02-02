<p align="center">
  <img src="logo.svg" alt="gosrvdir logo" width="480">
</p>

<p align="center"><strong>Simple directory server with file info</strong><br>Serve local directories with file sizes, dates, and inline preview.</p>

## Features

- 📁 **Directory listing** — File sizes and modification dates at a glance
- 🧭 **Breadcrumb navigation** — Click through the path hierarchy
- 👁️ **Inline preview** — PDFs, images, and text files display in browser
- 🎨 **Themeable** — 6 color schemes (Auto, Nord, Squirrel, Archlinux, Monokai, Zenburn)
- 🔒 **Basic Auth** — Optional authentication via `--auth` or `--auth-file` (htpasswd/bcrypt)
- 🐚 **Shell completions** — Fish, Bash, Zsh, and PowerShell supported
- ⚡ **Zero dependencies** — Single binary, no runtime required

## Installation

### From source

```bash
go install github.com/axelrhd/gosrvdir/cmd@latest
```

### With just

```bash
git clone https://github.com/axelrhd/gosrvdir.git
cd gosrvdir
just deploy  # Builds, installs to ~/.local/bin, and sets up shell completions
```

The `deploy` command will:
1. Build the binary
2. Install to `~/.local/bin/`
3. Interactively ask which shell completions to install (auto-detects your shell)

For non-interactive installation:
```bash
just deploy-bin                    # Binary only, no completions
just deploy-completion-for fish    # Install completions for specific shell
```

## Usage

```bash
gosrvdir                     # Serve current directory on port 8080
gosrvdir ./mydir             # Serve mydir on port 8080
gosrvdir -p 9000             # Serve current directory on port 9000
gosrvdir -p 9000 ./mydir     # Serve mydir on port 9000
gosrvdir --host 127.0.0.1    # Only listen on localhost
gosrvdir --theme nord        # Use Nord theme
gosrvdir --auth admin:secret # Basic Auth (inline, single user)
gosrvdir --auth-file .htpasswd # Basic Auth (htpasswd file)
```

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `-p, --port` | `8080` | Port to listen on |
| `--host` | `0.0.0.0` | Host/interface to bind |
| `--theme` | `auto` | Color theme (auto, nord, squirrel, archlinux, monokai, zenburn) |
| `--auth` | — | Inline Basic Auth (`user:password`) |
| `--auth-file` | — | Path to htpasswd file (bcrypt) |
| Positional | `.` | Directory to serve |

`--auth` and `--auth-file` are mutually exclusive. Without either flag, no authentication is required.

### Managing htpasswd files

```bash
gosrvdir htpasswd .htpasswd admin    # Add or update user (prompts for password)
```

## Shell Completions

Completions are installed automatically with `just deploy`.

If you installed via `go install` or downloaded a binary, generate completions manually:

```bash
# Fish
gosrvdir completion fish > ~/.config/fish/completions/gosrvdir.fish

# Bash
gosrvdir completion bash > ~/.local/share/bash-completion/completions/gosrvdir

# Zsh
gosrvdir completion zsh > ~/.local/share/zsh/site-functions/_gosrvdir

# PowerShell
gosrvdir completion pwsh > ~/.config/powershell/gosrvdir.ps1
# Then add to your profile: . ~/.config/powershell/gosrvdir.ps1
```

## Why gosrvdir?

| Tool | File Info | Inline Preview |
|------|-----------|----------------|
| `python3 -m http.server` | ❌ | ✅ |
| `miniserve` | ✅ | ❌ (forces download) |
| `gosrvdir` | ✅ | ✅ |

## License

[MIT](LICENSE)
