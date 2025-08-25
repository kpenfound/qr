package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"
	"golang.org/x/term"
)

const (
	// Unicode block characters for QR code rendering
	fullBlock  = "██"
	emptyBlock = "  "
)

type Config struct {
	text       string
	size       int
	outputFile string
	quiet      bool
	border     int
}

func main() {
	config := parseFlags()

	if config.text == "" {
		fmt.Fprintf(os.Stderr, "Error: text to encode is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if err := generateQR(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.text, "text", "", "Text to encode in QR code (required)")
	flag.StringVar(&config.text, "t", "", "Text to encode in QR code (shorthand)")
	flag.IntVar(&config.size, "size", 0, "Size scale 1-10 (0 for auto-detect, 1=smallest, 10=largest)")
	flag.IntVar(&config.size, "s", 0, "Size scale 1-10 (shorthand)")
	flag.StringVar(&config.outputFile, "output", "", "Output file (default: stdout)")
	flag.StringVar(&config.outputFile, "o", "", "Output file (shorthand)")
	flag.BoolVar(&config.quiet, "quiet", false, "Suppress extra output")
	flag.BoolVar(&config.quiet, "q", false, "Suppress extra output (shorthand)")
	flag.IntVar(&config.border, "border", 2, "Border size around QR code")
	flag.IntVar(&config.border, "b", 2, "Border size around QR code (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Generate QR codes in the terminal or save to file.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -text \"Hello, World!\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t \"https://example.com\" -s 5\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t \"Save to file\" -o qr.png\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  echo \"Pipe input\" | %s -t -\n", os.Args[0])
	}

	flag.Parse()

	// Handle reading from stdin if text is "-"
	if config.text == "-" {
		var input strings.Builder
		buffer := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buffer)
			if n > 0 {
				input.Write(buffer[:n])
			}
			if err != nil {
				break
			}
		}
		config.text = strings.TrimSpace(input.String())
	}

	return config
}

func generateQR(config Config) error {
	// If output file is specified, save as PNG
	if config.outputFile != "" {
		// Create QR code with proper border configuration
		qr, err := qrcode.New(config.text, qrcode.Medium)
		if err != nil {
			return fmt.Errorf("failed to generate QR code: %v", err)
		}

		// Apply border configuration
		if config.border <= 0 {
			qr.DisableBorder = true
		} else {
			// Set custom border size if supported by the library
			// Note: go-qrcode doesn't expose border width directly for PNG,
			// so we work with the default border when enabled
			qr.DisableBorder = false
		}

		return qr.WriteFile(256, config.outputFile)
	}

	// Generate QR code for terminal display
	qr, err := qrcode.New(config.text, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Always generate QR code without border for terminal display
	// We'll handle border manually in renderQRToTerminal
	qr.DisableBorder = true

	// Convert user-friendly size to actual QR dimensions
	var size int
	if config.size == 0 || config.size < 1 || config.size > 10 {
		// Use auto-detection for 0 or invalid sizes
		size = calculateOptimalSize()
	} else {
		// Use the mapped size for valid inputs (1-10)
		size = convertSizeScale(config.size)
	}

	// Get the QR code bitmap
	bitmap := qr.Bitmap()

	// Check if we're in a TTY environment
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	if !config.quiet && isTTY {
		fmt.Printf("QR Code for: %s\n", config.text)
		fmt.Println()
	}

	// Render QR code to terminal
	renderQRToTerminal(bitmap, isTTY, size, config.border)

	if !config.quiet && isTTY {
		fmt.Println()
	}

	return nil
}

func calculateOptimalSize() int {
	// Try to detect terminal size
	if term.IsTerminal(int(os.Stdout.Fd())) {
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			// Use smaller dimension and account for borders and text
			// Each QR module takes 2 characters width, so divide by 2
			maxSize := min(width/2-4, height-6)
			if maxSize > 10 && maxSize < 50 {
				return maxSize
			}
		}
	}

	// Default size for no-tty or if detection fails
	return 25
}

func renderQRToTerminal(matrix [][]bool, isTTY bool, targetSize int, border int) {
	// Scale the matrix if needed
	scaledMatrix := scaleMatrix(matrix, targetSize)

	// Add border if specified
	if border > 0 {
		scaledMatrix = addBorder(scaledMatrix, border)
	}

	for _, row := range scaledMatrix {
		var line strings.Builder
		for _, module := range row {
			if module {
				// Black module
				if isTTY {
					line.WriteString(fullBlock)
				} else {
					line.WriteString("██")
				}
			} else {
				// White module
				if isTTY {
					line.WriteString(emptyBlock)
				} else {
					line.WriteString("  ")
				}
			}
		}
		fmt.Println(line.String())
	}
}

func scaleMatrix(matrix [][]bool, targetSize int) [][]bool {
	// If targetSize is 0 or matrix is empty, return original
	if targetSize <= 0 || len(matrix) == 0 {
		return matrix
	}

	originalSize := len(matrix)

	// If target size matches original, return as-is
	if targetSize == originalSize {
		return matrix
	}

	// Calculate scale factor
	scale := float64(targetSize) / float64(originalSize)

	// Create scaled matrix
	scaled := make([][]bool, targetSize)
	for i := range scaled {
		scaled[i] = make([]bool, targetSize)
	}

	// Fill scaled matrix
	for y := 0; y < targetSize; y++ {
		for x := 0; x < targetSize; x++ {
			// Map back to original coordinates
			origY := int(float64(y) / scale)
			origX := int(float64(x) / scale)

			// Ensure we don't go out of bounds
			if origY >= originalSize {
				origY = originalSize - 1
			}
			if origX >= originalSize {
				origX = originalSize - 1
			}

			scaled[y][x] = matrix[origY][origX]
		}
	}

	return scaled
}

func addBorder(matrix [][]bool, borderSize int) [][]bool {
	if borderSize <= 0 || len(matrix) == 0 {
		return matrix
	}

	originalSize := len(matrix)
	newSize := originalSize + (borderSize * 2)

	// Create new matrix with border
	borderedMatrix := make([][]bool, newSize)
	for i := range borderedMatrix {
		borderedMatrix[i] = make([]bool, newSize)
	}

	// Copy original matrix to center, leaving borders as false (white)
	for y := 0; y < originalSize; y++ {
		for x := 0; x < originalSize; x++ {
			borderedMatrix[y+borderSize][x+borderSize] = matrix[y][x]
		}
	}

	return borderedMatrix
}

func convertSizeScale(userSize int) int {
	// Map user-friendly size (1-10) to valid QR code dimensions
	// These sizes ensure the QR code structure remains intact
	// Note: Starting from Version 2 QR (25) since Version 1 (21) doesn't produce valid QR codes
	if userSize < 1 || userSize > 10 {
		return 25 // Default fallback
	}

	sizeMap := []int{
		25, // 1
		29, // 2
		33, // 3
		37, // 4
		41, // 5
		45, // 6
		49, // 7
		53, // 8
		57, // 9
		61, // 10
	}

	return sizeMap[userSize-1]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
