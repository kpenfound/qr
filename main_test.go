package main

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestScaleMatrix(t *testing.T) {
	tests := []struct {
		name       string
		matrix     [][]bool
		targetSize int
		expected   [][]bool
	}{
		{
			name:       "empty matrix",
			matrix:     [][]bool{},
			targetSize: 5,
			expected:   [][]bool{},
		},
		{
			name:       "zero target size",
			matrix:     [][]bool{{true, false}, {false, true}},
			targetSize: 0,
			expected:   [][]bool{{true, false}, {false, true}},
		},
		{
			name:       "same size",
			matrix:     [][]bool{{true, false}, {false, true}},
			targetSize: 2,
			expected:   [][]bool{{true, false}, {false, true}},
		},
		{
			name:       "scale up 2x2 to 4x4",
			matrix:     [][]bool{{true, false}, {false, true}},
			targetSize: 4,
			expected: [][]bool{
				{true, true, false, false},
				{true, true, false, false},
				{false, false, true, true},
				{false, false, true, true},
			},
		},
		{
			name: "scale down 4x4 to 2x2",
			matrix: [][]bool{
				{true, true, false, false},
				{true, true, false, false},
				{false, false, true, true},
				{false, false, true, true},
			},
			targetSize: 2,
			expected:   [][]bool{{true, false}, {false, true}},
		},
		{
			name:       "single pixel matrix",
			matrix:     [][]bool{{true}},
			targetSize: 3,
			expected:   [][]bool{{true, true, true}, {true, true, true}, {true, true, true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scaleMatrix(tt.matrix, tt.targetSize)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d rows, got %d", len(tt.expected), len(result))
				return
			}

			for i := range result {
				if len(result[i]) != len(tt.expected[i]) {
					t.Errorf("Row %d: expected %d columns, got %d", i, len(tt.expected[i]), len(result[i]))
					continue
				}

				for j := range result[i] {
					if result[i][j] != tt.expected[i][j] {
						t.Errorf("Position [%d][%d]: expected %v, got %v", i, j, tt.expected[i][j], result[i][j])
					}
				}
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"a smaller", 1, 2, 1},
		{"b smaller", 3, 2, 2},
		{"equal", 5, 5, 5},
		{"negative numbers", -1, -2, -2},
		{"mixed signs", -1, 2, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := min(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("min(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCalculateOptimalSize(t *testing.T) {
	// This function depends on terminal detection, so we'll test the fallback behavior
	size := calculateOptimalSize()

	// Should return a reasonable default when terminal size can't be detected
	if size <= 0 || size > 100 {
		t.Errorf("calculateOptimalSize() returned unreasonable size: %d", size)
	}
}

func TestRenderQRToTerminal(t *testing.T) {
	tests := []struct {
		name       string
		matrix     [][]bool
		isTTY      bool
		targetSize int
	}{
		{
			name:       "simple 2x2 matrix TTY",
			matrix:     [][]bool{{true, false}, {false, true}},
			isTTY:      true,
			targetSize: 2,
		},
		{
			name:       "simple 2x2 matrix non-TTY",
			matrix:     [][]bool{{true, false}, {false, true}},
			isTTY:      false,
			targetSize: 2,
		},
		{
			name:       "empty matrix",
			matrix:     [][]bool{},
			isTTY:      true,
			targetSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			renderQRToTerminal(tt.matrix, tt.isTTY, tt.targetSize, 0)

			w.Close()
			os.Stdout = old

			buf := new(bytes.Buffer)
			buf.ReadFrom(r)
			output := buf.String()

			// Basic validation that output was generated
			if len(tt.matrix) > 0 && len(output) == 0 {
				t.Error("Expected some output for non-empty matrix")
			}
		})
	}
}

// Mock stdin for testing
func TestParseFlags_StdinInput(t *testing.T) {
	// Save original args and stdin
	oldArgs := os.Args
	oldStdin := os.Stdin

	// Create a pipe to simulate stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Write test data to stdin
	go func() {
		defer w.Close()
		w.Write([]byte("test input from stdin"))
	}()

	// Set up command line args to read from stdin
	os.Args = []string{"qr", "-text", "-"}

	// Reset flag package state
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	defer func() {
		os.Args = oldArgs
		os.Stdin = oldStdin
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	config := parseFlags()

	if config.text != "test input from stdin" {
		t.Errorf("Expected 'test input from stdin', got '%s'", config.text)
	}
}

func TestParseFlags_BasicFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Config
	}{
		{
			name: "basic text flag",
			args: []string{"qr", "-text", "hello world"},
			expected: Config{
				text:   "hello world",
				size:   0,
				border: 2,
			},
		},
		{
			name: "shorthand flags",
			args: []string{"qr", "-t", "test", "-s", "10", "-b", "1"},
			expected: Config{
				text:   "test",
				size:   10,
				border: 1,
			},
		},
		{
			name: "all flags",
			args: []string{"qr", "-text", "full test", "-size", "15", "-output", "test.png", "-quiet", "-border", "3"},
			expected: Config{
				text:       "full test",
				size:       15,
				outputFile: "test.png",
				quiet:      true,
				border:     3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original args
			oldArgs := os.Args
			defer func() {
				os.Args = oldArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			// Set test args
			os.Args = tt.args

			// Reset flag package
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			config := parseFlags()

			if config.text != tt.expected.text {
				t.Errorf("text: expected '%s', got '%s'", tt.expected.text, config.text)
			}
			if config.size != tt.expected.size {
				t.Errorf("size: expected %d, got %d", tt.expected.size, config.size)
			}
			if config.outputFile != tt.expected.outputFile {
				t.Errorf("outputFile: expected '%s', got '%s'", tt.expected.outputFile, config.outputFile)
			}
			if config.quiet != tt.expected.quiet {
				t.Errorf("quiet: expected %v, got %v", tt.expected.quiet, config.quiet)
			}
			if config.border != tt.expected.border {
				t.Errorf("border: expected %d, got %d", tt.expected.border, config.border)
			}
		})
	}
}

func TestGenerateQR_InvalidInput(t *testing.T) {
	config := Config{
		text: "", // Empty text should cause an error
	}

	err := generateQR(config)
	if err == nil {
		t.Error("Expected error for empty text, got nil")
	}
}

func TestGenerateQR_ValidInput(t *testing.T) {
	// Test with valid input but no file output (terminal output)
	config := Config{
		text:   "test message",
		size:   10,
		quiet:  true, // Quiet mode to reduce output
		border: 1,
	}

	// Capture stdout to avoid cluttering test output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := generateQR(config)

	w.Close()
	os.Stdout = old

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	if err != nil {
		t.Errorf("Expected no error for valid input, got: %v", err)
	}

	// Should produce some output
	if buf.Len() == 0 {
		t.Error("Expected some QR code output, got none")
	}
}

func TestGenerateQR_FileOutput(t *testing.T) {
	tempFile := "/tmp/test_qr.png"

	config := Config{
		text:       "test file output",
		outputFile: tempFile,
	}

	err := generateQR(config)

	// Clean up
	defer os.Remove(tempFile)

	if err != nil {
		t.Errorf("Expected no error for file output, got: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("Expected output file to be created")
	}
}

// Benchmark tests
func BenchmarkScaleMatrix(b *testing.B) {
	matrix := [][]bool{
		{true, false, true, false},
		{false, true, false, true},
		{true, false, true, false},
		{false, true, false, true},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scaleMatrix(matrix, 20)
	}
}

func BenchmarkGenerateQR(b *testing.B) {
	config := Config{
		text:   "benchmark test message",
		size:   15,
		quiet:  true,
		border: 1,
	}

	// Redirect stdout to avoid cluttering benchmark output
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateQR(config)
	}
}

func TestAddBorder(t *testing.T) {
	tests := []struct {
		name       string
		matrix     [][]bool
		borderSize int
		expected   [][]bool
	}{
		{
			name:       "no border",
			matrix:     [][]bool{{true, false}, {false, true}},
			borderSize: 0,
			expected:   [][]bool{{true, false}, {false, true}},
		},
		{
			name:       "border size 1",
			matrix:     [][]bool{{true, false}, {false, true}},
			borderSize: 1,
			expected: [][]bool{
				{false, false, false, false},
				{false, true, false, false},
				{false, false, true, false},
				{false, false, false, false},
			},
		},
		{
			name:       "empty matrix",
			matrix:     [][]bool{},
			borderSize: 2,
			expected:   [][]bool{},
		},
		{
			name:       "single pixel with border",
			matrix:     [][]bool{{true}},
			borderSize: 2,
			expected: [][]bool{
				{false, false, false, false, false},
				{false, false, false, false, false},
				{false, false, true, false, false},
				{false, false, false, false, false},
				{false, false, false, false, false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addBorder(tt.matrix, tt.borderSize)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d rows, got %d", len(tt.expected), len(result))
				return
			}

			for i := range result {
				if len(result[i]) != len(tt.expected[i]) {
					t.Errorf("Row %d: expected %d columns, got %d", i, len(tt.expected[i]), len(result[i]))
					continue
				}

				for j := range result[i] {
					if result[i][j] != tt.expected[i][j] {
						t.Errorf("Position [%d][%d]: expected %v, got %v", i, j, tt.expected[i][j], result[i][j])
					}
				}
			}
		})
	}
}

func TestConvertSizeScale(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"valid size 1", 1, 25},
		{"valid size 5", 5, 41},
		{"valid size 10", 10, 61},
		{"invalid size 0", 0, 25},
		{"invalid size 11", 11, 25},
		{"invalid size 15", 15, 25},
		{"invalid negative size", -1, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertSizeScale(tt.input)
			if result != tt.expected {
				t.Errorf("convertSizeScale(%d) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateQR_InvalidSizes(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"size too large", 15},
		{"size too small", -1},
		{"size zero", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				text:   "test message",
				size:   tt.size,
				quiet:  true,
				border: 1,
			}

			// Capture stdout to avoid cluttering test output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := generateQR(config)

			w.Close()
			os.Stdout = old

			buf := new(bytes.Buffer)
			buf.ReadFrom(r)

			if err != nil {
				t.Errorf("Expected no error for invalid size %d, got: %v", tt.size, err)
			}

			// Should produce some output (auto-detection should work)
			if buf.Len() == 0 {
				t.Error("Expected some QR code output for invalid size, got none")
			}
		})
	}
}

// Test helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(s, item) {
			return true
		}
	}
	return false
}
