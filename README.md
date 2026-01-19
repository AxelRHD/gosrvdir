<p align="center">
  <img src="logo.svg" alt="gosrvdir logo" width="480">
</p>

<p align="center"><strong>Simple directory server with file info</strong><br>Serve local directories with file sizes, dates, and inline preview.</p>

## Features

- 📁 **Directory listing** — File sizes and modification dates at a glance
- 🧭 **Breadcrumb navigation** — Click through the path hierarchy
- 👁️ **Inline preview** — PDFs, images, and text files display in browser
- 🎨 **Themeable** — 6 color schemes (Auto, Nord, Squirrel, Archlinux, Monokai, Zenburn)
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
just deploy  # Builds and installs to ~/.local/bin
```

## Usage

```bash
gosrvdir                     # Serve current directory on port 8080
gosrvdir ./mydir             # Serve mydir on port 8080
gosrvdir -p 9000             # Serve current directory on port 9000
gosrvdir -p 9000 ./mydir     # Serve mydir on port 9000
gosrvdir --host 127.0.0.1    # Only listen on localhost
gosrvdir --theme nord        # Use Nord theme
```

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `-p, --port` | `8080` | Port to listen on |
| `--host` | `0.0.0.0` | Host/interface to bind |
| `--theme` | `auto` | Color theme (auto, nord, squirrel, archlinux, monokai, zenburn) |
| Positional | `.` | Directory to serve |

## Why gosrvdir?

| Tool | File Info | Inline Preview |
|------|-----------|----------------|
| `python3 -m http.server` | ❌ | ✅ |
| `miniserve` | ✅ | ❌ (forces download) |
| `gosrvdir` | ✅ | ✅ |

## License

[MIT](LICENSE)
