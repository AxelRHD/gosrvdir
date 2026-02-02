# gosrvdir

Minimaler HTTP-Server zum Serven lokaler Verzeichnisse.

## Motivation

- `python3 -m http.server` zeigt keine Dateiinfos (Größe, Datum)
- `miniserve` erzwingt Download statt Inline-Anzeige im Browser

## Features

- Directory Listing mit Dateigröße und Änderungsdatum
- Inline-Anzeige von Dateien im Browser (PDFs etc.)
- Optionale Basic Auth (`--auth user:pass` oder `--auth-file htpasswd`)
- htpasswd-Subcommand zum Erstellen/Verwalten von Passwort-Dateien
- Minimale Konfiguration

## CLI

```bash
gosrvdir                     # pwd auf Port 8080
gosrvdir ./mydir             # mydir auf Port 8080
gosrvdir -p 9000             # pwd auf Port 9000
gosrvdir -p 9000 ./mydir     # mydir auf Port 9000
gosrvdir --host 127.0.0.1    # nur localhost
gosrvdir --auth admin:secret # Basic Auth (inline)
gosrvdir --auth-file .htpasswd # Basic Auth (htpasswd-Datei)
gosrvdir htpasswd .htpasswd admin # User anlegen/aktualisieren
```

## Parameter

| Flag | Default | Beschreibung |
|------|---------|--------------|
| `-p, --port` | `8080` | Port |
| `-h, --host` | `0.0.0.0` | Host/Interface |
| `--auth` | — | Inline-Auth (`user:password`) |
| `--auth-file` | — | Pfad zu htpasswd-Datei (bcrypt) |
| Positional | `.` | Verzeichnis |

`--auth` und `--auth-file` schließen sich gegenseitig aus.

## Stack

- Go
- urfave/cli für CLI
- Standardbibliothek für HTTP
- golang.org/x/crypto/bcrypt für Passwort-Hashing

## Directory Listing

Einfaches HTML mit:
- Dateiname (Link)
- Größe (human-readable: KB, MB, GB)
- Änderungsdatum
- Sortierung nach Name (vorerst)

## Wichtig

- Kein Content-Disposition Header setzen (Browser entscheidet)
- Symlinks folgen (normales Verhalten)
- Optionale Basic Auth (ohne Flags → kein Auth)
- Keine Upload-Funktion
