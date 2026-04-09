package main

import(
//	"fmt"
	"github.com/ph4mished/crayon"
)

func main(){
	word := crayon.Parse("[fg=red][0][reset]") //18 len

	for i:=1; i<=1000000; i++ {
	word.Sprint("ERROR: File Not Found So Go home and eat and sleep for tomorrow")
	//word.Sprint("Hi")
	}
	
}
