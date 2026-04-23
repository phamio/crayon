package crayon

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)


//===========================================
//  COLOR DOWNSAMPLING/DEGRADATION
//===========================================
func rgbTo256Index(r, g, b int) int {
    //Find the maximum and minimum channel values
	//to compute the range, spread across RGB channels.
	//A small range means the color is close to neutral/grey.
	maxC := r 
	if g > maxC { maxC = g }
	if b > maxC { maxC = b }

	minC := r 
	if g < minC { minC = g }
	if b < minC { minC = b }
    

	//Used to determine how dark the color is
	avg := (r + g + b ) / 3
    
	//====== GRAYSCALE RAMP ROUTING =========
	//Route to the 24-step grayscale ramp if 
	// - maxC-minC <= 20: The cube is too coarse for neutral tones
	//and introduces visible color casts such as pinkish, greenish, etc.
	//So the threshold was made wide enough (20)
	// - avg > 5:  it allows dark grays to correctly hit the ramp, whiles true or near blacks passes over (0, 0, 8)
	if maxC - minC <= 20 && avg > 5 {
		//Clamp avg to the valid grayscale ramp
		//Starts at RGB(8,8,8) and ends at RGB(238,238,238).

	if avg < 8 {
  		avg = 8
  	}
  	if avg > 238 {
  		avg = 238
  	}
	//Tries to map avg to grayscale ramp index 232-255
  	return  232 + (avg-8)/10
	//return  232 + ((avg-8)*23/247)
  }	

  //====== COLOR CUBE ROUTING ======
  // for colors where the channel spread exceeds 20.
    r6 := (r * 5 + 127) / 255
	g6 := (g * 5 + 127) / 255
	b6 := (b * 5 + 127)/ 255
    //fmt.Printf("RGB TO INDEX FROM COLOR HELPERS CODE (NOT TEST):  RGB=(%d,%d,%d)  | 256 = %d\n", r, g, b, 16 + 36*r6 + 6*g6 +b6)
	return 16 + 36*r6 + 6*g6 +b6

}


//==================================
// ANSI TERMINAL SUPPORT DETECTION
//==================================
//This only check what the terminal advertises, which isn't always accurate.
//It has to query the terminal directly
//Ways to query the terminal
//echo -e "\e[c"  //What ansi capabilites are supported
//echo -e "\e[>c" // Which terminal emulator you're actually in

func supportsTrueColor() bool {
	colorterm := os.Getenv("COLORTERM")
	return colorterm == "truecolor" || colorterm == "24bit"
}

func supports256Color() bool {
	term := os.Getenv("TERM")
	return strings.Contains(term, "256color")
}

func supportsNone() bool {
	noneTerm := os.Getenv("TERM")
	//return noneTerm == "dumb"
	return noneTerm == "dumb" || strings.Contains(noneTerm, "mono")
}

//===========================================
//  COLOR VALIDATION
//===========================================
func hasValidPrefix(inputCode string) bool {
	return (strings.HasPrefix(inputCode, "fg=") || strings.HasPrefix(inputCode, "bg="))
}

func hasSpecialPrefix(inputCode string, prefix string) bool {
	return (strings.HasPrefix(inputCode, fmt.Sprintf("fg=%s", prefix)) || strings.HasPrefix(inputCode, fmt.Sprintf("bg=%s", prefix)))
}

func isHexCode(hexCode string) bool{
	for _, ch := range hexCode {
		if !(byte(ch) >= '0' && byte(ch) <= '9' || byte(ch) >= 'a' && byte(ch) <= 'f' || byte(ch) >= 'A' && byte(ch) <= 'F'){
			return false
		}
	}
	return true
}

//this one also reads the value and throws it away
func isValidHex(hexCode string) bool {
	if len(hexCode) == 10 && hasSpecialPrefix(hexCode, "#") && isHexCode(hexCode[4:]) {
			return true
	}
	return false
}


func isValid256Code(paletteCode string) (int, bool) {
	if len(paletteCode) >= 4 && len(paletteCode) <= 6 && hasValidPrefix(paletteCode) {
		parsedInt, err := strconv.Atoi(paletteCode[3:])
		if err != nil {
			return 0, false
		}
		return parsedInt, parsedInt >= 0 && parsedInt <= 255
	}
	return 0, false
}

func isValidRGB(rgbCode string) ([]int, bool) {
	//includes positions 3,4,5,6 excludes position 7
	if len(rgbCode) >= 13 && len(rgbCode) <= 19 && hasValidPrefix(rgbCode) && strings.HasPrefix(rgbCode[3:7], "rgb(") && strings.HasSuffix(rgbCode, ")") {
		//extract content to see if each value is in 0..255 and are numbers
		seqNumbers, boolean := parseRGB(rgbCode)
		//true means successfully extracted and are numbers
		if boolean  && seqNumbers != nil{
			if len(seqNumbers) != 3 {
				return nil, false
			}

			for _, num := range seqNumbers {
				
				if num < 0 || num > 255 {
					return nil, false
				}
			}
			return seqNumbers, true
		}
		return nil, false
	}
	return nil, false
}



// this function was made to validate words in []
func isSupportedColor(input string) bool {
	_, inColorMap := colorMap[input]
	_, inResetMap := resetMap[input]
	_, inStyleMap := styleMap[input]
	_, validRGB := isValidRGB(input)
	_, valid256 := isValid256Code(input)

	return inColorMap || inResetMap || inStyleMap || isValidHex(input) || valid256 || validRGB
}

func parseRGB(rgbCode string) ([]int, bool) {
	//fg=rgb(rrr,ggg,bbb)
	var num int
	var err error
	var result []int
	end := len(rgbCode) - 1
	numbers := strings.Split(rgbCode[7:end], ",")

	if len(numbers) == 3 {
	for _, numStr := range numbers {
		num, err = strconv.Atoi(numStr)
		  if err != nil {
			return nil, false
		  }
		  result = append(result, num)
	  }
	  return result, true
    }
	return nil, false

}

// ======================================
// COLOR PARSING
// ======================================
func parseAnsi(colorCode string, ansiAppend string, isAnsi16 bool) string {

	if strings.HasPrefix(colorCode, "bg=") {
		if isAnsi16{
			ansiInt, _ := strconv.Atoi(ansiAppend)
		    ansiInt = ansiInt + 10
	    	return fmt.Sprintf("\033[%dm", ansiInt)
		}
		return fmt.Sprintf("\033[48;%sm", ansiAppend)

	} else if strings.HasPrefix(colorCode, "fg=") {
		if isAnsi16{
			return fmt.Sprintf("\033[%sm", ansiAppend)
		}
		return fmt.Sprintf("\033[38;%sm", ansiAppend)
	}
    
	return ""
}

func parseRGBToAnsiCode(rgbCode string, RGB []int) string {
	if supportsTrueColor() {
		return parseAnsi(rgbCode, fmt.Sprintf("2;%d;%d;%d", RGB[0], RGB[1], RGB[2]), false)
	}
	//256 palette fallback
	if supports256Color(){
		return parseAnsi(rgbCode, fmt.Sprintf("5;%d", rgbTo256Index(RGB[0], RGB[1], RGB[2])), false)
	}
	if !supportsNone() {
	//ansi 16 fallback
	return parseAnsi(rgbCode, fmt.Sprint(ansi256ToAnsi16Lut[rgbTo256Index(RGB[0], RGB[1], RGB[2])]), true)
	}
	return ""
}

func parseHexToAnsiCode(hexCode string) string {
	//fg=#RRGGBB
		R, _ := strconv.ParseInt(hexCode[4:6], 16, 32)
		G, _ := strconv.ParseInt(hexCode[6:8], 16, 32)
	    B, _ := strconv.ParseInt(hexCode[8:10], 16, 32)
		return parseRGBToAnsiCode(hexCode, []int{int(R), int(G), int(B)})
}

/* Note:
    #foreground colors use 38 and background colors use 48. the 2 is for truecolor support
so its \e[38;2;R;G;Bm or for background \e[48;2;R;G;Bm
so the second row of number tells what color mode it is (2: rgb(24 bits), 245)
 2 is for truecolor supported numbers that is rgb and its 24 bits using a range of 0-255
 5 is for 256 palette(index 196)
 256 palette support syntax will be [fg=214] = foreground color and [bg=214] = background color*/

func parse256ColorCode(colorCode string, paletteCode int) string {
	if supports256Color(){
		return parseAnsi(colorCode, fmt.Sprintf("5;%d", paletteCode), false)
	}
	if !supportsNone() {
	return parseAnsi(colorCode, fmt.Sprint(ansi256ToAnsi16Lut[paletteCode]), true)
	}
	return ""
}


func parseColor(color string) string {
	//this function is meant to receive string like "bold" "fg=red" and other colors and
	//convert them to their ansi codes
	if code, exists := colorMap[color]; exists {
		return fmt.Sprintf("\033[%sm", code)
	}

	if code, exists := styleMap[color]; exists {
		return fmt.Sprintf("\033[%sm", code)
	}

	if code, exists := resetMap[color]; exists {
		return fmt.Sprintf("\033[%sm", code)
	}


	if palette, ok := isValid256Code(color); ok{
		return parse256ColorCode(color, palette)
	}

	if isValidHex(color) {
		return parseHexToAnsiCode(color)
	}
    
	if rgb, ok := isValidRGB(color); ok{
		return parseRGBToAnsiCode(color, rgb)
	}
	return ""
}
