# goPdfGen

Generates a valid, uncompressed PDF of an exact size, for testing upload
limits and other size-sensitive paths.

- Exact byte count: target size in MB, converted as `size * 1024 * 1024`.
- Optional searchable signature text embedded in the page content.
- Deterministic output: same seed and size produce an identical file.

## Usage

```bash
./goPdfGen -size 30 -out /path/to/name.pdf
```

```bash
./goPdfGen -size 4096 -signature "NEEDLE-7f3a9c" -out big.pdf
```

## Flags

| Flag | Default | Purpose |
|------|---------|---------|
| `-size` | `5` | Target size in megabytes (exact byte count = size * 1024 * 1024). |
| `-out` | `out.pdf` | Output path. |
| `-seed` | `0` | Random seed for the page text; `0` uses the current time. Same seed + size gives an identical file. |
| `-signature` | (none) | Text embedded in the page content so it can be found by full-text search. |
| `-sig-page` | `last` | Where the signature goes: `first`, `last`, `all`, or a page number. |

## Notes

- Output is uncompressed text pages, so it is roughly 45% gzip-compressible.
  If the path under test applies transfer compression, size on the wire will
  be smaller.
- Peak memory is several times the output size (about 6x measured at 500 MB).

## Build

```bash
go build -o goPdfGen .
```

## Releases

Tagged releases ship binaries for Linux, Windows, and macOS (amd64 + arm64)
with a `SHA256SUMS` file.

## License

Apache License 2.0. See [LICENSE](LICENSE).
