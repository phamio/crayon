package crayon

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// FIXES THAT ARE IN SESSION (NEEDED BEFORE v1.0.0)
// - 256 fallback if terminal doesnt have true color support [DONE]

//===========================================
//  COLOR VALIDATION
//===========================================

func isValidHex(hexCode string) bool {
	//fg=#RRGGBB
	if len(hexCode) == 10 && (strings.HasPrefix(hexCode, "fg=#") || strings.HasPrefix(hexCode, "bg=#")) {
		matched, _ := regexp.MatchString(`^[0-9a-fA-F]+$`, hexCode[4:])
		if len(hexCode[4:]) == 6 || matched {
			return true
		}
	}
	return false
}

func isValid256Code(paletteCode string) bool {
	//low = fg=0
	//high = fg=255
	if len(paletteCode) >= 4 && len(paletteCode) <= 6 && (strings.HasPrefix(paletteCode, "fg=") || strings.HasPrefix(paletteCode, "bg=")) {
		parsedInt, err := strconv.Atoi(paletteCode[3:])
		if err != nil {
			return false
		}
		return parsedInt >= 0 && parsedInt <= 255
	}
	return false
}

func isValidRGB(rgbCode string) bool {
	//low = fg=rgb(r,g,b)
	//high = fg=rgb(rrr,ggg,bbb)

	//includes positions 3,4,5,6 excludes position 7
	if len(rgbCode) >= 13 && len(rgbCode) <= 19 && (strings.HasPrefix(rgbCode, "fg=") || strings.HasPrefix(rgbCode, "bg=")) {
		if !strings.HasPrefix(rgbCode[3:], "rgb(") && !strings.HasSuffix(rgbCode, ")") {
			return false
		}
		//extract content to see if each value is in 0..255 and are numbers
		seqNumbers, boolean := readRGB(rgbCode)
		//true means successfully extracted and are numbers
		if boolean {
			for _, num := range seqNumbers {
				if num >= 0 && num <= 255 {
					//checks first match and returns

					return true
				}
				return false
			}
		}
	}
	return false
}

func supportsTrueColor() bool {
	colorterm := os.Getenv("COLORTERM")
	return colorterm == "truecolor" || colorterm == "24bit"
}

// this function was made to validate words in []
func isSupportedColor(input string) bool {
	_, inColorMap := ColorMap[input]
	_, inResetMap := ResetMap[input]
	_, inStyleMap := StyleMap[input]

	return inColorMap || inResetMap || inStyleMap || isValidHex(input) || isValid256Code(input) || isValidRGB(input)
}

func readRGB(rgbCode string) ([]int, bool) {
	//fg=rgb(rrr,ggg,bbb)
	var result []int
	end := len(rgbCode) - 1
	numbers := strings.Split(rgbCode[7:end], ",")
	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Printf("Error parsing %s: %v", numStr, err)
			return nil, false
		}
		result = append(result, num)
	}
	return result, true
}

// ======================================
// COLOR PARSING
// ======================================
func parseAnsi(colorCode string, ansiAppend string) string {
	if strings.HasPrefix(colorCode, "bg=") {
		return fmt.Sprintf("\033[48;%sm", ansiAppend)
	} else if strings.HasPrefix(colorCode, "fg=") {
		return fmt.Sprintf("\033[38;%sm", ansiAppend)
	}
	return ""
}

func parseRGBToAnsiCode(rgbCode string) string {
	RGB, _ := readRGB(rgbCode)
	if supportsTrueColor() {
		return parseAnsi(rgbCode, fmt.Sprintf("2;%d;%d;%d", RGB[0], RGB[1], RGB[2]))
	}
	//Use 256 color palette fallback
	//its obvious the fallback is broken, if given rgb is not in
	//delta match, it won't catch.
	//im just trying to learn and understand more aboyt terminal colouring.
	fallBackCode := hexTo256Fallback(rgbToHex(int64(RGB[0]), int64(RGB[1]), int64(RGB[2])))
	return parseAnsi(rgbCode, fmt.Sprintf("5;%s", fallBackCode))
}

func parseHexToAnsiCode(hexCode string) string {
	//fg=#RRGGBB
	if len(hexCode) == 10 {
		if supportsTrueColor() {
			R, _ := strconv.ParseInt(hexCode[4:6], 16, 32)
			G, _ := strconv.ParseInt(hexCode[6:8], 16, 32)
			B, _ := strconv.ParseInt(hexCode[8:10], 16, 32)

			return parseAnsi(hexCode, fmt.Sprintf("2;%d;%d;%d", R, G, B))
		}
		//Use 256 color palette fallback
		fallBackCode := hexTo256Fallback(hexCode[4:])
		return parseAnsi(hexCode, fmt.Sprintf("5;%s", fallBackCode))
	}
	return ""
}

/* Note:
    #foreground colors use 38 and background colors use 48. the 2 is for truecolor support
so its \e[38;2;R;G;Bm or for background \e[48;2;R;G;Bm
so the second row of number tells what color mode it is (2: rgb(24 bits), 245)
 2 is for truecolor supported numbers that is rgb and its 24 bits using a range of 0-255
 5 is for 256 palette(index 196)
 256 palette support syntax will be [fg=214] = foreground color and [bg=214] = background color*/

func parse256ColorCode(colorCode string) string {
	fmt.Println("PARSING ANSI")
	return parseAnsi(colorCode, fmt.Sprintf("5;%s", colorCode[3:]))
}

func ParseColor(color string) string {
	//this function is meant to receive string like "bold" "fg=red" and other colors and
	//convert them to their ansi codes
	if code, exists := ColorMap[color]; exists {
		return fmt.Sprintf("\033[%sm", code)
	}

	if code, exists := StyleMap[color]; exists {
		return fmt.Sprintf("\033[%sm", code)
	}

	if code, exists := ResetMap[color]; exists {
		return fmt.Sprintf("\033[%sm", code)
	}

	if isValid256Code(color) {
		return parse256ColorCode(color)
	}

	if isValidHex(color) {
		return parseHexToAnsiCode(color)
	}

	if isValidRGB(color) {
		return parseRGBToAnsiCode(color)
	}
	return ""
}
