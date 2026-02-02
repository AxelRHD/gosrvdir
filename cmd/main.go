package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/axelrhd/gosrvdir"
	"github.com/urfave/cli/v3"
)

var appVersion = "dev"

func main() {
	cmd := &cli.Command{
		Name:                  "gosrvdir",
		Usage:                 "Simple directory server with file info",
		Version:               appVersion,
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   8080,
				Usage:   "Port to listen on",
			},
			&cli.StringFlag{
				Name:  "host",
				Value: "0.0.0.0",
				Usage: "Host/interface to bind",
			},
			&cli.StringFlag{
				Name:  "theme",
				Value: "auto",
				Usage: "Color theme (auto, nord, squirrel, archlinux, monokai, zenburn)",
			},
			&cli.StringFlag{
				Name:  "auth",
				Usage: "Basic auth credentials (user:password)",
			},
			&cli.StringFlag{
				Name:  "auth-file",
				Usage: "Path to htpasswd file",
			},
		},
		ArgsUsage: "[directory]",
		Commands: []*cli.Command{
			{
				Name:      "htpasswd",
				Usage:     "Create or update an htpasswd file",
				ArgsUsage: "<file> <username>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.NArg() != 2 {
						return fmt.Errorf("usage: gosrvdir htpasswd <file> <username>")
					}
					return gosrvdir.RunHtpasswd(cmd.Args().Get(0), cmd.Args().Get(1))
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			auth := cmd.String("auth")
			authFile := cmd.String("auth-file")
			if auth != "" && authFile != "" {
				return fmt.Errorf("--auth and --auth-file are mutually exclusive")
			}

			dir := "."
			if cmd.NArg() > 0 {
				dir = cmd.Args().Get(0)
			}

			// Validate directory exists
			info, err := os.Stat(dir)
			if err != nil {
				return fmt.Errorf("cannot access directory: %w", err)
			}
			if !info.IsDir() {
				return fmt.Errorf("%s is not a directory", dir)
			}

			cfg := gosrvdir.Config{
				Host:     cmd.String("host"),
				Port:     int(cmd.Int("port")),
				Dir:      dir,
				Theme:    cmd.String("theme"),
				Auth:     auth,
				AuthFile: authFile,
			}

			return gosrvdir.Serve(cfg)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
