package main

import (
	"fmt"
	"os"
)

func main() {

	for {
		line := ""
		fmt.Scanln(&line)
		// do stuff
		//fmt.Println("hello", line)
		//f, _ := os.Create("ptest.txt")
		fmt.Fprintln(os.Stdout, "3")
	}

}
