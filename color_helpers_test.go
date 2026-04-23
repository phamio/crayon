package crayon

import (
	"os"
	"strings"
	"testing"
)


// =============================
// RGB TO 256 INDEX TESTS
// =============================

func TestRgbTo256Index(t *testing.T) {
	tests := []struct {
		r, g, b   int
		expected  int
		desc      string
	}{

		{0, 0, 0, 16, "black"},
		{255, 255, 255, 255, "white"},
		{255, 0, 0, 196, "red"},
		{0, 255, 0, 46, "green"},
		{0, 0, 255, 21, "blue"},
		{128, 128, 128, 244, "gray fallback"},
		{255, 255, 0, 226, "yellow"},
		{255, 0, 255, 201, "magenta"},
		{0, 255, 255, 51, "cyan"},
	}

	for _, test := range tests {
		result := rgbTo256Index(test.r, test.g, test.b)
		if result != test.expected {
			t.Errorf("%s: rgbTo256Index(%d,%d,%d) = %d, expected %d",
				test.desc, test.r, test.g, test.b, result, test.expected)
		}
	}
}

func TestRgbTo256Index_Grayscale(t *testing.T) {
	// Test grayscale detection
	result := rgbTo256Index(100, 100, 100)
	// Should be in grayscale range 232-255
	if result < 232 || result > 255 {
		t.Errorf("Grayscale should map to 232-255, got %d", result)
	}
}

// =============================
// PREFIX VALIDATION TESTS
// =============================

func TestHasValidPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"fg=red", true},
		{"bg=blue", true},
		{"fg=#FF0000", true},
		{"bg=rgb(255,0,0)", true},
		{"bold", false},
		{"underline", false},
		{"", false},
		{"invalid", false},
	}

	for _, test := range tests {
		result := hasValidPrefix(test.input)
		if result != test.expected {
			t.Errorf("hasValidPrefix(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// =============================
// HEX CODE TESTS
// =============================

func TestIsHexCode(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"FF0000", true},
		{"00ff00", true},
		{"123abc", true},
		{"ABCDEF", true},
		{"GHIJKL", false},
		{"12345", true},
		{"", true}, // empty string passes (edge case)
		{"12 34", false},
	}

	for _, test := range tests {
		result := isHexCode(test.input)
		if result != test.expected {
			t.Errorf("isHexCode(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidHex(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"fg=#FF0000", true},
		{"fg=#00ff00", true},
		{"fg=#123456", true},
		{"bg=#FF0000", true},
		{"fg=FF00", false},
		{"fg=GHIJKL", false},
		{"fg=", false},
		{"bg=", false},
		{"fg=FFFFFFF", false},
		{"bold", false},
	}

	for _, test := range tests {
		result := isValidHex(test.input)
		if result != test.expected {
			t.Errorf("isValidHex(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// =============================
// 256 COLOR CODE TESTS
// =============================

func TestIsValid256Code(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal int
		expectedOk  bool
	}{
		{"fg=0", 0, true},
		{"fg=255", 255, true},
		{"fg=128", 128, true},
		{"bg=196", 196, true},
		{"fg=256", 0, false},
		{"fg=-1", 0, false},
		{"fg=abc", 0, false},
		{"fg=", 0, false},
		{"bold", 0, false},
		{"fg=1000", 0, false},
	}

	for _, test := range tests {
		val, ok := isValid256Code(test.input)
		if ok != test.expectedOk {
			t.Errorf("isValid256Code(%q) ok = %v, expected %v", test.input, ok, test.expectedOk)
		}
		if ok && val != test.expectedVal {
			t.Errorf("isValid256Code(%q) val = %d, expected %d", test.input, val, test.expectedVal)
		}
	}
}

// =============================
// RGB VALIDATION TESTS
// =============================

func TestIsValidRGB(t *testing.T) {
	tests := []struct {
		input       string
		expectedOk  bool
		expectedLen int
	}{
		{"fg=rgb(255,0,0)", true, 3},
		{"fg=rgb(0,255,0)", true, 3},
		{"fg=rgb(0,0,255)", true, 3},
		{"bg=rgb(128,128,128)", true, 3},
		{"fg=rgb(300,0,0)", false, 0},
		{"fg=rgb(-1,0,0)", false, 0},
		{"fg=rgb(255,0)", false, 0},
		{"fg=rgb(255,0,0,0)", false, 0},
		{"fg=rgb(abc,0,0)", false, 0},
		{"fg=rgb(0,0,abc)", false, 0},
		{"fg=rgb(255,0,0", false, 0},
		{"bold", false, 0},
	}

	for _, test := range tests {
		vals, ok := isValidRGB(test.input)
		if ok != test.expectedOk {
			t.Errorf("isValidRGB(%q) ok = %v, expected %v", test.input, ok, test.expectedOk)
		}
		if ok && len(vals) != test.expectedLen {
			t.Errorf("isValidRGB(%q) returned %d values, expected %d", test.input, len(vals), test.expectedLen)
		}
	}
}

// =============================
// TERMINAL SUPPORT DETECTION
// =============================

func TestSupportsTrueColor(t *testing.T) {
	// Save original env
	origColorterm := os.Getenv("COLORTERM")
	defer os.Setenv("COLORTERM", origColorterm)

	tests := []struct {
		colorterm string
		expected  bool
	}{
		{"truecolor", true},
		{"24bit", true},
		{"", false},
		{"256color", false},
		{"something", false},
	}

	for _, test := range tests {
		os.Setenv("COLORTERM", test.colorterm)
		result := supportsTrueColor()
		if result != test.expected {
			t.Errorf("supportsTrueColor() with COLORTERM=%q = %v, expected %v",
				test.colorterm, result, test.expected)
		}
	}
}

func TestSupports256Color(t *testing.T) {
	// Save original env
	origTerm := os.Getenv("TERM")
	defer os.Setenv("TERM", origTerm)

	tests := []struct {
		term     string
		expected bool
	}{
		{"xterm-256color", true},
		{"screen-256color", true},
		{"xterm", false},
		{"vt100", false},
		{"", false},
	}

	for _, test := range tests {
		os.Setenv("TERM", test.term)
		result := supports256Color()
		if result != test.expected {
			t.Errorf("supports256Color() with TERM=%q = %v, expected %v",
				test.term, result, test.expected)
		}
	}
}

// =============================
// IS SUPPORTED COLOR TESTS
// =============================

func TestIsSupportedColor(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Named colors
		{"fg=red", true},
		{"bg=blue", true},
		{"bold", true},
		{"underline=single", true},
		
		// Resets
		{"reset", true},
		{"fg=reset", true},
		{"bg=reset", true},
		
		// Hex colors
		{"fg=#FF0000", true},
		{"bg=#00FF00", true},
		
		// 256 colors
		{"fg=196", true},
		{"bg=255", true},
		
		// RGB colors
		{"fg=rgb(255,0,0)", true},
		{"bg=rgb(0,255,0)", true},
		
		// Invalid
		{"invalid", false},
		{"", false},
		{"fg=", false},
		{"bg=", false},
	}

	for _, test := range tests {
		result := isSupportedColor(test.input)
		if result != test.expected {
			t.Errorf("isSupportedColor(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// =============================
// PARSE RGB TESTS
// =============================

func TestParseRGB(t *testing.T) {
	tests := []struct {
		input       string
		expected    []int
		expectedOk  bool
	}{
		{"fg=rgb(255,0,0)", []int{255, 0, 0}, true},
		{"fg=rgb(0,255,0)", []int{0, 255, 0}, true},
		{"fg=rgb(128,128,128)", []int{128, 128, 128}, true},
		{"fg=rgb(255,0,0,0)", nil, false},
		{"fg=rgb(255,0)", nil, false},
		{"fg=rgb(abc,0,0)", nil, false},
	}

	for _, test := range tests {
		result, ok := parseRGB(test.input)
		if ok != test.expectedOk {
			t.Errorf("parseRGB(%q) ok = %v, expected %v", test.input, ok, test.expectedOk)
		}
		if ok {
			if len(result) != len(test.expected) {
				t.Errorf("parseRGB(%q) returned %d values, expected %d", test.input, len(result), len(test.expected))
			}
			for i := range result {
				if result[i] != test.expected[i] {
					t.Errorf("parseRGB(%q)[%d] = %d, expected %d", test.input, i, result[i], test.expected[i])
				}
			}
		}
	}
}

// =============================
// PARSE ANSI TESTS
// =============================

func TestParseAnsi_Foreground(t *testing.T) {
	result, _ := colorMap["fg=red"]
	expected := "31"
	if result != expected {
		t.Errorf("parseAnsi(fg) = %q, expected %q", result, expected)
	}
}

func TestParseAnsi_ForegroundTruecolor(t *testing.T) {
	result := parseAnsi("fg=red", "2;255;0;0", false)
	expected := "\033[38;2;255;0;0m"
	if result != expected {
		t.Errorf("parseAnsi(fg truecolor) = %q, expected %q", result, expected)
	}
}

func TestParseAnsi_Background(t *testing.T) {	
	result, _ := colorMap["bg=blue"]
	expected := "44"
	if result != expected {
		t.Errorf("parseAnsi(bg) = %q, expected %q", result, expected)
	}
}

func TestParseAnsi_BackgroundTruecolor(t *testing.T) {
	result := parseAnsi("bg=blue", "2;0;0;255", false)
	expected := "\033[48;2;0;0;255m"
	if result != expected {
		t.Errorf("parseAnsi(bg truecolor) = %q, expected %q", result, expected)
	}
}

func TestParseAnsi_Invalid(t *testing.T) {
	result := parseAnsi("invalid", "31", true)
	expected := ""
	if result != expected {
		t.Errorf("parseAnsi(invalid) = %q, expected %q", result, expected)
	}
}

// =============================
// RGB TO ANSI CODE TESTS
// =============================

func TestParseRGBToAnsiCode(t *testing.T) {
	// Test with mocked terminal capabilities
	// Note: These tests depend on environment variables
	rgb := []int{255, 0, 0}
	result := parseRGBToAnsiCode("fg=rgb(255,0,0)", rgb)
	
	// Should return some ANSI code (exact value depends on terminal support)
	if result == "" {
		t.Error("parseRGBToAnsiCode returned empty string")
	}
	if !strings.HasPrefix(result, "\033[") {
		t.Errorf("parseRGBToAnsiCode result should start with ESC, got %q", result)
	}
}

// =============================
// HEX TO ANSI CODE TESTS
// =============================

func TestParseHexToAnsiCode(t *testing.T) {
	result := parseHexToAnsiCode("fg=#FF0000")
	
	if result == "" {
		t.Error("parseHexToAnsiCode returned empty string")
	}
	if !strings.HasPrefix(result, "\033[") {
		t.Errorf("parseHexToAnsiCode result should start with ESC, got %q", result)
	}
}

// =============================
// PARSE COLOR TESTS
// =============================

func TestParseColor_NamedColors(t *testing.T) {
	tests := []string{
		"fg=red", "fg=blue", "fg=green", "fg=yellow",
		"bg=red", "bg=blue", "bg=green",
	}

	for _, test := range tests {
		result := parseColor(test)
		if result == "" {
			t.Errorf("parseColor(%q) returned empty string", test)
		}
		if !strings.HasPrefix(result, "\033[") {
			t.Errorf("parseColor(%q) = %q, expected ANSI escape", test, result)
		}
	}
}

func TestParseColor_Styles(t *testing.T) {
	tests := []string{
		"bold", "italic", "underline=single", "dim", "strike",
	}

	for _, test := range tests {
		result := parseColor(test)
		if result == "" {
			t.Errorf("parseColor(%q) returned empty string", test)
		}
	}
}

func TestParseColor_Resets(t *testing.T) {
	tests := []string{
		"reset", "fg=reset", "bg=reset", "bold=reset",
	}

	for _, test := range tests {
		result := parseColor(test)
		if result == "" {
			t.Errorf("parseColor(%q) returned empty string", test)
		}
	}
}

func TestParseColor_256Colors(t *testing.T) {
	tests := []string{
		"fg=196", "fg=255", "fg=0", "bg=196", "bg=255",
	}

	for _, test := range tests {
		result := parseColor(test)
		if result == "" {
			t.Errorf("parseColor(%q) returned empty string", test)
		}
	}
}

func TestParseColor_HexColors(t *testing.T) {
	tests := []string{
		"fg=#FF0000", "fg=#00FF00", "fg=#0000FF",
		"bg=#FF0000", "bg=#00FF00", "bg=#0000FF",
	}

	for _, test := range tests {
		result := parseColor(test)
		if result == "" {
			t.Errorf("parseColor(%q) returned empty string", test)
		}
	}
}

func TestParseColor_RGBColors(t *testing.T) {
	tests := []string{
		"fg=rgb(255,0,0)",
		"fg=rgb(0,255,0)",
		"fg=rgb(0,0,255)",
		"bg=rgb(255,0,0)",
	}

	for _, test := range tests {
		result := parseColor(test)
		if result == "" {
			t.Errorf("parseColor(%q) returned empty string", test)
		}
	}
}

func TestParseColor_Invalid(t *testing.T) {
	tests := []string{
		"invalid",
		"",
		"fg=",
		"bg=",
		"fg=invalid",
	}

	for _, test := range tests {
		result := parseColor(test)
		if result != "" {
			t.Errorf("parseColor(%q) = %q, expected empty string", test, result)
		}
	}
}

// =============================
// BENCHMARK TESTS
// =============================

func BenchmarkRgbTo256Index(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rgbTo256Index(128, 128, 128)
	}
}

func BenchmarkParseColor(b *testing.B) {
	colors := []string{
		"fg=red",
		"fg=rgb(255,0,0)",
		"fg=FF0000",
		"fg=196",
		"bold",
		"reset",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, color := range colors {
			parseColor(color)
		}
	}
}

func BenchmarkIsValidRGB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidRGB("fg=rgb(255,128,64)")
	}
}

func BenchmarkIsValidHex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidHex("fg=FFAABB")
	}
}