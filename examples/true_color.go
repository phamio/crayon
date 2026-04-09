//to test if true color fallback works
package main
import(
	"fmt"
	"github.com/ph4mished/crayon"
)

func main(){
//	col := crayon.Parse("[fg=#ff5fd7]HELLO [fg=#d7875f]WORLD[reset]")
	 rgb := crayon.Parse("[fg=rgb(222,230,240)]HELLO[reset]")
	 c_256 := crayon.Parse("[fg=255]HELLO[reset]")
	 hex := crayon.Parse("[fg=#875f00]HELLO[reset]")
	 
	fmt.Println("RGB: ", rgb.Sprint())
	fmt.Println("256: ", c_256.Sprint())
	fmt.Println("HEX: ", hex.Sprint())
}
