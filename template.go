package crayon

import (
	"fmt"
	"io"
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

	//Handle unclosed sequence at the end of input
	if inReadSequence && len(contentSequence) > 0 {
		//Treat unclosed bracket as literal
		parts = flushText(parts, currentText)
		parts = append(parts, TempPart{Text: "[" + contentSequence, Index: -1, FormatStr: ""})
		currentText = ""
	}
	return parts, currentText
}


//=============================
// PARSE - BRACKET HANDLERS
//=============================

func handleOpenBracket(i int, input string, parts []TempPart, currentText string) ([]TempPart, string, string, bool) {
    //Bracket peeling logic - Peel brackets to see if content is a color or not
	//consider first '[' as a text, move until, content is found.
	if i+1 < len(input) && input[i+1] == '[' {

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
			parts = append(parts, TempPart{Text: parseColor(w), Index: -1, FormatStr: ""})
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

	//unrecognized -  pass through as literal
	return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
}

//=============================
// PARSE - PLACEHOLDER
//=============================
//extract placeholders
//placeholders will support padding too. [0:<20] = left alignment, [0:>20] = right align

//========= WHY PADDINGS WERE ADOPTED =========
//crayon added inline padding because doing so with fmt.Printf was cumbersome considering repeated output
//pad := crayon.Parse("[fg=red]Error: [0][reset]")
//using fmt.Printf("%-20s", pad.Sprint("File Not Found"))
//This left aligns the whole "\033[31mError: File Not Found\033[0m" instead of only "File Not Found"

//Although there's a fix, its verbose
// use pad.Println(fmt.Sprintf("%-20s", "File Not Found"))

//So crayon opted for {Define once, pad many times}
// pad := crayon.Parse("[fg=red]Error: [0:<20][reset]")
//  pad.Println("File Not Found") which correctly left aligns only "File Not Found"

//crayon's padding is nothing special, it just does this
//padIt := fmt.Sprintf("%-20s", "File Not Found") in the backend
// pad.Println(padIt)
//Saving you less typing strokes and efficient for repeated outputs
//========= END =========


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
	align, width, boolean:= parseAlignWidth(padStr)
	if !boolean {
		return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
	}
	return append(parts, TempPart{Text: "", Index: index, FormatStr: buildFormatStr(align, width)})
}


// =============================
// PARSE - HELPERS
// =============================
func parseAlignWidth(input string) (rune, int, bool) {
	//'input' is aleady stripped of its placeholder
	if len(input) < 2 {
		return 0, 0, false
	}
	align := rune(input[0])
	if align != '<' && align != '>' {
		return 0, 0, false
	}

	widthStr := input[1:]
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return 0, 0, false
	}

	if width <= 0 {
		return 0, 0, false
	}

	return align, width, true
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
