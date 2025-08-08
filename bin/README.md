# Bin Directory

This directory contains development builds and binaries for the Open Workbench Platform.

## Contents

- **Development builds**: Locally built binaries for testing and development
- **Cross-platform builds**: Binaries built for different operating systems and architectures

## Usage

### Development

```bash
# Build for current platform
make build

# Run the development build
./bin/om --help
```

### Cross-Platform Builds

```bash
# Build for all platforms
make build-all

# Available binaries:
# - om-darwin-amd64      (macOS Intel)
# - om-darwin-arm64      (macOS Apple Silicon)
# - om-linux-amd64       (Linux)
# - om-windows-amd64.exe (Windows AMD64)
# - om-windows-arm64.exe (Windows ARM64)
```

## Notes

- This directory is not ignored by git to allow for development builds
- Production releases are handled by GoReleaser and published to GitHub Releases
- Always use the latest release from GitHub for production use

## Building

See the main [README](../README.md) and [Contributing Guide](../docs/CONTRIBUTING.md) for detailed build instructions.
