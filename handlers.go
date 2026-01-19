package gosrvdir

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type Handler struct {
	Dir   string
	Theme string
}

type FileInfo struct {
	Name    string
	Path    string
	Size    string
	ModTime string
	IsDir   bool
}

type ListingData struct {
	Path    string
	Theme   string
	Entries []FileInfo
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean and resolve path
	urlPath := path.Clean(r.URL.Path)
	if urlPath == "" {
		urlPath = "/"
	}

	filePath := filepath.Join(h.Dir, filepath.FromSlash(urlPath))

	// Security: ensure we don't escape the root directory
	if !strings.HasPrefix(filePath, h.Dir) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if info.IsDir() {
		h.serveDirectory(w, r, filePath, urlPath)
	} else {
		h.serveFile(w, r, filePath)
	}
}

func (h *Handler) serveDirectory(w http.ResponseWriter, r *http.Request, filePath, urlPath string) {
	// Ensure trailing slash for directories
	if !strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
		return
	}

	entries, err := os.ReadDir(filePath)
	if err != nil {
		http.Error(w, "Cannot read directory", http.StatusInternalServerError)
		return
	}

	var files []FileInfo

	// Add parent directory link if not at root
	if urlPath != "/" {
		files = append(files, FileInfo{
			Name:  "..",
			Path:  path.Dir(strings.TrimSuffix(urlPath, "/")),
			IsDir: true,
		})
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := entry.Name()
		entryPath := path.Join(urlPath, name)

		fi := FileInfo{
			Name:    name,
			Path:    entryPath,
			ModTime: info.ModTime().Format("2006-01-02 15:04"),
			IsDir:   entry.IsDir(),
		}

		if entry.IsDir() {
			fi.Name += "/"
			fi.Path += "/"
		} else {
			fi.Size = formatSize(info.Size())
		}

		files = append(files, fi)
	}

	// Sort: directories first, then by name
	sort.Slice(files, func(i, j int) bool {
		// Keep ".." at the top
		if files[i].Name == ".." {
			return true
		}
		if files[j].Name == ".." {
			return false
		}
		// Directories before files
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		// Alphabetical (case-insensitive)
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	data := ListingData{
		Path:    urlPath,
		Theme:   h.Theme,
		Entries: files,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	RenderListing(w, data)
}

func (h *Handler) serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	// Don't set Content-Disposition — let browser decide (inline preview)
	http.ServeFile(w, r, filePath)
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
