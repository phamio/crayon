package main
import(
	"github.com/ph4mished/crayon"
	//"fmt"
)

func main(){
	esc := crayon.Parse("[fg=red][<<fg=red>>][[fg=red]]ERROR[reset]")
	esc.Println()
}
