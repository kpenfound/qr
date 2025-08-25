# Contributing to QR Code Generator

Thank you for your interest in contributing to this QR code generator project! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.24.1 or later
- Git

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/qr.git
   cd qr
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build the project:
   ```bash
   go build -o qr
   ```

## Testing

We maintain comprehensive tests to ensure code quality and prevent regressions. Please run tests before submitting any changes.

### Running Tests

```bash
# Run all tests with verbose output
go test -v

# Run tests with coverage report
go test -cover

# Run tests with detailed coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run only specific test functions
go test -run TestScaleMatrix -v

# Run benchmark tests
go test -bench=.

# Run benchmarks with memory allocation stats
go test -bench=. -benchmem
```

### Test Structure

Our test suite includes:

- **Unit Tests**: Test individual functions in isolation
- **Integration Tests**: Test complete workflows end-to-end
- **Benchmark Tests**: Performance testing for critical functions
- **Table-driven Tests**: Comprehensive coverage of edge cases

### Writing Tests

When adding new functionality, please include appropriate tests:

1. **Unit tests** for new functions
2. **Table-driven tests** for functions with multiple scenarios
3. **Error handling tests** for edge cases
4. **Benchmark tests** for performance-critical code

#### Example Test Structure

```go
func TestNewFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := NewFunction(tt.input)

            if tt.wantErr && err == nil {
                t.Error("Expected error but got none")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("Unexpected error: %v", err)
            }
            if result != tt.expected {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

## Code Style and Standards

### Go Standards

- Follow standard Go formatting (`go fmt`)
- Use `go vet` to catch common issues
- Follow Go naming conventions
- Add comments for exported functions and types

### Code Quality Tools

Run these before submitting:

```bash
# Format code
go fmt ./...

# Run static analysis
go vet ./...

# Run tests
go test -v ./...

# Check for common issues (if you have golangci-lint installed)
golangci-lint run
```

## Submitting Changes

### Before You Submit

1. **Run all tests**: Ensure `go test -v` passes
2. **Check code formatting**: Run `go fmt ./...`
3. **Run static analysis**: Run `go vet ./...`
4. **Test manually**: Build and test the binary with various inputs
5. **Update documentation**: Update README.md if adding new features

### Pull Request Process

1. **Create a feature branch**: `git checkout -b feature/your-feature-name`
2. **Make your changes** with appropriate tests
3. **Commit with clear messages**:
   ```
   feat: add new QR code size validation

   - Add validation for QR code size limits
   - Include tests for edge cases
   - Update help text with size constraints
   ```
4. **Push to your fork**: `git push origin feature/your-feature-name`
5. **Create a Pull Request** with:
   - Clear description of changes
   - Test results
   - Any breaking changes noted

### Commit Message Format

Use conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for adding tests
- `refactor:` for code refactoring
- `perf:` for performance improvements

## Project Structure

```
qr/
├── main.go          # Main application logic
├── main_test.go     # Comprehensive test suite
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
├── .gitignore       # Git ignore patterns
├── README.md        # Project documentation
└── CONTRIBUTING.md  # This file
```

## Feature Areas

### Current Functionality

- QR code generation from text input
- Terminal output with Unicode blocks
- PNG file output
- Configurable size and borders
- Stdin input support
- Terminal size auto-detection

### Areas for Contribution

- **Additional output formats** (SVG, ASCII art modes)
- **Error correction level options**
- **Color customization** for terminal output
- **Batch processing** capabilities
- **Configuration file support**
- **Performance optimizations**
- **Cross-platform terminal detection improvements**

Thank you for contributing to make this QR code generator better!
