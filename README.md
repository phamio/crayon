
# Crayon

A comprehensive Go library for adding colors, styles, and formatting to terminal output with support for multiple color formats and truecolor detection. [This is a Go port of the Spectra color library](https://github.com/phamio/spectra)


# Installation

```bash
go get github.com/phamio/crayon
```

# Features

- Multiple Color Systems: Named colors, hex codes, RGB, 256-color palette.
- TrueColor Detection: Automatic detection of terminal truecolor support.
- Terminal Safe: Graceful fallbacks when color not supported(truecolor-to-256-palette fallback).
- Simple API: Easy-to-use functions for text styling and coloring.
- Comprehensive Styles: Bold, italic, underline, blink, reverse, hidden, strike-through.
- Granular Resets: Individual and full reset codes for precise control.
- Escape aren't needed: Texts in [] that aren't colors/styles are left as it is.
- Inline Padding: Left and right alignment, declared directly on placeholders.
- Cross-Platform: Full color support on Windows (Windows Terminal, cmd), Linux and macOS.

---

# Core Concepts

## Template System

The library follows a template-first approach: parse color templates once with or without placeholders([0], [1], etc), then reuse them with different data to replace placeholders.
**Placeholders are like slots**

## Color Toggling

Respects the NO_COLOR environment variable and detects when output is redirected. It can be manually controlled to suit user preference.

## Padding & Alignment
When printing repeated lines of output, values rarely line up on their own.
Padding fixes that by reserving a fixed width for each value so columns stay aligned across every print call.

Padding is applied directly on placeholders

---

# Quick Start

```go
package main

import (
    "fmt"
    "github.com/phamio/crayon"
)

func main() {
    // Parse and use color codes directly
    red := crayon.ParseColor("fg=red")
    bold := crayon.ParseColor("bold")
    reset := crayon.ParseColor("reset")
    
    fmt.Printf("%sThis is red and bold!%s\n", red + bold, reset)
    
    // Or use the main functions
    crayon.Parse("[fg=blue]Hello in blue![reset]").Println()
    crayon.Parse("[bg=yellow fg=black bold]Bold black text on yellow background.[reset]").Println()


    //Or pre-parse the color template with placeholders for reuse. This is the heart of the library's performance.

    // Parse once
    template := crayon.Parse("[fg=red bold]Error: [0][reset]")

    // Reuse multiple times
    template.Println("File not found")
    template.Println("Permission denied")
    template.Println("Network timeout")

    //padding or alignment
    row := crayon.Parse("[fg=cyan bold][0:<20][fg=yellow][1:>10][reset]")
    
    row.Println("Alice", "admin")
    row.Println("Bob", "user")
    row.Println("Charlie", "guest")

    //escapes
    //escapes are enclosed in [<content>] where "content" is the word to be escaped
    escapes := crayon.Parse("[fg=cyan][<fg=red>]Hello World[reset]")
}
```
# Output

![quick_start](images/quick_start_merge.png)

---

# Complete Usage Examples

## Basic Template with Placeholders

```go
package main

import (
    "fmt"
    "github.com/phamio/crayon"
)

func main() {
    // Simple template with one placeholder
    greeting := crayon.Parse("[fg=green]Hello, [0][reset]!")
    
    greeting.Println("Alice")
    greeting.Println("Bob")
    greeting.Println("World")
    
    // Complex template with multiple placeholders
    logTemplate := crayon.Parse("[0] [fg=blue][1][reset]: [fg=yellow][2][reset]")
    
    // Different log levels
    logTemplate.Println("[INFO]", "main", "Application started")
    logTemplate.Println("[WARN]", "auth", "Token expiring soon")
    logTemplate.Println("[ERROR]", "db", "Connection failed")
}
```

# Output

![basic_index_merge](images/basic_index_merge.png)

---

## Basic Text Coloring
```go
package main

import "github.com/phamio/crayon"

func main(){
  // Simple colored text
  crayon.Parse("[fg=green]Success message![reset]").Println()
  crayon.Parse("[fg=red bold]Error: Something went wrong![reset]").Println()
  crayon.Parse("[fg=cyan italic]Info message[reset]").Println()

  // Background colors
  crayon.Parse("[bg=blue fg=white]White text on blue background[reset]").Println()
  crayon.Parse("[bg=lightgreen fg=black]Black text on light green[reset]").Println()
}
```

## Advanced Color Formats

```go
package main

import "github.com/phamio/crayon"

func main(){
// Hex colors (requires truecolor support)
crayon.Parse("[fg=#FF5733]Orange hex color[reset]").Println()
crayon.Parse("[bg=#3498db]Blue background[reset]").Println()

// RGB colors
crayon.Parse("[fg=rgb(255,105,180)]Hot pink text[reset]").Println()
crayon.Parse("[bg=rgb(50,205,50)]Lime green background[reset]").Println()

// 256-color palette
crayon.Parse("[fg=214]Orange from 256-color palette[reset]").Println()
crayon.Parse("[bg=196]Red background from palette[reset]").Println()
}
```

### Padding Example1

```go
package main

import (
    "github.com/phamio/crayon"
)

func main(){
    row := crayon.Parse("[fg=cyan bold][0:<20][fg=yellow][1:>10][reset]")
    
    row.Println("Alice", "admin")
    row.Println("Bob", "user")
    row.Println("Charlie", "guest")
}
```
# Output

![pad_row](images/pad_row_merge.png)
---

### Padding Example2

```go
package main

import (
    "fmt"
    "github.com/phamio/crayon"
)
  var(
  header = crayon.Parse("[bold fg=cyan][0][reset]")
  command  = crayon.Parse("[fg=yellow][0:<25][fg=green][1][reset]")
  flag = crayon.Parse("[fg=yellow][0][reset], [fg=yellow][1:<20] [fg=green][2][reset]")
  )

func ShowHelp() {
    header.Println("MyApp Help")
    fmt.Println()
    
    header.Println("Usage:")
    fmt.Println("  myapp [command] [options]")
    fmt.Println()
    
    header.Println("Commands:")
    command.Println("start", "Start the application")
    command.Println("stop", "Stop the application")
    command.Println("status", "Check application status")
    command.Println()
    
    header.Println("Options:")
    flag.Println("-h", "--help", "Show this help")
    flag.Println("-v", "--version", "Show version")
    flag.Println("-d", "--debug", "Enable debug mode")
}

func main(){
	ShowHelp()
}

```
# Output

![help_flag](images/help_flag_merge.png)
---


## Text Styles

```go
package main

import "github.com/phamio/crayon"

func main(){
    // Combine styles
    crayon.Parse("[bold underline=single] Bold and underlined[reset]").Println()
    crayon.Parse("[italic dim]", "Dim italic text. [italic=reset dim=reset][strike]Strikethrough  text only[reset]").Println()
    crayon.Parse("[blink=slow hidden]Slow blinking hidden text[reset]").Println()

    // Reset specific attributes
    crayon.Parse("[bold fg=blue]Blue bold text. [bold=reset]No longer bold, but still blue. [fg=reset]No color, but other styles remain[reset]").Println()
}
```

# Output

![basic_color_merge](images/basic_color_merge.png)
---

## Color Toggling

```go
package main

import (
    "fmt"
    "os"
    "github.com/phamio/crayon"
)

func main() {
    // Create color toggle - respects NO_COLOR env var and when output is redirected by default
    toggle := crayon.NewColorToggle()
    
    // Parse templates using the toggle
    successTemplate := toggle.Parse("[fg=green]✓ [0][reset]")
    errorTemplate := toggle.Parse("[fg=red]✗ [0][reset]")
    
    // These will only show colors if appropriate
    successTemplate.Println("Operation completed")
    errorTemplate.Println("Operation failed")
    
    // Manual control
    forceColors := crayon.NewColorToggle(true)   // Always show colors
    noColors := crayon.NewColorToggle(false)     // Never show colors
    
    // Use in CLI applications
    useColor := os.Getenv("NO_COLOR") == ""
    appToggle := crayon.NewColorToggle(useColor)
    
    helpTemplate := appToggle.Parse("[bold fg=cyan][0][reset] [fg=green][1][reset]")
    helpTemplate.Println("Usage:", "myapp [options]")
}
```

# Output

![color_toggle](images/color_toggle_merge.png)
---

## Advanced Template Examples

```go
package main

import (
    "fmt"
    "time"
    "github.com/phamio/crayon"
)

func main() {
    // Status indicator with conditional colors
    statusTemplate := crayon.Parse("[0] [1][reset]")
    
    items := []struct{
        name string
        status string
    }{
        {"Database", "Online"},
        {"API Server", "Offline"},
        {"Cache", "Degraded"},
    }
    
    for _, item := range items {
        var statusColor string
        switch item.status {
        case "Online":
            statusColor = "[fg=green bold]"
        case "Offline":
            statusColor = "[fg=red bold]"
        default:
            statusColor = "[fg=yellow]"
        }
        
        statusColored := crayon.Parse(statusColor + item.status).Sprint()
        statusTemplate.Println(item.name + ":", statusColored)
    }
    
    // Progress bar template
    progressTemplate := crayon.Parse("[fg=cyan][0][reset]/[fg=cyan][1][reset] [fg=green][2][reset]%")
    
    total := 100
    for i := 0; i <= total; i += 10 {
        percent := i * 100 / total
        fmt.Printf("\r%s", progressTemplate.Sprint(i, total, percent))
        time.Sleep(100 * time.Millisecond)
    }
    fmt.Println()
}
```

# Output

![adv_temp](images/adv_temp_merge.png)
---

## Building Complex UIs

```go
package main

import (
    "fmt"
    "strings"
    "github.com/phamio/crayon"
)

func main() {
    
    // Table with colored headers
    headerTemplate := crayon.Parse("[bold fg=cyan][0][reset]")
    rowTemplate := crayon.Parse("[0]  [fg=yellow][1][reset]  [fg=green][2][reset]")
    
    headerTemplate.Println(strings.Repeat("─", 40))
    headerTemplate.Println("USER MANAGEMENT")
    headerTemplate.Println(strings.Repeat("─", 40))
    
    rowTemplate.Println("Alice", "admin", "active")
    rowTemplate.Println("Bob", "user", "active")
    rowTemplate.Println("Charlie", "guest", "inactive")
    
    // Nested templates
    errorTemplate := crayon.Parse("[bold fg=red][0][reset]: [1]")
    suggestionTemplate := crayon.Parse("[fg=yellow]Suggestion: [0][reset]")
    
    errors := []struct{
        code string
        msg string
        suggestion string
    }{
        {"E001", "File not found", "Check the file path"},
        {"E002", "Permission denied", "Run with sudo or check permissions"},
        {"E003", "Out of memory", "Close other applications"},
    }
    
    for _, err := range errors {
        errorTemplate.Println(err.code, err.msg)
        fmt.Println("  " + suggestionTemplate.Sprint(err.suggestion))
        fmt.Println()
    }
}
```

# Output

![complex_ui](images/complex_ui_merge.png)
---

## Project Structure Example

```go
// file: styles/styles.go - Define your color scheme
package styles

import "github.com/phamio/crayon"

var Toggle = crayon.NewColorToggle()

var Templates = struct {
    Success  crayon.CompiledTemplate
    Error    crayon.CompiledTemplate
    Warning  crayon.CompiledTemplate
    Info     crayon.CompiledTemplate
    Header   crayon.CompiledTemplate
    Flag     crayon.CompiledTemplate
}{
    Success:  Toggle.Parse("[fg=green bold]✓ [0][reset]"),
    Error:    Toggle.Parse("[fg=red bold]✗ [0][reset]"),
    Warning:  Toggle.Parse("[fg=yellow bold]⚠ [0][reset]"),
    Info:     Toggle.Parse("[fg=blue][0][reset]"),
    Header:   Toggle.Parse("[bold fg=cyan][0][reset]"),
    Flag:     Toggle.Parse("[fg=yellow][0][reset], [fg=yellow][1][reset]: [2]"),
}
```



```go
package cmd

//file: cmd/help.go - Use the color templates
import (
    "fmt"
    "yourproject/styles"
)

func ShowHelp() {
    styles.Templates.Header.Println("MyApp Help")
    fmt.Println()
    
    styles.Templates.Header.Println("Usage:")
    fmt.Println("  myapp [command] [options]")
    fmt.Println()
    
    styles.Templates.Header.Println("Commands:")
    styles.Templates.Flag.Println("start", "", "Start the application")
    styles.Templates.Flag.Println("stop", "", "Stop the application")
    styles.Templates.Flag.Println("status", "", "Check application status")
    fmt.Println()
    
    styles.Templates.Header.Println("Options:")
    styles.Templates.Flag.Println("-h", "--help", "Show this help")
    styles.Templates.Flag.Println("-v", "--version", "Show version")
    styles.Templates.Flag.Println("-d", "--debug", "Enable debug mode")
}
```


```go
package main

//file: main.go - Main application
import (
    "fmt"
    "os"
    "yourproject/styles"
    "yourproject/cmd"
)

func main() {
    if len(os.Args) > 1 && os.Args[1] == "--help" {
        cmd.ShowHelp()
        return
    }
    
    // Use color templates throughout
    styles.Templates.Success.Println("Application started")
    
    // Process...
    
    styles.Templates.Info.Println("Processing completed")
    styles.Templates.Success.Println("All tasks finished")
}
```


## CLI Applications

```go
package main
import(
    "fmt"
    "os"
    "github.com/phamio/crayon"
)
// Best practice for CLI applications
func main() {
    // Check for --no-color flag
    noColorFlag := false
    for _, arg := range os.Args {
        if arg == "--no-color" {
            noColorFlag = true
            break
        }
    }
    
    // Respect both flag and environment variable
    useColor := !noColorFlag && os.Getenv("NO_COLOR") == ""
    
    // Create toggle
    toggle := crayon.NewColorToggle(useColor)
    
    // All templates use this toggle
    templates := struct {
        Success crayon.CompiledTemplate
        Error   crayon.CompiledTemplate
        Header  crayon.CompiledTemplate
    }{
        Success: toggle.Parse("[fg=green]✓ [0][reset]"),
        Error:   toggle.Parse("[fg=red]✗ [0][reset]"),
        Header:  toggle.Parse("[bold][0][reset]"),
    }
    
    // Use templates - they'll respect the toggle
    templates.Header.Println("My Application")
    templates.Success.Println("Started successfully")
    
    // If --no-color was used or NO_COLOR is set,
    // outputs will be plain text without escape codes
}
```

# Output

![cli_app](images/cli_app_merge.png)
---

## Error Handling in Templates

```go
package main

import (
    "fmt"
    "github.com/phamio/crayon"
)

func main() {
    // Template for showing validation errors
    validationTemplate := crayon.Parse("[fg=red]• [0]: [1][reset]")
    
    errors := map[string]string{
        "username": "Must be at least 3 characters",
        "email":    "Invalid email format",
        "password": "Must contain uppercase and numbers",
    }
    
    crayon.Parse("[bold fg=yellow]Validation Errors:[reset]").Println()
    for field, message := range errors {
        validationTemplate.Println(field, message)
    }
    
    // Template with conditional formatting
    scoreTemplate := crayon.Parse("[0]: [1]")
    
    scores := []struct{
        name string
        score int
    }{
        {"Alice", 95},
        {"Bob", 75},
        {"Charlie", 45},
        {"Diana", 60},
    }
    
    for _, s := range scores {
        var scoreColor string
        switch {
        case s.score >= 90:
            scoreColor = "[fg=green bold]"
        case s.score >= 70:
            scoreColor = "[fg=yellow]"
        default:
            scoreColor = "[fg=red]"
        }
        
        coloredScore := crayon.Parse(scoreColor + fmt.Sprint(s.score)+ "[reset]").Sprint()
        scoreTemplate.Println(s.name, coloredScore)
    }
}
```

# Output

![error_handle](images/err_handle_merge.png)
---

## Pattern to avoid
```go
// Good pattern
var appTemplates struct {
    Success crayon.Template
    Error   crayon.Template
}

func init() {
    toggle := crayon.NewColorToggle()
    appTemplates.Success = toggle.Parse("[fg=green] [0][reset]")
    appTemplates.Error = toggle.Parse("[fg=red] [0][reset]")
}

// Bad pattern (parsing in hot loop)
func processItems(items []string) {
    for _, item := range items {
        // DON'T DO THIS - parses every iteration!
        tmpl := crayon.Parse("[fg=blue]" + item + "[reset]")
        tmpl.Println()
    }
}
```

## Performance Comparison

```go
package main

import (
    "fmt"
    "time"
    "github.com/phamio/crayon"
)

func main() {
    const iterations = 1000000
    
    // Method 1: Parse once, apply many
    template := crayon.Parse("[bold fg=red][0][reset] [fg=green][1][reset]")
    
    start := time.Now()
    for i := 0; i < iterations; i++ {
        template.Sprint(fmt.Sprintf("Item%d", i), fmt.Sprintf("Value%d", i))
    }
    fmt.Printf("Template reuse: %v\n", time.Since(start))
    
    // Method 2: Parse every time
    start = time.Now()
    for i := 0; i < iterations; i++ {
        crayon.Parse(fmt.Sprintf("[bold fg=red]Item%d[reset] [fg=green]Value%d[reset]", i, i)).Sprint()
    }
    fmt.Printf("Parse every time: %v\n", time.Since(start))
}
```

## Performance Comparison Result

```bash
Template reuse: 1.75269893s
Parse every time: 20.433925041s
```


# Spectra Syntax Reference

## Basic Colors
**Foreground Colors**
| Command | Effect |
|---------|--------|
| `fg=black` | Black text |
| `fg=red` | Red text |
| `fg=green` | Green text |
| `fg=yellow` | Yellow text |
| `fg=blue` | Blue text |
| `fg=magenta` | Magenta text |
| `fg=cyan` | Cyan text |
| `fg=white` | White text |
| `fg=darkgray` | Dark gray text |
| `fg=lred` | Light red text |
| `fg=lgreen` | Light green text |
| `fg=lyellow` | Light yellow text |
| `fg=lblue` | Light blue text |
| `fg=lmagenta` | Light magenta text |
| `fg=lcyan` | Light cyan text |
| `fg=lwhite` | Light white text |


**Background Colors**
| Command | Effect |
|---------|--------|
| `bg=black` | Black background |
| `bg=red` | Red background |
| `bg=green` | Green background |
| `bg=yellow` | Yellow background |
| `bg=blue` | Blue background |
| `bg=magenta` | Magenta background |
| `bg=cyan` | Cyan background |
| `bg=white` | White background |
| `bg=darkgray` | Dark gray background |
| `bg=lred` | Light red background |
| `bg=lgreen` | Light green background |
| `bg=lyellow` | Light yellow background |
| `bg=lblue` | Light blue background |
| `bg=lmagenta` | Light magenta background |
| `bg=lcyan` | Light cyan background |
| `bg=lwhite` | Light white background |


## Text Styles
| Command | Effect |
|---------|--------|
| `bold` | Bold/bright text |
| `dim` | Dim/faint text |
| `italic` | Italic text |
| `underline=single` | Single underlined text |
| `underline=double` | Double underlined text |
| `blink=slow` | Slow blinking text |
| `blink=fast` | Fast blinking text |
| `reverse` | Reverse video (swap foreground and background colors) |
| `hidden` | Hidden text |
| `strike` | Strikethrough text |

## Reset Commands
| Command | Effect |
|---------|--------|
| `reset` | Reset all colors and styles |
| `fg=reset` | Reset foreground color only |
| `bg=reset` | Reset background color only |
| `bold=reset` | Reset bold style only |
| `dim=reset` | Reset dim style only |
| `italic=reset` | Reset italic style only |
| `underline=reset` | Reset underline style only |
| `blink=reset` | Reset blink style only |
| `reverse=reset` | Reset reverse style only |
| `hidden=reset` | Reset hidden style only |
| `strike=reset` | Reset strikethrough style only |


## Advanced Features
| Command | Effect |
|---------|--------|
| `fg=#RRGGBB` | Hex color for foreground |
| `bg=#RRGGBB` | Hex color for background |
| `fg=rgb(RR,GG,BB)` |RGB color for foreground |
| `bg=rgb(RR,GG,BB)` | RGB color for background |
| `fg=NNN` | 256-color palette (0-255) for foreground |
| `bg=NNN` | 256-color palette (0-255) for background |


## Padding Syntax Reference
| Syntax | Effect |
|---------|--------|
| `[a:bc]` | a=index(placeholder), b=direction, c=width|
| `[0:<20]` | placeholder=0, direction=left align, width=20 |
| `[0:>10]` | placeholder=0, direction=right align, placeholder 0, width=10 |




# Tips and Best Practices

1. Parse Once: Always parse templates at initialization, not in loops
2. Use Toggles: Respect user preferences with color toggling
3. Template Reuse: Create templates for consistent styling
4. Placeholder Limits: The current implementation supports [0] through [999]
5. Testing: Test both color and no-color outputs


# Limitations

1. Terminal Dependency: Colors only work in terminals that support ANSI escape codes(Unix/Linux/Windows platforms)
2. TrueColor Requirement: Hex and RGB colors require terminal with truecolor support
3. Style Support: Some terminals lack support for certain styles (blink, double underline).

# Platform Support

- Linux/macOS terminals (full support)
- Windows Terminal/WSL, CMD, Powershell (full support)
- Legacy Windows CMD (limited - some styles may not render)

# Contributing

We welcome contributions! Here's how you can help:

1. Report Bugs: Open an issue with reproduction steps
2. Suggest Features: Share your ideas for improvements
3. Submit PRs:
   - Fork the repository
   - Create a feature branch
   - Add tests for your changes
   - Ensure code follows Go conventions
   - Submit a pull request


# Development Setup

```bash
# Clone the repository
git clone https://github.com/phamio/crayon.git
cd crayon

# Run tests
go test 

```

# Areas Needing Improvement

1. Better Windows compatibility
2. Performance optimization


# License

MIT License - see [LICENSE](LICENSE) file for details.

# Acknowledgments

- ANSI escape code  specifications
- Thanks to [kurahaupo](https://gist.github.com/kurahaupo/6ce0eaefe5e730841f03cb82b061daa2) for explanations on true color.
- The Go community for testing and feedback
- All contributors who have helped improve this library

---

**Note**: Always test color output in different terminals to ensure compatibility with your users' environments. Consider providing a --no-color flag in your applications for users who prefer plain text.


