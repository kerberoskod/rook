# Rook — Dependency Migration CLI

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

**Rook** scans your project files, checks dependency versions against registries, and helps you update them — all from the terminal.

## Installation

```bash
# Ensure you have Go 1.23+
go install github.com/kerberoskod/rook@main
```

Or build from source:

```bash
git clone https://github.com/kerberoskod/rook.git
cd rook
go build -o rook .
```

## Usage

```bash
# Scan dependencies in the current directory
rook scan

# Scan a specific project
rook scan --path /path/to/project

# Check for outdated dependencies
rook check

# Exit with error if any are outdated
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
| **Cargo (Rust)** | `Cargo.toml` |
| **Pubspec (Dart/Flutter)** | `pubspec.yaml` |

## Testing

```bash
go test ./... -v -count=1
```

## License

MIT
