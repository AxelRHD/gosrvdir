package gosrvdir

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Host     string
	Port     int
	Dir      string
	Theme    string
	Auth     string
	AuthFile string
}

func Serve(cfg Config) error {
	absDir, err := filepath.Abs(cfg.Dir)
	if err != nil {
		return fmt.Errorf("cannot resolve path: %w", err)
	}

	var creds Credentials
	if cfg.Auth != "" {
		parts := strings.SplitN(cfg.Auth, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --auth format, expected user:password")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(parts[1]), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hashing password: %w", err)
		}
		creds = Credentials{parts[0]: string(hash)}
	} else if cfg.AuthFile != "" {
		var err error
		creds, err = ParseHtpasswd(cfg.AuthFile)
		if err != nil {
			return fmt.Errorf("reading auth file: %w", err)
		}
	}

	handler := &Handler{
		Dir:   absDir,
		Theme: cfg.Theme,
		Creds: creds,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf("Serving %s at http://%s\n", absDir, addr)

	return http.ListenAndServe(addr, handler)
}
