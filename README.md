# mediactl

`mediactl` is a Go CLI for sidecar-driven media workflows. It exports `ffprobe`
JSON snapshots next to your media files and remuxes MP4 sources into MKV while
applying `handler_name` changes you edit directly in the sidecar.

## What it does

- Export ffprobe JSON sidecars for MKV and MP4 libraries.
- Remux MP4 into MKV and update `handler_name` tags on audio streams.
- Walk directories recursively with predictable sidecar naming.
- Log verbose ffmpeg output and the exact command before execution.

## Requirements

- Go 1.22+
- `ffprobe` on `$PATH` (FFmpeg suite)
- `ffmpeg` on `$PATH` (for remux/import)
- Read/write permissions for the media tree

## Install

```bash
# Build locally
GOOS=linux GOARCH=amd64 go build -buildvcs=false -o mediactl .

# Or install into $GOBIN
GOOS=linux GOARCH=amd64 go install ./...
```

## Usage

```bash
mediactl --help
mediactl sidecar --help
```

### Export sidecars

Export MKV sidecars:

```bash
mediactl sidecar export --format mkv /path/to/media
```

Export MP4 sidecars:

```bash
mediactl sidecar export --format mp4 /path/to/media
```

Export behavior:

- Walks the provided directory (defaults to `.` if omitted).
- Runs ffprobe with a focused set of stream fields.
- Writes `file.ext.json` next to the media file.
- If a sidecar exists, writes `file.ext.1.json`, `file.ext.2.json`, etc.
- Normalizes any `HANDLER_NAME` tag keys to lowercase `handler_name`.

### Import sidecars (remux into MKV)

After editing a sidecar by hand, remux the corresponding MP4 into MKV with
updated `handler_name` tags for audio streams:

```bash
mediactl sidecar import "demo-data/Movie Title.mp4.json"
```

Import behavior:

- The media file is derived from the sidecar name.
  - `Movie.mp4.json` -> `Movie.mp4`
  - `Movie.mp4.1.json` -> `Movie.mp4`
- Output is `Movie.remux.mkv`.
- Only audio streams with changed `handler_name` values are updated.
- Data streams are dropped (Matroska does not support them).
- The ffmpeg command is printed before execution; ffmpeg runs in verbose mode.

## Project layout

- `main.go`: CLI entrypoint.
- `cmd/`: Cobra commands (`sidecar export`, `sidecar import`).
- `core/sidecar/`: sidecar discovery and naming helpers.
- `infra/exec/`: ffprobe and ffmpeg runners.

## Development

```bash
go build ./...
```

```bash
go test ./...
```

## License

MIT © 2025 Bartholomaeuss (see `LICENSE`).
