# mediactl

`mediactl` is a tiny Go CLI built with Cobra that helps you inspect Matroska (MKV) media files. The first available workflow runs `ffprobe` across all `.mkv` files in a directory tree and produces JSON sidecars that you can diff, catalog, or feed into other tooling.

## Quick Start

```bash
git clone https://example.com/your/mediactl.git
cd mediactl
GOOS=windows GOARCH=amd64 go build -buildvcs=false -o mediactl.exe .
```

The binary lives at `./mediactl` after `go build`. Install it globally with `go install`.

## Requirements

- Go 1.22+
- `ffprobe` (part of the FFmpeg suite) available on `$PATH`
- Read/write access to the directories you want to scan

## Usage

`mediactl` exposes a hierarchical command structure. Use `--help` at any level to see available options.

```bash
mediactl --help
mediactl mkv --help
```

### Generate ffprobe sidecars for MKVs

```bash
mediactl mkv ffprobe /path/to/videos
```

- Recursively finds every `.mkv` under the given directory (current directory if omitted).
- Runs `ffprobe -v error -show_entries stream=index,codec_name,codec_type,profile:stream_disposition=default:stream_tags=title -print_format json` on each file.
- Writes a JSON sidecar (`yourfile.mkv.json`) next to the original MKV.
- Prints the path to every sidecar it writes; if no MKVs are found you get a friendly notice instead.

Use a scratch directory or copies of production files the first time you try it to ensure the output structure matches your expectations.

### Generate ffprobe sidecars for MP4s

```bash
mediactl mp4 ffprobe /path/to/videos
```

Behavior matches the MKV command but filters on `.mp4` files.

## Development

The repo is designed to work without extra setup:

```bash
go build ./...
```

### Testing

`cmd/mkv_ffprobe_test.go` contains unit tests for the recursive ffprobe walker. They replace the real `ffprobe` call with a fake runner, so the tests do not require FFmpeg.

```bash
go test ./...
```

### Project Layout

- `main.go` wires Cobra’s root command.
- `cmd/root.go` declares the root CLI node.
- `cmd/mkv.go` hosts the `mkv` namespace.
- `cmd/mkv_ffprobe.go` implements the recursive ffprobe sidecar export logic.
- `cmd/mkv_ffprobe_test.go` provides regression tests for the walker.

## Contributing

1. Fork the repo and create a feature branch.
2. Make your changes with clear commits.
3. Run `go test ./...`.
4. Open a PR with context on the media workflows you’re targeting.

## License

MIT © 2025 Bartholomaeuss (see `LICENSE`).
