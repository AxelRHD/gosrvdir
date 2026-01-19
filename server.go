package gosrvdir

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type Config struct {
	Host  string
	Port  int
	Dir   string
	Theme string
}

func Serve(cfg Config) error {
	absDir, err := filepath.Abs(cfg.Dir)
	if err != nil {
		return fmt.Errorf("cannot resolve path: %w", err)
	}

	handler := &Handler{
		Dir:   absDir,
		Theme: cfg.Theme,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf("Serving %s at http://%s\n", absDir, addr)

	return http.ListenAndServe(addr, handler)
}
