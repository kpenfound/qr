# QR Code CLI Generator

A simple and efficient command-line tool for generating QR codes in Go that works perfectly in various terminal sizes and no-tty environments.

## Features

- 🖥️ **Terminal-friendly**: Automatically detects terminal size and adapts QR code size
- 🔧 **No-TTY support**: Works seamlessly when piped or redirected
- 📱 **Multiple output formats**: Terminal display or PNG file output
- ⚡ **Fast and lightweight**: Built with Go for optimal performance
- 🎯 **Flexible sizing**: Auto-detect or user-friendly size scale (1-10)
- 🎨 **Border control**: Configurable border sizes

## Installation

```bash
go install github.com/kpenfound/qr@latest
```

## Usage

### Basic Usage

```bash
# Generate QR code for text
qr -text "Hello, World!"

# Using shorthand flag
qr -t "https://example.com"
```

### Advanced Options

```bash
# Specify size manually (1-10 scale)
qr -t "Custom size" -s 5

# Save to PNG file
qr -t "Save to file" -o qrcode.png

# Remove border
qr -t "No border" -b 0

# Quiet mode (suppress extra output)
qr -t "Quiet mode" -q
```

### Pipe Input

```bash
# Read from stdin
echo "Pipe input" | qr -t -

# Chain with other commands
curl -s https://api.github.com/users/octocat | jq -r .html_url | qr -t -
```


## Command Line Options

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--text` | `-t` | *required* | Text to encode in QR code |
| `--size` | `-s` | `0` (auto) | Size scale 1-10 (0 for auto-detect, 1=smallest, 10=largest) |
| `--output` | `-o` | stdout | Output PNG file path |
| `--quiet` | `-q` | `false` | Suppress extra output |
| `--border` | `-b` | `2` | Border size (0 to disable) |

### Size Options

The size parameter uses a user-friendly scale that maps to valid QR code dimensions:

| Size | Dimensions | Description |
|------|------------|-------------|
| `1` | 21×21 | Smallest valid QR code |
| `2` | 25×25 | Very small |
| `3` | 29×29 | Small |
| `4` | 33×33 | Small-medium |
| `5` | 37×37 | Medium (good default) |
| `6` | 41×41 | Medium-large |
| `7` | 45×45 | Large |
| `8` | 49×49 | Very large |
| `9` | 53×53 | Extra large |
| `10` | 57×57 | Largest |

> **Note**: All sizes produce valid, scannable QR codes. The dimensions correspond to official QR code version standards.

## Examples

### Terminal Display
```bash
qr -t "https://github.com"
```

### Small QR Code
```bash
qr -t "Small" -s 1
```

### Medium QR Code
```bash
qr -t "Medium" -s 5
```

### Large QR Code without Border
```bash
qr -t "Large no border" -s 10 -b 0
```

### Save WiFi Credentials
```bash
qr -t "WIFI:T:WPA;S:MyNetwork;P:MyPassword;;" -o wifi.png
```

### Generate from File Contents
```bash
cat url_list.txt | head -1 | qr -t - -o first_url.png
```

## License

MIT License
