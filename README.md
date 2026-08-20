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

## Notes

- Output is uncompressed text pages, so it is roughly 45% gzip-compressible. If the path under test applies transfer compression, size on the wire will be smaller.
- Peak memory is several times the output size (about 6x measured at 500 MB).

## Build

```bash
go build -o goPdfGen .
```
