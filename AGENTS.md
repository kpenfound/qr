# AGENTS.md - QR Code CLI Generator

This file contains essential information for AI agents working with this QR Code CLI Generator project.

## Project Overview

A simple and efficient command-line tool for generating QR codes in Go that works perfectly in various terminal sizes and no-tty environments.

**Key Features:**
- Terminal-friendly QR code display using Unicode blocks
- PNG file output capability
- Automatic terminal size detection
- No-TTY environment support (pipes, redirects)
- Configurable size scales (1-10) and borders
- Stdin input support

## Key Files and Purpose

- **`main.go`** - Main application logic containing all core functionality
  - CLI argument parsing (`parseFlags()`)
  - QR code generation (`generateQR()`)
  - Terminal rendering (`renderQRToTerminal()`)
  - Matrix scaling and border functions
  - Terminal detection and sizing logic
- **`main_test.go`** - Comprehensive test suite
- **`README.md`** - Comprehensive user documentation with examples
- **`CONTRIBUTING.md`** - Developer guidelines, testing procedures, code standards

## Build and Test

### Building
```bash
# Build for current platform
go build -o qr
```

### Install new dependencies
```bash
# Use the go tool to manage dependencies
go get <dependency_name>
```

### Testing
```bash
# Run all tests with verbose output
go test -v
```

### Code Quality
```bash
# Format code
go fmt ./...
```

## Usage Examples

```bash
# Basic usage
qr -text "Hello, World!"
qr -t "https://example.com"

# Size and output control
qr -t "Custom size" -s 5
qr -t "Save to file" -o qrcode.png
qr -t "No border" -b 0

# Pipe input
echo "Pipe input" | qr -t -
curl -s https://api.github.com/users/octocat | jq -r .html_url | qr -t -
```

## Development Guidelines

### Code Standards
- Always use `go fmt` to format generated code
- Always run tests to know your task is complete
- Always clean up test artifacts like pngs if you create them while testing
- Update this AGENTS.md file when missing critical information
- Update the README.md when you change how the tool can be used
- Update the CONTRIBUTING.md when you change the project architecture
- Add new test cases for new functions

**❌ NOT PERMITTED:**
- Never try to commit your code
- Never remove existing functionality unless specifically asked
- Never reduce test coverage
- Never break backward compatibility
- Never leave the projects root directory
