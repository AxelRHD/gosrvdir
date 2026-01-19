package gosrvdir

import (
	"io"
	"strings"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func RenderListing(w io.Writer, data ListingData) error {
	return Page(data).Render(w)
}

func Page(data ListingData) g.Node {
	title := data.Path
	if title == "/" {
		title = "/"
	}

	return c.HTML5(c.HTML5Props{
		Title:    title + " – gosrvdir",
		Language: "en",
		Head: []g.Node{
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			g.El("style", g.Raw(cssStyles)),
		},
		Body: []g.Node{
			g.Attr("data-theme", data.Theme),
			Nav(
				ThemeSwitcher(data.Theme),
			),
			Header(
				Breadcrumbs(data.Path),
			),
			Main(
				FileTable(data.Entries),
			),
			g.El("script", g.Raw(jsThemeSwitcher)),
		},
	})
}

func Breadcrumbs(path string) g.Node {
	if path == "/" {
		return Div(Class("breadcrumbs"),
			Span(Class("crumb root"), g.Text("/")),
		)
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	var crumbs []g.Node

	// Root link
	crumbs = append(crumbs,
		A(Class("crumb root"), Href("/"), g.Text("~")),
	)

	// Build path progressively
	currentPath := ""
	for i, part := range parts {
		currentPath += "/" + part
		crumbs = append(crumbs,
			Span(Class("separator"), g.Text("/")),
		)
		if i == len(parts)-1 {
			// Last element - not a link
			crumbs = append(crumbs,
				Span(Class("crumb current"), g.Text(part)),
			)
		} else {
			crumbs = append(crumbs,
				A(Class("crumb"), Href(currentPath+"/"), g.Text(part)),
			)
		}
	}

	return Div(Class("breadcrumbs"), g.Group(crumbs))
}

func ThemeSwitcher(current string) g.Node {
	themes := []struct{ value, label string }{
		{"auto", "Auto"},
		{"nord", "Nord"},
		{"squirrel", "Squirrel"},
		{"archlinux", "Archlinux"},
		{"monokai", "Monokai"},
		{"zenburn", "Zenburn"},
	}

	var options []g.Node
	for _, t := range themes {
		opt := Option(Value(t.value), g.Text(t.label))
		if t.value == current {
			opt = Option(Value(t.value), g.Text(t.label), Selected())
		}
		options = append(options, opt)
	}

	return Div(Class("theme-switcher"),
		Label(g.Attr("for", "theme-select"), g.Text("Theme")),
		Select(
			ID("theme-select"),
			g.Attr("onchange", "setTheme(this.value)"),
			g.Group(options),
		),
	)
}

func FileTable(entries []FileInfo) g.Node {
	var rows []g.Node

	for _, entry := range entries {
		var icon, class string
		if entry.Name == ".." {
			icon = "⬆️"
			class = "name parent"
		} else if entry.IsDir {
			icon = "📁"
			class = "name dir"
		} else {
			icon = fileIcon(entry.Name)
			class = "name file"
		}

		rows = append(rows, Tr(
			Td(Class(class),
				Span(Class("icon"), g.Text(icon)),
				A(Href(entry.Path), g.Text(entry.Name)),
			),
			Td(Class("size"), g.Text(entry.Size)),
			Td(Class("date"), g.Text(entry.ModTime)),
		))
	}

	return Table(
		THead(
			Tr(
				Th(Class("name"), g.Text("Name")),
				Th(Class("size"), g.Text("Size")),
				Th(Class("date"), g.Text("Modified")),
			),
		),
		TBody(g.Group(rows)),
	)
}

func fileIcon(name string) string {
	lower := strings.ToLower(name)

	switch {
	// Images
	case strings.HasSuffix(lower, ".jpg"),
		strings.HasSuffix(lower, ".jpeg"),
		strings.HasSuffix(lower, ".png"),
		strings.HasSuffix(lower, ".gif"),
		strings.HasSuffix(lower, ".svg"),
		strings.HasSuffix(lower, ".webp"),
		strings.HasSuffix(lower, ".ico"):
		return "🖼️"

	// Documents
	case strings.HasSuffix(lower, ".pdf"):
		return "📄"
	case strings.HasSuffix(lower, ".doc"),
		strings.HasSuffix(lower, ".docx"),
		strings.HasSuffix(lower, ".odt"):
		return "📝"

	// Code
	case strings.HasSuffix(lower, ".go"),
		strings.HasSuffix(lower, ".py"),
		strings.HasSuffix(lower, ".js"),
		strings.HasSuffix(lower, ".ts"),
		strings.HasSuffix(lower, ".rs"),
		strings.HasSuffix(lower, ".c"),
		strings.HasSuffix(lower, ".cpp"),
		strings.HasSuffix(lower, ".h"),
		strings.HasSuffix(lower, ".java"),
		strings.HasSuffix(lower, ".rb"),
		strings.HasSuffix(lower, ".php"),
		strings.HasSuffix(lower, ".sh"),
		strings.HasSuffix(lower, ".fish"):
		return "📜"

	// Config/Data
	case strings.HasSuffix(lower, ".json"),
		strings.HasSuffix(lower, ".yaml"),
		strings.HasSuffix(lower, ".yml"),
		strings.HasSuffix(lower, ".toml"),
		strings.HasSuffix(lower, ".xml"),
		strings.HasSuffix(lower, ".ini"),
		strings.HasSuffix(lower, ".conf"):
		return "⚙️"

	// Archives
	case strings.HasSuffix(lower, ".zip"),
		strings.HasSuffix(lower, ".tar"),
		strings.HasSuffix(lower, ".gz"),
		strings.HasSuffix(lower, ".bz2"),
		strings.HasSuffix(lower, ".xz"),
		strings.HasSuffix(lower, ".7z"),
		strings.HasSuffix(lower, ".rar"):
		return "📦"

	// Audio
	case strings.HasSuffix(lower, ".mp3"),
		strings.HasSuffix(lower, ".wav"),
		strings.HasSuffix(lower, ".flac"),
		strings.HasSuffix(lower, ".ogg"),
		strings.HasSuffix(lower, ".m4a"):
		return "🎵"

	// Video
	case strings.HasSuffix(lower, ".mp4"),
		strings.HasSuffix(lower, ".mkv"),
		strings.HasSuffix(lower, ".avi"),
		strings.HasSuffix(lower, ".mov"),
		strings.HasSuffix(lower, ".webm"):
		return "🎬"

	// Markdown/Text
	case strings.HasSuffix(lower, ".md"),
		strings.HasSuffix(lower, ".txt"),
		strings.HasSuffix(lower, ".rst"):
		return "📃"

	// HTML/CSS
	case strings.HasSuffix(lower, ".html"),
		strings.HasSuffix(lower, ".htm"),
		strings.HasSuffix(lower, ".css"):
		return "🌐"

	default:
		return "📄"
	}
}

const cssStyles = `
:root {
  --bg: #eceff4;
  --bg-card: #fff;
  --text: #2e3440;
  --text-muted: #4c566a;
  --link-dir: #5e81ac;
  --link-file: #4c566a;
  --header-bg: #e5e9f0;
  --row-hover: #e5e9f0;
  --border: #d8dee9;
  --select-bg: #fff;
  --accent: #5e81ac;
}

[data-theme="nord"] {
  --bg: #2e3440;
  --bg-card: #3b4252;
  --text: #eceff4;
  --text-muted: #d8dee9;
  --link-dir: #88c0d0;
  --link-file: #81a1c1;
  --header-bg: #434c5e;
  --row-hover: #434c5e;
  --border: #4c566a;
  --select-bg: #3b4252;
  --accent: #88c0d0;
}

[data-theme="squirrel"] {
  --bg: #faf8f5;
  --bg-card: #fff;
  --text: #3d3d3d;
  --text-muted: #666;
  --link-dir: #d02474;
  --link-file: #555;
  --header-bg: #f0eeeb;
  --row-hover: #f5f3f0;
  --border: #e0ddd8;
  --select-bg: #fff;
  --accent: #d02474;
}

[data-theme="archlinux"] {
  --bg: #383c4a;
  --bg-card: #404552;
  --text: #fefefe;
  --text-muted: #ccc;
  --link-dir: #03a9f4;
  --link-file: #ea95ff;
  --header-bg: #2f343f;
  --row-hover: #4b5162;
  --border: #4c566a;
  --select-bg: #404552;
  --accent: #03a9f4;
}

[data-theme="monokai"] {
  --bg: #272822;
  --bg-card: #3e3d32;
  --text: #f8f8f2;
  --text-muted: #a6a68a;
  --link-dir: #66d9ef;
  --link-file: #a6e22e;
  --header-bg: #1e1f1c;
  --row-hover: #49483e;
  --border: #49483e;
  --select-bg: #3e3d32;
  --accent: #f92672;
}

[data-theme="zenburn"] {
  --bg: #3f3f3f;
  --bg-card: #4f4f4f;
  --text: #dcdccc;
  --text-muted: #9f9f8f;
  --link-dir: #8cd0d3;
  --link-file: #cc9393;
  --header-bg: #2b2b2b;
  --row-hover: #5f5f5f;
  --border: #6f6f6f;
  --select-bg: #4f4f4f;
  --accent: #f0dfaf;
}

[data-theme="auto"] {
  --bg: #eceff4;
  --bg-card: #fff;
  --text: #2e3440;
  --text-muted: #4c566a;
  --link-dir: #5e81ac;
  --link-file: #4c566a;
  --header-bg: #e5e9f0;
  --row-hover: #e5e9f0;
  --border: #d8dee9;
  --select-bg: #fff;
  --accent: #5e81ac;
}

@media (prefers-color-scheme: dark) {
  [data-theme="auto"] {
    --bg: #2e3440;
    --bg-card: #3b4252;
    --text: #eceff4;
    --text-muted: #d8dee9;
    --link-dir: #88c0d0;
    --link-file: #81a1c1;
    --header-bg: #434c5e;
    --row-hover: #434c5e;
    --border: #4c566a;
    --select-bg: #3b4252;
    --accent: #88c0d0;
  }
}

* {
  box-sizing: border-box;
}

body {
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background: var(--bg);
  color: var(--text);
  margin: 0;
  padding: 0;
  line-height: 1.6;
  min-height: 100vh;
}

nav {
  display: flex;
  justify-content: flex-end;
  padding: 0.75rem 1.5rem;
  border-bottom: 1px solid var(--border);
  background: var(--bg-card);
}

.theme-switcher {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
  color: var(--text-muted);
}

select {
  background: var(--select-bg);
  color: var(--text);
  border: 1px solid var(--border);
  padding: 0.35rem 0.5rem;
  border-radius: 4px;
  font-size: 0.85rem;
  cursor: pointer;
}

select:focus {
  outline: 2px solid var(--accent);
  outline-offset: 1px;
}

header {
  padding: 1rem 1.5rem;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
}

.breadcrumbs {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.25rem;
  font-size: 1.1rem;
}

.crumb {
  color: var(--link-dir);
  text-decoration: none;
  font-weight: 500;
}

a.crumb:hover {
  text-decoration: underline;
}

.crumb.root {
  font-size: 1.2rem;
}

.crumb.current {
  color: var(--text);
  font-weight: 600;
}

.separator {
  color: var(--text-muted);
  margin: 0 0.15rem;
}

main {
  padding: 1rem 1.5rem 2rem;
}

table {
  width: 100%;
  border-collapse: collapse;
  background: var(--bg-card);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
}

thead {
  background: var(--header-bg);
}

th {
  text-align: left;
  padding: 0.75rem 1rem;
  font-weight: 600;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  color: var(--text-muted);
  border-bottom: 1px solid var(--border);
}

td {
  padding: 0.6rem 1rem;
  border-bottom: 1px solid var(--border);
}

tbody tr:last-child td {
  border-bottom: none;
}

tbody tr:hover td {
  background: var(--row-hover);
}

.name {
  width: 55%;
}

.size {
  width: 15%;
  text-align: right;
  color: var(--text-muted);
  font-variant-numeric: tabular-nums;
}

.date {
  width: 30%;
  color: var(--text-muted);
  font-variant-numeric: tabular-nums;
}

th.size, th.date {
  text-align: right;
}

td.date {
  text-align: right;
}

.icon {
  margin-right: 0.5rem;
}

a {
  text-decoration: none;
}

.dir a {
  color: var(--link-dir);
  font-weight: 500;
}

.parent a {
  color: var(--text-muted);
}

.file a {
  color: var(--link-file);
}

td a:hover {
  text-decoration: underline;
}

@media (max-width: 700px) {
  nav, header, main {
    padding-left: 1rem;
    padding-right: 1rem;
  }

  .breadcrumbs {
    font-size: 1rem;
  }

  th, td {
    padding: 0.5rem 0.6rem;
  }

  .date {
    display: none;
  }

  .name {
    width: 70%;
  }

  .size {
    width: 30%;
  }
}

@media (max-width: 400px) {
  .theme-switcher label {
    display: none;
  }
}
`

const jsThemeSwitcher = `
function setTheme(theme) {
  document.body.setAttribute('data-theme', theme);
  localStorage.setItem('gosrvdir-theme', theme);
}

(function() {
  const saved = localStorage.getItem('gosrvdir-theme');
  if (saved) {
    document.body.setAttribute('data-theme', saved);
    const select = document.getElementById('theme-select');
    if (select) select.value = saved;
  }
})();
`
