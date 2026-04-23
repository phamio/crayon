
# Crayon

A comprehensive Go library for adding colors, styles, and formatting to terminal output with support for multiple color formats and truecolor detection. [This is a Go port of the Spectra color library](https://github.com/phamio/spectra)


# Installation

```bash
go get -u github.com/phamio/crayon
```

# Features

- Multiple Color Systems: Named colors, hex codes, RGB, 256-color palette.
- Full ColorFallback Chain: Automatic downsampling across all color levels — truecolor -> 256 color palette -> ANSI 16 colors.
- TrueColor Detection: Automatic detection of terminal truecolor support.
- Simples and Readable: Easy-to-use template system for text styling and coloring.
- Comprehensive Styles: Bold, italic, underline, blink, reverse, hidden, strike-through.
- Granular Resets: Individual and full reset codes for precise control.
- Escape System: Texts in [] that aren't colors/styles are left as it is.
- Inline Padding: Left and right alignment, declared directly on placeholders.
- Cross-Platform: Full color support on Windows (Windows Terminal, cmd, PowerShell), Linux and macOS.
- Color Toggling: Respects NO_COLOR environment variable and detects when output is redirected.

---


# Quick Start

```go
package main

import (
    "fmt"
    "github.com/phamio/crayon"
)

func main() {
    // basic ansi16
    crayon.Parse("[fg=blue bg=yellow]Hello in blue on yellow![reset]").Println()

    // 256-color palette
    crayon.Parse("[fg=4 bg=3]Hello in blue on yellow![reset]").Println()
    crayon.Parse("[fg=214]Orange from 256-color palette[reset]").Println()
    crayon.Parse("[bg=196]Red background from palette[reset]").Println()

    // RGB colors
    crayon.Parse("[fg=rgb(0,0,128) bg=rgb(255,255,0)]Hello in blue on yellow![reset]").Println()
    crayon.Parse("[fg=rgb(255,105,180)]Hot pink text[reset]").Println()
    crayon.Parse("[bg=rgb(50,205,50)]Lime green background[reset]").Println()

    // Hex colors
    crayon.Parse("[fg=#000080 bg=ffff00]Hello in blue on yellow![reset]").Println()
    crayon.Parse("[fg=#FF5733]Orange hex color[reset]").Println()
    crayon.Parse("[bg=#3498db]Blue background[reset]").Println()


    // Styles
    crayon.Parse("[bold underline=single] Bold and underlined[reset]").Println()
    crayon.Parse("[italic dim strikethrough]Dim italic strikethrough  text only[reset]").Println()
    
    // Reset specific attributes
    crayon.Parse("[blink=slow]Slow blinking[hidden] hidden text [hidden=reset] Shown text[reset]").Println()
    crayon.Parse("[bold fg=blue]Blue bold text. [bold=reset]No longer bold, but still blue. [fg=reset]No color, but other styles remain[reset]").Println()
}
```
# Output

![quick_start](images/quick_start_merge.png)

---

# Template System
The library follows a template-first approach: parse color templates once with or without placeholders([0], [1], etc), then reuse them with different data to replace placeholders.

Templates are pre-parsed color/style definitions with slots. Pre-parsing templates is the heart of the library
Placeholders are like empty slots in templates that awaits dynamic data to be plugged in later- You define the template once, then insert varying content at runtime without re-parsing the styling.

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



    logTemplate := crayon.Parse("[0] [fg=blue][1][reset]: [fg=yellow][2][reset]")
    template := crayon.Parse("[[fg=red bold]Error]: [0][reset]")
    

    // Different log levels
    logTemplate.Println("[INFO]", "main", "Application started")
    logTemplate.Println("[WARN]", "auth", "Token expiring soon")
    logTemplate.Println("[ERROR]", "db", "Connection failed")

    

    // Reuse multiple times
    template.Println("File not found")
    template.Println("Permission denied")
    template.Println("Network timeout")
    

    // Progress bar template
    progressTemplate := crayon.Parse("[fg=cyan][0][reset]/[fg=cyan][1][reset] [fg=green][2][reset]%")
    
    total := 100
    for i := 0; i <= total; i += 10 {
        percent := i * 100 / total
        fmt.Printf("\r%s", progressTemplate.Sprint(i, total, percent))
        time.Sleep(100 * time.Millisecond)
    }
}
```

---



# Padding
Padding is used for text formatting. Crayon padding measures the visible text length, ignoring ANSI color codes, so colored text aligns correctly in columns.
Padding is applied directly on placeholders

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

## Color Toggling
Crayon automatically detects whether color should be enabled based on NO_COLOR environment variable and whether output is going to a real terminal(TTY). You can also control this manually.

```go
package main

import (
    "fmt"
    "os"
    "github.com/phamio/crayon"
)

func main() {
    // Crayon decides - colors on for TTY, off for piped output
    toggle := crayon.NewColorToggle()
    
    successTemplate := toggle.Parse("[fg=green]✓ [0][reset]")
    //this  successTemplate := toggle.Parse("[fg=green]✓ [0][reset]") is same as successTemplate := crayon.Parse("[fg=green]✓ [0][reset]") the latter implictly uses crayon default color toggle

    errorTemplate := toggle.Parse("[fg=red]✗ [0][reset]")
    
    successTemplate.Println("Operation completed")
    errorTemplate.Println("Operation failed")
    

    // Manual control
    forceOn := crayon.NewColorToggle(true)   // Always show colors
    forceOff := crayon.NewColorToggle(false)     // Never show colors
    
    
    // Use in CLI applications
    //Respect both --no-color flag and NO_COLOR environment variable
    noColorFlag = false
    for _, arg := range os.Args {
        if arg == "--no-color" {
            noColorFlag = true
            break
        }
    }
    useColor := !noColorFlag && os.Getenv("NO_COLOR") == ""
    appToggle := crayon.NewColorToggle(useColor)
    
    helpTemplate := appToggle.Parse("[bold fg=cyan][0][reset] [fg=green][1][reset]")
    helpTemplate.Println("Usage:", "myapp [options]")
}
```

---




# Patterns to avoid
Crayon performance comes from parsing templates once and reusing them. Parsing inside a loop throws that advantage away

``` go
package main
import "github.com/phamio/crayon"

// BAD: parses every iteration, slow
func processItems(items []string) {
    for _, item := range items {
        // DON'T DO THIS - parses every iteration!
        tmpl := crayon.Parse("[fg=blue]" + item + "[reset]")
        tmpl.Println()
    }
}

// GOOD: parse once, reuse in loop
var itemTemplate = crayon.Parse("[fg=blue][0][reset]")
func processItems(items []string) {
    for _, item := range items {
        tmpl.Println(item)
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
    template := crayon.Parse("[bold fg=red]Item [0][reset] [fg=green]Value[1][reset]")
    
    start := time.Now()
    for i := 0; i < iterations; i++ {
        template.Sprint(i, i)
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
| `[0:>10]` | placeholder=0, direction=right align, width=10 |




# Tips and Best Practices
1. Parse Once: Always parse templates at initialization.
3. Template Reuse: Create templates for consistent styling
4. Placeholder Limits: The current implementation supports [0] through [999]


# Limitations
1. Terminal Dependency: Colors only work in terminals that support ANSI escape codes. Legacy Windows cmd has limited style support.
3. Style Support: Some styles like blink and double underline are not universally supported across all terminals — behaviour may vary.

---


# Platform Support
- Linux — full support
- macOS — full support
- Windows Terminal — full support
- Windows cmd — full support, limited styles
- Powershell — full support
- Legacy Windows CMD — limited — some styles may not render

---

# Contributing
Contributions are welcome! Here's how you can help:

## Report Bugs
Open an issue with clear description and reproduction steps.

## Suggest Features
Share your ideas and the problem they solve

## Submit Pull-Requests
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

1. Windows testing — verifying behaviours across CMD, Powershell and Windows terminal.
2. Performance optimization.


# License

MIT License - see [LICENSE](LICENSE) file for details.

# Acknowledgments
- [mitchellh/colorstring](https://github.com/mitchellh/colorstring) — early inspiration for the bracket syntax approach
- ANSI escape code specifications — the foundation everything is built on.
- Spectra(https://github.com/phamio/crayon) — the Nim library Crayon was ported from.
- Thanks to [kurahaupo](https://gist.github.com/kurahaupo/6ce0eaefe5e730841f03cb82b061daa2) for explanations on true color.
- The Go community for testing and feedback
- All contributors who have helped improve this library



