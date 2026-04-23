package crayon

import (
	"os"
	"testing"
	"strings"

)

// =============================
// NEW COLOR TOGGLE TESTS
// =============================

func TestNewColorToggle_Default(t *testing.T) {
	// Save original env
	origNoColor := os.Getenv("NO_COLOR")
	origTerm := os.Getenv("TERM")
	defer func() {
		os.Setenv("NO_COLOR", origNoColor)
		os.Setenv("TERM", origTerm)
	}()

	// Test with TTY-like environment (no NO_COLOR, TERM set)
	os.Unsetenv("NO_COLOR")
	os.Setenv("TERM", "xterm-256color")
	
	toggle := NewColorToggle()
	if toggle == nil {
		t.Fatal("NewColorToggle() returned nil")
	}
	// Note: EnableColor depends on term.IsTerminal which is hard to mock
	// So we just check that toggle exists and has the field
}

func TestNewColorToggle_WithTrue(t *testing.T) {
	toggle := NewColorToggle(true)
	if toggle == nil {
		t.Fatal("NewColorToggle(true) returned nil")
	}
	if !toggle.EnableColor {
		t.Error("NewColorToggle(true) should enable colors")
	}
}

func TestNewColorToggle_WithFalse(t *testing.T) {
	toggle := NewColorToggle(false)
	if toggle == nil {
		t.Fatal("NewColorToggle(false) returned nil")
	}
	if toggle.EnableColor {
		t.Error("NewColorToggle(false) should disable colors")
	}
}

func TestNewColorToggle_MultipleArgs(t *testing.T) {
	// Only first argument should matter
	toggle := NewColorToggle(true, false)
	if toggle == nil {
		t.Fatal("NewColorToggle(true, false) returned nil")
	}
	if !toggle.EnableColor {
		t.Error("NewColorToggle(true, false) should use first argument (true)")
	}
}

func TestNewColorToggle_NoArgs(t *testing.T) {
	toggle := NewColorToggle()
	if toggle == nil {
		t.Fatal("NewColorToggle() returned nil")
	}
	// Should use auto-detection (can't predict value, just check it's set)
	// EnableColor will be true or false based on environment
}

// =============================
// PARSE FUNCTION TESTS
// =============================

func TestParse_BasicText(t *testing.T) {
	template := Parse("Hello World")
	
	if template.TotalLength != 11 {
		t.Errorf("Expected TotalLength 11, got %d", template.TotalLength)
	}
	
	if len(template.Parts) != 1 {
		t.Errorf("Expected 1 part, got %d", len(template.Parts))
	}
	
	if template.Parts[0].Text != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", template.Parts[0].Text)
	}
}

func TestParse_WithColorTag(t *testing.T) {
	template := Parse("[fg=red]Hello")
	
	// Should have at least one part (the color code)
	if len(template.Parts) < 1 {
		t.Errorf("Expected at least 1 part, got %d", len(template.Parts))
	}
	
	// TotalLength should be original input length
	if template.TotalLength != 13 { // "[fg=red]Hello" = 12 chars
		t.Errorf("Expected TotalLength 13, got %d", template.TotalLength)
	}
}

func TestParse_WithPlaceholders(t *testing.T) {
	template := Parse("Hello [0], welcome [1]")
	
	if len(template.Parts) != 4 {
		t.Errorf("Expected 4 parts, got %d", len(template.Parts))
	}
	
	// Check placeholder indices
	found0 := false
	found1 := false
	for _, part := range template.Parts {
		if part.Index == 0 {
			found0 = true
		}
		if part.Index == 1 {
			found1 = true
		}
	}
	
	if !found0 {
		t.Error("Placeholder [0] not found")
	}
	if !found1 {
		t.Error("Placeholder [1] not found")
	}
}


func TestParse_ComplexTemplate(t *testing.T) {
	template := Parse("[fg=blue bold][0:<20][reset] [fg=green][1][reset]")
	
	if template.TotalLength != 49 { // Length of the input string
		t.Errorf("Expected TotalLength 49, got %d", template.TotalLength)
	}
	
	// Should have multiple parts for colors, placeholders, and resets
	if len(template.Parts) < 5 {
		t.Errorf("Expected at least 5 parts, got %d", len(template.Parts))
	}
}

func TestParse_EmptyString(t *testing.T) {
	template := Parse("")
	
	if template.TotalLength != 0 {
		t.Errorf("Expected TotalLength 0, got %d", template.TotalLength)
	}
	
	if len(template.Parts) != 0 {
		t.Errorf("Expected 0 parts, got %d", len(template.Parts))
	}
}

func TestParse_OnlyBrackets(t *testing.T) {
	template := Parse("[]")
	
	if template.TotalLength != 2 {
		t.Errorf("Expected TotalLength 2, got %d", template.TotalLength)
	}
	
	// Empty brackets should be treated as literal
	foundLiteral := false
	for _, part := range template.Parts {
		if part.Text == "[]" {
			foundLiteral = true
			break
		}
	}
	
	if !foundLiteral {
		t.Error("Empty brackets '[]' should be treated as literal")
	}
}

// =============================
// COLOR TOGGLE PARSE TESTS
// =============================

func TestColorToggle_Parse_ColorsEnabled(t *testing.T) {
	toggle := NewColorToggle(true)
	template := toggle.Parse("[fg=red]Hello")
	
	// With colors enabled, color tags should be converted to ANSI codes
	// The Parts should include the ANSI code (non-empty Text)
	hasAnsiCode := false
	for _, part := range template.Parts {
		if part.Text != "" && part.Text[0] == '\033' {
			hasAnsiCode = true
			break
		}
	}
	
	if !hasAnsiCode {
		t.Error("With colors enabled, expected ANSI codes in parts")
	}
}

func TestColorToggle_Parse_ColorsDisabled(t *testing.T) {
	toggle := NewColorToggle(false)
	template := toggle.Parse("[fg=red]Hello")
	
	// With colors disabled, color tags should produce empty strings
	for _, part := range template.Parts {
		if part.Text != "" && part.Text[0] == '\033' {
			t.Error("With colors disabled, found ANSI code in parts")
		}
	}
}

func TestColorToggle_Parse_NilToggle(t *testing.T) {
	var toggle *ColorToggle = nil
	template := toggle.Parse("[fg=red]Hello")
	
	// Should handle nil gracefully (create new toggle with auto-detection)
	if template.TotalLength == 0 {
		t.Error("Parse on nil toggle should still produce a template")
	}
}

func TestColorToggle_Parse_PreservesText(t *testing.T) {
	toggle := NewColorToggle(false)
	input := "Hello [fg=red]World[reset]!"
	template := toggle.Parse(input)
	
	// The text content should be preserved even without colors
	result := template.apply()
	
	// Should contain all the text (without ANSI codes)
	if !containsString(result, "Hello") {
		t.Error("Expected 'Hello' in result")
	}
	if !containsString(result, "World") {
		t.Error("Expected 'World' in result")
	}
	if !containsString(result, "!") {
		t.Error("Expected '!' in result")
	}
}

// =============================
// INTEGRATION TESTS
// =============================

func TestIntegration_ParseAndApply(t *testing.T) {
	template := Parse("[fg=cyan][0][reset] [bold][1][reset]")
	
	result := template.apply("Hello", "World")
	
	// Should contain both words
	if !containsString(result, "Hello") {
		t.Error("Expected 'Hello' in result")
	}
	if !containsString(result, "World") {
		t.Error("Expected 'World' in result")
	}
}

func TestIntegration_ParseWithPadding(t *testing.T) {
	template := Parse("[0:<20] is [1:>10]")
	
	result := template.apply("Left", "Right")
	
	// Should have padding
	if len(result) < 30 {
		t.Errorf("Expected padded result length >= 30, got %d", len(result))
	}
}

func TestIntegration_MultipleParses(t *testing.T) {
	template1 := Parse("[fg=red]Error: [0][reset]")
	template2 := Parse("[fg=green]Success: [0][reset]")
	
	result1 := template1.apply("File not found")
	result2 := template2.apply("Operation complete")
	
	if !containsString(result1, "Error") {
		t.Error("Template1 missing 'Error'")
	}
	if !containsString(result2, "Success") {
		t.Error("Template2 missing 'Success'")
	}
}

func TestIntegration_ParseComplexStyles(t *testing.T) {
	template := Parse("[bold underline=single fg=rgb(255,0,0)][0][reset]")
	
	result := template.apply("Important")
	
	if result == "" {
		t.Error("Expected non-empty result")
	}
	// Should contain ANSI codes (unless colors disabled)
	if !containsString(result, "Important") {
		t.Error("Expected 'Important' in result")
	}
}

// =============================
// EDGE CASE TESTS
// =============================

func TestParse_UnclosedBracket(t *testing.T) {
	template := Parse("[fg=red Unclosed bracket")
	
	// Unclosed bracket should be treated as literal text
	foundLiteral := false
	for _, part := range template.Parts {
		if strings.Contains(part.Text, "[fg=red Unclosed bracket") {
			foundLiteral = true
			break
		}
	}
	if !foundLiteral {
		t.Error("Unclosed bracket should be treated as literal text")
	}
}

func TestParse_NestedBrackets(t *testing.T) {
	template := Parse("[[fg=red]]")
	
	// Colors should still be parsed in double brackets
	parsedColor := false
	for _, part := range template.Parts {
		if part.Text == "[" {
			parsedColor = true
			break
		}
	}
	
	if !parsedColor {
		t.Error("[[fg=red]] should be parsed as color and show empty brackets")
	}
}

func TestParse_WhitespaceOnly(t *testing.T) {
	template := Parse("   ")
	
	if template.TotalLength != 3 {
		t.Errorf("Expected TotalLength 3, got %d", template.TotalLength)
	}
	
	if len(template.Parts) != 1 {
		t.Errorf("Expected 1 part, got %d", len(template.Parts))
	}
	
	if template.Parts[0].Text != "   " {
		t.Errorf("Expected '   ', got '%s'", template.Parts[0].Text)
	}
}

func TestParse_SpecialCharacters(t *testing.T) {
	template := Parse("Hello\nWorld\t![0]\r\n")
	
	if template.TotalLength != 18 {
		t.Errorf("Expected TotalLength 18, got %d", template.TotalLength)
	}
	
	// Should preserve special characters
	result := template.apply()
	if !containsString(result, "\n") {
		t.Error("Expected newline preserved")
	}
	if !containsString(result, "\t") {
		t.Error("Expected tab preserved")
	}
}

// =============================
// BENCHMARK TESTS
// =============================

func BenchmarkParse_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("Hello World")
	}
}

func BenchmarkParse_WithColors(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("[fg=red]Hello [bold]World[reset]")
	}
}

func BenchmarkParse_WithPlaceholders(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("[fg=cyan][0][reset] [fg=yellow][1][reset]")
	}
}

func BenchmarkParse_ComplexTemplate(b *testing.B) {
	tmpl := "[fg=blue bold][0:<20][reset] [fg=green][1:>10][reset] [fg=red]![reset]"
	for i := 0; i < b.N; i++ {
		Parse(tmpl)
	}
}

func BenchmarkColorToggle_Parse(b *testing.B) {
	toggle := NewColorToggle(true)
	tmpl := "[fg=red]Hello [bold]World[reset]"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		toggle.Parse(tmpl)
	}
}

// =============================
// HELPER FUNCTIONS
// =============================

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// =============================
// TABLE DRIVEN TESTS
// =============================

func TestParse_VariousInputs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		minParts int
	}{
		{"Empty", "", 0},
		{"Plain text", "Hello", 1},
		{"Single color", "[fg=red]text", 2},
		{"Multiple colors", "[fg=red][bg=blue]text", 3},
		{"Placeholder", "[0]", 1},
		{"Padded placeholder", "[0:<20]", 1},
		{"Escape", "[<fg=red>]", 1},
		{"Mixed", "Start [fg=red][0][reset] End", 4},
		{"Complex", "[bold fg=cyan][0:<20][reset] [fg=yellow][1][reset]!", 6},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			template := Parse(test.input)
			
			if len(template.Parts) < test.minParts {
				t.Errorf("Expected at least %d parts, got %d", test.minParts, len(template.Parts))
			}
			
			if template.TotalLength != len(test.input) {
				t.Errorf("TotalLength = %d, expected %d", template.TotalLength, len(test.input))
			}
		})
	}
}

func TestNewColorToggle_VariousArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []bool
		expected *bool // nil means don't check specific value
	}{
		{"No args", []bool{}, nil},
		{"True", []bool{true}, boolPtr(true)},
		{"False", []bool{false}, boolPtr(false)},
		{"Multiple", []bool{true, false}, boolPtr(true)},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var toggle *ColorToggle
			if len(test.args) == 0 {
				toggle = NewColorToggle()
			} else {
				toggle = NewColorToggle(test.args...)
			}
			
			if toggle == nil {
				t.Fatal("Toggle should not be nil")
			}
			
			if test.expected != nil {
				if toggle.EnableColor != *test.expected {
					t.Errorf("EnableColor = %v, expected %v", toggle.EnableColor, *test.expected)
				}
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}