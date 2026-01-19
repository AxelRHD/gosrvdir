# gosrvdir

Minimaler HTTP-Server zum Serven lokaler Verzeichnisse.

## Motivation

- `python3 -m http.server` zeigt keine Dateiinfos (Größe, Datum)
- `miniserve` erzwingt Download statt Inline-Anzeige im Browser

## Features

- Directory Listing mit Dateigröße und Änderungsdatum
- Inline-Anzeige von Dateien im Browser (PDFs etc.)
- Minimale Konfiguration

## CLI

```bash
gosrvdir                     # pwd auf Port 8080
gosrvdir ./mydir             # mydir auf Port 8080
gosrvdir -p 9000             # pwd auf Port 9000
gosrvdir -p 9000 ./mydir     # mydir auf Port 9000
gosrvdir --host 127.0.0.1    # nur localhost
```

## Parameter

| Flag | Default | Beschreibung |
|------|---------|--------------|
| `-p, --port` | `8080` | Port |
| `-h, --host` | `0.0.0.0` | Host/Interface |
| Positional | `.` | Verzeichnis |

## Stack

- Go
- urfave/cli für CLI
- Standardbibliothek für HTTP

## Directory Listing

Einfaches HTML mit:
- Dateiname (Link)
- Größe (human-readable: KB, MB, GB)
- Änderungsdatum
- Sortierung nach Name (vorerst)

## Wichtig

- Kein Content-Disposition Header setzen (Browser entscheidet)
- Symlinks folgen (normales Verhalten)
- Keine Authentifizierung nötig
- Keine Upload-Funktion
