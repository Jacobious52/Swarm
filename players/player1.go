package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	color := os.Args[3]
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		line := ""
		fmt.Scanln(&line)
		count := strings.Count(line, color)
		var moves string
		for i := 0; i < count; i++ {
			moves += strconv.Itoa(rand.Intn(7))
		}
		fmt.Fprintln(os.Stdout, moves)
	}

}
