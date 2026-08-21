# goPdfGen

Generates a valid, uncompressed PDF of an exact size in MB, for testing upload limits and other size-sensitive paths.

## Usage

```bash
./goPdfGen -size 30 -out /path/to/name.pdf
```

## Flags

- `-size` target size in megabytes (exact byte count = size * 1024 * 1024). Default `5`.
- `-out` output path. Default `out.pdf`.
- `-seed` random seed for the page text; `0` (default) uses the current time. Same seed + size = identical file.
- `-signature` text embedded in the page content so it can be found by full-text search. Empty (default) = none.
- `-sig-page` where the signature goes: `first`, `last` (default), `all`, or a page number.

```bash
./goPdfGen -size 4096 -signature "NEEDLE-7f3a9c" -out big.pdf
```

## Notes

- Output is uncompressed text pages, so it is roughly 45% gzip-compressible. If the path under test applies transfer compression, size on the wire will be smaller.
- Peak memory is several times the output size (about 6x measured at 500 MB).

## Releases

Tagged releases ship binaries for Linux, Windows, and macOS (amd64 + arm64) with a `SHA256SUMS` file.

## Build

```bash
go build -o goPdfGen .
```
