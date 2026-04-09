package crayon

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//===== FIXES BEFORE v1.0.0 ========
// Cross platform support [DONE]
// RGB-TO-256 fallback [DONE]
// Dumb terminals [YET TO KNOW ITS WORKINGS]
// Escape system [NOT YET]
// Fast Parsing [IN PROGRESS]
// 256 to ansi colors for terminals that dont support 256 [IN PROGRESS]
// ========== END =============


//===========================================
//  RGB TO 256 PALETTE FALLBACk
//===========================================

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func rgbTo256Index(r, g, b int) int {
    r6 := (r * 5 + 127) / 255
	g6 := (g * 5 + 127) / 255
	b6 := (b * 5 + 127)/ 255
	//cubeIndex := 16 + 36*r6 + 6*g6 +b6


  //check if it's close enough to gray
  if abs(r-g) < 10 && abs(g-b) < 10 {
  	avg := (r + g + b) / 3
  	if avg < 8 {
  		//return 16 //closest is black in the color cube
  		avg = 8
  	}
  	if avg > 238 {
  		//return 231 //closest is white in the color cube
  		avg = 238
  	}
  	return  232 + (avg-8)/10
  }
	
	return 16 + 36*r6 + 6*g6 +b6
}



//===========================================
//  COLOR VALIDATION
//===========================================
func isHex(hexCode string) bool{
	for _, ch := range hexCode {
		if !isHexDigit(byte(ch)){
			return false
		}
	}
	return true
}

func isHexDigit(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func isValidHex(hexCode string) bool {
	if len(hexCode) == 10 && (strings.HasPrefix(hexCode, "fg=#") || strings.HasPrefix(hexCode, "bg=#")) {
		if len(hexCode[4:]) == 6 || isHex(hexCode[4:]) {
			return true
		}
	}
	return false
}


func isValid256Code(paletteCode string) bool {
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
				
				if num < 0 || num > 255 {
					return false
				}
			}
		}
		return true
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
	//256 palette fallback
		return parseAnsi(rgbCode, fmt.Sprintf("5;%d", rgbTo256Index(RGB[0], RGB[1], RGB[2])))
}

func parseHexToAnsiCode(hexCode string) string {
	//fg=#RRGGBB
	if len(hexCode) == 10 {
		R, _ := strconv.ParseInt(hexCode[4:6], 16, 32)
		G, _ := strconv.ParseInt(hexCode[6:8], 16, 32)
	    B, _ := strconv.ParseInt(hexCode[8:10], 16, 32)

		if supportsTrueColor() {
			return parseAnsi(hexCode, fmt.Sprintf("2;%d;%d;%d", R, G, B))
		}
		//256 palette fallback
		return parseAnsi(hexCode, fmt.Sprintf("5;%d", rgbTo256Index(int(R), int(G), int(B))))
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
	return parseAnsi(colorCode, fmt.Sprintf("5;%s", colorCode[3:]))
}


// will be made a private function in v0.7.0
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

	if isValidRGB(color) /*reads and throws values away*/ {
		//got no way to reuse values that isValidRGB read because prefix or color is needed too, hence re-reading it again
		return parseRGBToAnsiCode(color)
	}
	return ""
}
