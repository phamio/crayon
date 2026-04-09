// Package crayon provides terminal colors and styles for Go.
// A better go doc is needed to be written
package crayon

import (
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type TempPart struct {
	Text      string
	Index     int
	FormatStr string
}

type CompiledTemplate struct {
	Parts       []TempPart
	TotalLength int
}

type ColorToggle struct {
	EnableColor bool
}

//=============================
// COLOR TOGGLE
//=============================

func autoDetect() bool {
	if _, exists := os.LookupEnv("NO_COLOR"); exists {
		return false
	}
	return term.IsTerminal(int(os.Stdout.Fd()))
}

//[TEMPORARILY]
//for escapes, it will be [<content>] so that anyone can use eg. [12:30] (time literal) without getting errors.
//such escape will be used for color and styles too

//Escapes
// [[content]] //needs lookahead and a whole lot of complexities. the parser was made to fight such escape
//when it was new and needed no escape system.
// \[content\]  //just prefix and suffix look but may conflict with already existing golang string adages.

//=============================
// PARSE - LOOP
//=============================

func parseLoop(input string, enableColor bool) ([]TempPart, string) {
	var (
		parts           []TempPart
		currentText     string
		contentSequence string
		inReadSequence  bool
	)

	for i, ch := range input {
		char := string(ch)

		switch {
		case char == "[" && !inReadSequence:
			parts, currentText, contentSequence, inReadSequence = handleOpenBracket(i, input, parts, currentText)

		case ch == ']' && inReadSequence:
			parts, inReadSequence = handleCloseBracket(contentSequence, parts, enableColor)
			contentSequence = ""

		case inReadSequence:
			contentSequence += char

		default:
			currentText += char
		}
	}
	return parts, currentText
}

//[[content]] will be the accepted escape, the parser checks if its at least 2 chars "[[", same for
// at least 2 chars "]]"
//=============================
// PARSE - BRACKET HANDLERS
//=============================

func handleOpenBracket(i int, input string, parts []TempPart, currentText string) ([]TempPart, string, string, bool) {
	//check if the next value is "["
	// [[fg=color]] should never be an escape, because first ']' terminates early without lookahead [WILL BE FIXED]
	//besides escape is just meant to be an opt in

	//consider first '[' as a text, move until, content is found.
	if i+1 < len(input) && input[i+1] == '[' {
		//for escapes


		fmt.Println("NEXT INPUT: ", string(input[i+2]))
		fmt.Println("LAST INPUT: ", string(input[len(input)-6:]))
		fmt.Println("Text: ", currentText)
		currentText += "["
		return parts, currentText, "", false
	}
	//flush current text before entering sequence
	parts = flushText(parts, currentText)
	return parts, "", "", true
}

func handleCloseBracket(contentSequence string, parts []TempPart, enableColor bool) ([]TempPart, bool) {
	allWords := strings.Fields(contentSequence)

	if isColorSequence(allWords) {
		parts = handleColorSequence(parts, allWords, enableColor)
	} else {
		parts = handleNonColorSequence(parts, contentSequence)
	}
	return parts, false
}

//=============================
// PARSE - SEQUENCE HANDLERS
//=============================

func isColorSequence(words []string) bool {
	if len(words) == 0 {
		return false
	}
	for _, w := range words {
		if !isSupportedColor(w) {
			return false
		}
	}
	return true
}

func handleColorSequence(parts []TempPart, words []string, enableColor bool) []TempPart {
	if enableColor {
		for _, w := range words {
			parts = append(parts, TempPart{Text: ParseColor(w), Index: -1, FormatStr: ""})
		}
	} else {
		parts = append(parts, TempPart{Text: "", Index: -1, FormatStr: ""})
	}
	return parts
}

func handleNonColorSequence(parts []TempPart, contentSequence string) []TempPart {
	if isValidPlaceholder(contentSequence) {
		return handlePlaceholder(parts, contentSequence)
	}

	//for padded placeholders
	if strings.Contains(contentSequence, ":") && !strings.HasPrefix(contentSequence, "<") && !strings.HasSuffix(contentSequence, ">") {
		return handlePaddedPlaceholder(parts, contentSequence)
	}

	//for escapes
	if strings.HasPrefix(contentSequence, "<") && strings.HasSuffix(contentSequence, ">") {
		return handleEscape(parts, contentSequence)
	}
	//unrecognized -  pass through as literal
	return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
}

//=============================
// PARSE - PLACEHOLDER
//=============================
//extract placeholders
//placeholders will support padding too. [0:<20] = left alignment, [0:>20] = right align

//========= WHY PADDINGS WERE ADOPTED =========
//crayon added inline padding because doing so with fmt.Printf was close to imposible
//pad := crayon.Parse("[fg=red]Error: [0][reset]")
//fmt.Printf("%-20s", pad.Sprint("File Not Found"))
//This left aligns the whole "Error: File Not Found" instead of only "File Not Found"
//So crayon opted in for
// pad := crayon.Parse("[fg=red]Error: [0:<20][reset]")
//  pad.Println("File Not Found") which correctly left aligns only "File Not Found"
//========= END =========

//Overflow handling will slow down crayon. I'm still on the fence of throwing it away or using it
//It will slow down crayon because calculation will be moved to apply, thats not the work of apply
//[0:<20!] = right alignment  with truncation, [0:>20~] = left align with elipsis(...),
//[0:<20?] = right alignment  with warn to stderr

func isValidPlaceholder(input string) bool {
	return len(input) > 0 && allDigits(input)
}

func handlePlaceholder(parts []TempPart, contentSequence string) []TempPart {
	//digit boundary guard to prevent overflow
	index, err := strconv.Atoi(contentSequence)
	if err == nil && index >= 0 && index <= 999 {
		return append(parts, TempPart{Text: "", Index: index, FormatStr: ""})
	}
	//out of range - treat as literal
	return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
}

func handlePaddedPlaceholder(parts []TempPart, contentSequence string) []TempPart {
	//[0:>20] stripped of its brackets ==> 0:>20
	splitWord := strings.SplitN(contentSequence, ":", 2) // ==> [0 >20]
	if len(splitWord) != 2 {
		return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
	}
	indexStr := splitWord[0]
	padStr := splitWord[1]
	//parse indexStr
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index > 999 {
		return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
	}

	//parse the padStr
	align, width, err := parseAlignWidth(padStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		//for the sake of no better escape yet, accept as literal after error propagation
		//once escape is created, enforce its use by printing error and stopping.
		//Not everything is meant to be treated as literal, some are meant to be brought to notice.
		return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
	}
	return append(parts, TempPart{Text: "", Index: index, FormatStr: buildFormatStr(align, width)})
}

// =============================
// PARSE - ESCAPE
// =============================
// This is only a temporary escape [<content>] => <content> 
// This will be the main escape [[content]]
// strip it of its angle brackets
func handleEscape(parts []TempPart, contentSequence string) []TempPart {
	contentSequence = strings.TrimPrefix(contentSequence, "<")
	contentSequence = strings.TrimSuffix(contentSequence, ">")
	return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})

}

// =============================
// PARSE - HELPERS
// =============================
func parseAlignWidth(input string) (rune, int, error) {
	//'input' is aleady stripped of its placeholder
	if len(input) < 2 {
		return 0, 0, fmt.Errorf("padding spec: too short - expected '<N' or '>N' (got '%s')", input)
	}
	align := rune(input[0])
	if align != '<' && align != '>' {
		return 0, 0, fmt.Errorf("padding spec: expected '<' or '>', got '%c'", align)
	}

	widthStr := input[1:]
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return 0, 0, fmt.Errorf("padding spec: width must be a number (got '%s')", widthStr)
	}

	if width <= 0 {
		return 0, 0, fmt.Errorf("padding spec: must be greater than zero (got %d)", width)
	}

	return align, width, nil
}

func buildFormatStr(align rune, width int) string {
	switch align {
	//left align
	case '<':
		return fmt.Sprintf("%%-%ds", width)
		//right align
	case '>':
		return fmt.Sprintf("%%%ds", width)
	}
	return ""
}

func flushText(parts []TempPart, currentText string) []TempPart {
	if len(currentText) > 0 {
		parts = append(parts, TempPart{Text: currentText, Index: -1})
	}
	return parts
}

func allDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

//=============================
// APPLY
//=============================
//what if sum of apply(args) is greater than estimation 
func (temp CompiledTemplate) apply(args ...any) string {
	//Calculate estimated size for optimization
	var totalArgLength int
	for _, arg := range args {
		totalArgLength += len(fmt.Sprint(arg))
	}

	estimatedSize := temp.TotalLength + totalArgLength
	var result strings.Builder
	result.Grow(estimatedSize)

	for _, part := range temp.Parts {
		if part.Index < 0 {
			result.WriteString(part.Text)
		} else if part.Index < len(args) {
			value := fmt.Sprint(args[part.Index])
			if part.FormatStr != "" {				
				value = fmt.Sprintf(part.FormatStr, value)
			}
			result.WriteString(value)
		}
	}
	return result.String()
}

// =======================
// PRINT
// =======================
func (temp CompiledTemplate) Println(args ...any) {
	fmt.Println(temp.apply(args...))
}

func (temp CompiledTemplate) Print(args ...any) {
	fmt.Print(temp.apply(args...))
}

// =======================
// FPRINT
// =======================
func (temp CompiledTemplate) Fprintln(w io.Writer, args ...any) (n int, err error) {
	return fmt.Fprintln(w, temp.apply(args...))
}

func (temp CompiledTemplate) Fprint(w io.Writer, args ...any) (n int, err error) {
	return fmt.Fprint(w, temp.apply(args...))
}

// =======================
// SPRINT
// =======================
func (temp CompiledTemplate) Sprintln(args ...any) string {
	return fmt.Sprintln(temp.apply(args...))
}

func (temp CompiledTemplate) Sprint(args ...any) string {
	return temp.apply(args...)
}
