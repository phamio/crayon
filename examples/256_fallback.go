package main

import (
	"fmt"
	"math"
)

func main(){

	testCases := [][]int {
		{100, 100, 101}, //Near Gray (1 diff)
		{100, 100, 105}, //Near gray (5 Diff)
		{100, 100, 110}, // Near gray(10 diff - boundary)
		{128, 128, 128}, //Perfect gray
		{100, 100, 115}, //Cool gray
		{180, 170, 160}, //Warm gray
		{140, 135, 145}, //Purple gray
		{75, 80, 70}, //Olive-gray
		{0, 0, 1}, //Almost black
		{254, 254, 255}, //Almost white
		{0, 255, 0}, //Pure green
		{0, 128, 0}, //Mid green
		{0, 128, 128}, //Teal
		{139, 69, 19}, //Brown
		{150, 140, 150},
		{200, 210, 200},
		{0, 0, 8},
		{5, 0, 0},
		{8, 8, 0},
		{248, 255, 248},
		{250, 242, 250},
		{120, 128, 120},
		{160, 152, 160},
		{180, 180, 169},
		{245, 255, 245},
		{8, 8, 8},
		{12, 12, 12},
		{100, 100,115},
		{6, 6, 6},
		{255, 0, 0},
		{0, 128, 255},
		{200, 50, 100},
		{100, 200, 50},
		{100, 90, 80},
		{200, 185, 180},
		{50, 60, 70},
		{130, 120, 110},
		{3, 0, 5},
		{5, 5, 4},
		{0, 4, 0},
		{10, 10, 10},
		{50, 50, 50},
		{200, 200, 200},
		{0, 0, 0},
		{255, 255, 255},
		{255, 0, 0},
		{0, 255, 0},
		{0, 0, 255},
		{128, 128, 128},
		{255, 255, 0},
		{255, 0, 255},
		{0, 255, 255},

	}

	fmt.Println("Testing RGB to 256-color conversion")

	for _, rgb := range testCases {
		r, g, b := rgb[0], rgb[1], rgb[2]

		//Get indexes from each library
		gchalkIdx := callGChalk(uint8(r), uint8(g), uint8(b))
		oldCrayonIdx := callOldCrayon(r, g, b)
		secOldCrayonIdx := callSecOldCrayon(r, g, b)
		thirdOldCrayonIdx := callThirdOldCrayon(r, g, b)
		newCrayonIdx := callNewCrayon(r, g, b)
		//termenvIdx := callTermenv(r, g, b)

		fmt.Printf("RGB{%3d,%3d,%3d}\n", r, g, b)
		fmt.Printf("       TRUECOLOR:     \033[48;2;%d;%d;%dm      \033[0m\n", r, g, b)
		fmt.Printf("          gchalk: %3d %s\n", gchalkIdx, colorBlock(int(gchalkIdx)))
		fmt.Printf("      old crayon: %3d %s\n", oldCrayonIdx, colorBlock(oldCrayonIdx))
		fmt.Printf("  2nd old crayon: %3d %s\n", secOldCrayonIdx, colorBlock(secOldCrayonIdx))
		fmt.Printf("  3rd old crayon: %3d %s\n", thirdOldCrayonIdx, colorBlock(thirdOldCrayonIdx))
		fmt.Printf("      new crayon: %3d %s\n", newCrayonIdx, colorBlock(newCrayonIdx))
		//fmt.Printf("  termenv: %3d %s\n", termenvIdx, colorBlock(termenvIdx))
	}
}


func colorBlock(idx int) string {
	return fmt.Sprintf("\033[48;5;%dm      \033[0m", idx)
}

func callGChalk(red uint8, green uint8, blue uint8) uint8 {
	// Originally from // From https://github.com/Qix-/color-convert/blob/3f0e0d4e92e235796ccb17f6e85c72094a651f49/conversions.js

	// We use the extended greyscale palette here, with the exception of
	// black and white. normal palette only has 4 greyscale shades.
	if red == green && green == blue {
		if red < 8 {
			return 16
		}

		if red > 248 {
			return 231
		}

		return uint8(math.Round(((float64(red)-8)/247)*24)) + 232
	}

	return 16 +
		uint8(
			(36*math.Round(float64(red)/255*5))+
				(6*math.Round(float64(green)/255*5))+
				math.Round(float64(blue)/255*5))
}


//=========================
// CRAYON SECTION
//========================

//Old crayon doesn;t detect all greys. Grey catching is the priority as the cube doesnt handle them well.
//RGB(75,80,70): old crayon gives greenish tint, which isn't visually accurate. Proper grey detection would have saved the day
//RGB(180,170,160): old crayon gives pinkish which is questionable too. 

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// RGB to 256 palette fallback
func callOldCrayon(r, g, b int) int {
    r6 := (r * 5 + 127) / 255
	g6 := (g * 5 + 127) / 255
	b6 := (b * 5 + 127)/ 255
	//cubeIndex := 16 + 36*r6 + 6*g6 +b6


  //check if it's close enough to gray
  if abs(r-g) < 10 && abs(g-b) < 10 {
  	avg := (r + g + b) / 3
  	if avg < 8 {
  		avg = 8
  	}
  	if avg > 238 {
  		avg = 238
  	}
  	return  232 + (avg-8)/10
  }	
    //fmt.Printf("RGB TO INDEX FROM COLOR HELPERS CODE (NOT TEST):  RGB=(%d,%d,%d)  | 256 = %d\n", r, g, b, 16 + 36*r6 + 6*g6 +b6)
	return 16 + 36*r6 + 6*g6 +b6

}

// RGB to 256 palette fallback
func callSecOldCrayon(r, g, b int) int {
    //Find the maximum and minimum channel values
	//to compute the range, spread across RGB channels.
	//A small range means the color is close to neutral/grey.
	maxC := r 
	if g > maxC { maxC = g }
	if b > maxC { maxC = b }

	minC := r 
	if g < minC { minC = g }
	if b < minC { minC = b }
    
	//Average luminance of the color
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



func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RGB to 256 palette fallback
func callThirdOldCrayon(r, g, b int) int {
    //r6 := (r * 5 + 127) / 255
	//g6 := (g * 5 + 127) / 255
	//b6 := (b * 5 + 127)/ 255
	//cubeIndex := 16 + 36*r6 + 6*g6 +b6
	maxC := max(r, max(g, b))
	//maxC := r 
	//if g > maxC { maxC = g }
	//if b > maxC { maxC = b }
	minC := min(r, min(g, b))

	//minC := r 
	//if g < minC { minC = g }
	//if b < minC { minC = b }

	avg := (r + g + b ) / 3
	spread := maxC - minC

	if spread < 15 {
  //check if it's close enough to gray
  //if abs(r-g) < 10 && abs(g-b) < 10 {
  	//avg := (r + g + b) / 3
    //if avg > 245 {
	//	return 255
	//}
	grayVal := avg

	if grayVal < 8 {
  		grayVal = 8
  	}
  	if grayVal > 247 {
  		grayVal = 247
  	}
  	return  232 + ((grayVal-8)*23/247)
  }	
   //cube mapping for colors
    r6 := (r * 5 + 127) / 255
	g6 := (g * 5 + 127) / 255
	b6 := (b * 5 + 127)/ 255
    //fmt.Printf("RGB TO INDEX FROM COLOR HELPERS CODE (NOT TEST):  RGB=(%d,%d,%d)  | 256 = %d\n", r, g, b, 16 + 36*r6 + 6*g6 +b6)
	return 16 + 36*r6 + 6*g6 +b6

}




// RGB to 256 palette fallback
func callNewCrayon(r, g, b int) int {
    //r6 := (r * 5 + 127) / 255
	//g6 := (g * 5 + 127) / 255
	//b6 := (b * 5 + 127)/ 255
	//cubeIndex := 16 + 36*r6 + 6*g6 +b6
	maxC := r 
	if g > maxC { maxC = g }
	if b > maxC { maxC = b }

	minC := r 
	if g < minC { minC = g }
	if b < minC { minC = b }

	avg := (r + g + b ) / 3

	if maxC - minC < 15 && avg > 15 {
  //check if it's close enough to gray
  //if abs(r-g) < 10 && abs(g-b) < 10 {
  	//avg := (r + g + b) / 3
    if avg > 245 {
		return 255
	}

	//if avg < 8 {
  	//	avg = 8
  	//}
  	//if avg > 238 {
  	//	avg = 238
  	//}
  	//return  232 + (avg-8)/10
	return  232 + ((avg-8) * 23 / 247)
  }	
    r6 := (r * 5 + 127) / 255
	g6 := (g * 5 + 127) / 255
	b6 := (b * 5 + 127)/ 255
    //fmt.Printf("RGB TO INDEX FROM COLOR HELPERS CODE (NOT TEST):  RGB=(%d,%d,%d)  | 256 = %d\n", r, g, b, 16 + 36*r6 + 6*g6 +b6)
	return 16 + 36*r6 + 6*g6 +b6

}
