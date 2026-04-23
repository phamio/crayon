// Package crayon provides a comprehensive terminal styling library for Go,
// supporting multiple color formats, styles, and automatic fallback across
// different terminal capabilities.
//
// # Overview
//
// Crayon is a Go port of the Spectra color library, offering a template-based
// approach to terminal styling. It supports named colors, hex codes, RGB values,
// 256-color palette, and automatic truecolor detection with intelligent fallback
// chains.
//
// # Key Features
//
//   - Multiple Color Systems: Named colors, hex codes, RGB, 256-color palette
//   - Full Color Fallback Chain: Automatic downsampling across all color levels
//     (truecolor → 256 color palette → ANSI 16 colors)
//   - TrueColor Detection: Automatic detection of terminal truecolor support
//   - Simple Template System: Easy-to-use templates with placeholders
//   - Comprehensive Styles: Bold, italic, underline, blink, reverse, hidden, strike-through
//   - Granular Resets: Individual and full reset codes for precise control
//   - Escape System: Texts in [] that aren't colors/styles are left as-is
//   - Inline Padding: Left and right alignment directly on placeholders
//   - Cross-Platform: Full support for Windows Terminal, cmd, PowerShell, Linux, and macOS
//   - Color Toggling: Respects NO_COLOR environment variable and TTY detection
//
// # Quick Start
//
//	package main
//
//	import (
//	    "github.com/phamio/crayon"
//	)
//
//	func main() {
//	    // Basic colors
//	    crayon.Parse("[fg=blue bg=yellow]Hello World![reset]").Println()
//
//	    // RGB colors
//	    crayon.Parse("[fg=rgb(255,105,180)]Hot pink text[reset]").Println()
//
//	    // Hex colors
//	    crayon.Parse("[fg=#FF5733]Orange text[reset]").Println()
//
//	    // 256-color palette
//	    crayon.Parse("[fg=214]Orange from palette[reset]").Println()
//
//	    // Styles
//	    crayon.Parse("[bold underline=single]Bold and underlined[reset]").Println()
//	}
//
// # Template System
//
// Crayon uses a template-first approach for optimal performance. Parse templates
// once and reuse them with different data:
//
//	// Parse once
//	errorTemplate := crayon.Parse("[fg=red bold]Error: [0][reset]")
//	warningTemplate := crayon.Parse("[fg=yellow]Warning: [0][reset]")
//
//	// Reuse multiple times
//	errorTemplate.Println("File not found")
//	errorTemplate.Println("Permission denied")
//	warningTemplate.Println("Disk space low")
//
// # Placeholders and Padding
//
// Placeholders [0] through [999] can be used for dynamic content. Padding  is applied inline:
//
//	// Left align with width 20
//	template := crayon.Parse("[0:<20]")
//
//	// Right align with width 10
//	template := crayon.Parse("[0:>10]")
//
//	// Combined with colors
//	row := crayon.Parse("[fg=cyan][0:<20][fg=yellow][1:>10][reset]")
//	row.Println("Name", "Score")
//
//
//
// # Color Toggling
//
// Crayon automatically respects NO_COLOR environment variable and TTY detection.
// Manual control is also available:
//
//	// Auto-detect (default)
//	toggle := crayon.NewColorToggle()
//
//	// Force colors on/off
//	forceOn := crayon.NewColorToggle(true)
//	forceOff := crayon.NewColorToggle(false)
//
//	// Use with templates
//	template := toggle.Parse("[fg=green]Success[reset]")
//
// # Performance Best Practices
//
// Always parse templates once and reuse them:
//
//	// GOOD: Parse once, use many times
//	var itemTemplate = crayon.Parse("[fg=blue][0][reset]")
//	for _, item := range items {
//	    itemTemplate.Println(item)
//	}
//
//	// BAD: Parsing in a loop (slow)
//	for _, item := range items {
//	    crayon.Parse("[fg=blue]" + item + "[reset]").Println()
//	}
//
// # Supported Color Formats
//
// Named Colors (Foreground):
//
//	fg=black, fg=red, fg=green, fg=yellow, fg=blue
//	fg=magenta, fg=cyan, fg=white, fg=darkgray
//	fg=lred, fg=lgreen, fg=lyellow, fg=lblue
//	fg=lmagenta, fg=lcyan, fg=lwhite
//
// Named Colors (Background):
//
//	bg=black, bg=red, bg=green, bg=yellow, bg=blue
//	bg=magenta, bg=cyan, bg=white, bg=darkgray
//	bg=lred, bg=lgreen, bg=lyellow, bg=lblue
//	bg=lmagenta, bg=lcyan, bg=lwhite
//
// Text Styles:
//
//	bold, dim, italic, underline=single, underline=double
//	blink=slow, blink=fast, reverse, hidden, strike
//
// Reset Commands:
//
//	reset               // Reset all
//	fg=reset, bg=reset  // Reset colors only
//	bold=reset, italic=reset, underline=reset  // Reset specific styles
//	blink=reset, reverse=reset, hidden=reset, strike=reset
//
// Hex Colors:
//
//	fg=#RRGGBB, bg=#RRGGBB
//	Example: fg=#FF5733, bg=#00FF00
//
// RGB Colors:
//
//	fg=rgb(R,G,B), bg=rgb(R,G,B)
//	Example: fg=rgb(255,100,50), bg=rgb(0,255,0)
//
// 256-Color Palette:
//
//	fg=0-255, bg=0-255
//	Example: fg=196 (red), fg=46 (green), fg=21 (blue)
//
// # Platform Support
//
//   - Linux: Full support
//   - macOS: Full support
//   - Windows Terminal: Full support
//   - Windows cmd: Full support (with VT enabled)
//   - PowerShell: Full support
//   - Legacy Windows CMD: Limited style support
//
// # Limitations
//
//   - Terminal Dependency: Colors only work in terminals that support ANSI escape codes
//   - Style Support: Some styles (blink, double underline) may not work in all terminals
//   - Placeholder Limit: Supports [0] through [999]
//
// # Environment Variables
//
//   - NO_COLOR: When set, disables all color output
//   - COLORTERM: Used for truecolor detection (truecolor or 24bit)
//   - TERM: Used for 256-color and dumb terminal detection (contains "256color")
//
// # Examples
//
// Example 1: Colored Table
//
//	header := crayon.Parse("[bold fg=cyan][0][reset]")
//	row := crayon.Parse("[0]  [fg=yellow][1][reset]  [fg=green][2][reset]")
//
//	header.Println("USER MANAGEMENT")
//	row.Println("Alice", "admin", "active")
//	row.Println("Bob", "user", "active")
//
// Example 2: Log Formatter
//
//	logTemplate := crayon.Parse("[0] [fg=blue][1][reset]: [fg=yellow][2][reset]")
//	logTemplate.Println("[INFO]", "main", "Application started")
//	logTemplate.Println("[WARN]", "auth", "Token expiring soon")
//	logTemplate.Println("[ERROR]", "db", "Connection failed")
//
// Example 3: Progress Bar
//
//	progress := crayon.Parse("[fg=cyan][0][reset]/[fg=cyan][1][reset] [fg=green][2][reset]%%")
//	for i := 0; i <= 100; i += 10 {
//	    fmt.Printf("\r%s", progress.Sprint(i, 100, i))
//	    time.Sleep(100 * time.Millisecond)
//	}
//
// Example 4: CLI Help Output
//
//	header := crayon.Parse("[bold fg=cyan][0][reset]")
//	command := crayon.Parse("[fg=yellow][0:<25][fg=green][1][reset]")
//
//	header.Println("Available Commands:")
//	command.Println("start", "Start the application")
//	command.Println("stop", "Stop the application")
//	command.Println("status", "Check status")
//
// # Contributing
//
// Contributions are welcome! Areas needing improvement include:
//   - Windows testing across different terminals
//   - Performance optimizations
//   - Additional test coverage
//
// # License
//
// MIT License - see LICENSE file for details.
package crayon