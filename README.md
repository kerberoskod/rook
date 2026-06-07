# Rook — Dependency Migration CLI

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
[![CI](https://github.com/kerberoskod/rook/actions/workflows/ci.yml/badge.svg)](https://github.com/kerberoskod/rook/actions/workflows/ci.yml)
[![Release](https://github.com/kerberoskod/rook/actions/workflows/release.yml/badge.svg)](https://github.com/kerberoskod/rook/actions/workflows/release.yml)

**Rook** scans your project files, checks dependency versions against registries, and helps you update them — all from the terminal.

## Installation

```bash
go install github.com/kerberoskod/rook@latest
```

Or download a prebuilt binary from [releases](https://github.com/kerberoskod/rook/releases).

## Usage

```bash
# Scan dependencies in the current directory
rook scan

# Scan a specific project
rook scan --path /path/to/project

# Check for outdated dependencies
rook check

# Exit with error if any are outdated (for CI)
rook check --strict

# Update to latest versions
rook update

# Preview updates without modifying files
rook update --dry-run

# JSON output
rook scan --json
rook check --json
```

## Supported Formats

| Manager | File |
|---------|------|
| **npm** | `package.json` |
| **Maven** | `pom.xml` |
| **Go** | `go.mod` |
| **pip** | `requirements.txt` |

## License

MIT
