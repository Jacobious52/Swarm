package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	yellow = 'Y'
	blue   = 'B'
)

func print(board [][]byte, dst io.Writer) {
	str := ""
	for _, n := range board {
		for _, m := range n {
			str += string(m) + ","
		}
	}
	fmt.Fprintln(dst, str)
}

func read(src io.Reader) string {
	moves := ""
	fmt.Fscanln(src, &moves)
	return moves
}

func initBoard(w, h int) [][]byte {
	board := make([][]byte, w)
	for i := range board {
		board[i] = make([]byte, h)
	}

	rand.Seed(42)

	cellsP1, cellsP2 := 10, 10

	for cellsP1 > 0 && cellsP2 > 0 {
		for i := 0; i < w; i++ {
			for j := 0; j < h; j++ {
				board[i][j] = '0'

				if cellsP1 > 0 && rand.Intn(100) < 5 {
					board[i][j] = yellow
					cellsP1--
					continue
				}
				if cellsP2 > 0 && rand.Intn(100) < 5 {
					board[i][j] = blue
					cellsP2--
					continue
				}
			}
		}
	}
	return board
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("not enough args", len(os.Args))
		os.Exit(1)
	}

	w, _ := strconv.Atoi(os.Args[1])
	h, _ := strconv.Atoi(os.Args[2])
	p1, p2 := os.Args[3], os.Args[4]

	// init the players, p1 is yellow, p2 is blue
	player1 := exec.Command(p1, strconv.Itoa(w), strconv.Itoa(h), string(yellow))
	player2 := exec.Command(p2, strconv.Itoa(w), strconv.Itoa(h), string(blue))

	// create board
	board := initBoard(w, h)

	// setup player in out pipes
	player1In, _ := player1.StdinPipe()
	player1Out, _ := player1.StdoutPipe()
	defer player1In.Close()
	defer player1Out.Close()

	player2In, _ := player2.StdinPipe()
	player2Out, _ := player2.StdoutPipe()
	defer player2In.Close()
	defer player2Out.Close()

	// start the bot processes
	e1 := player1.Start()
	if e1 != nil {
		fmt.Println("player1", e1.Error())
		os.Exit(1)
	}
	e2 := player2.Start()
	if e2 != nil {
		fmt.Println("player1", e2.Error())
		os.Exit(1)
	}

	// start game loop
	gameOver := false
	for !gameOver {
		print(board, player1In)
		print(board, player2In)

		player1Moves := read(player1Out)
		player2Moves := read(player2Out)

		fmt.Println("p1", string(player1Moves))
		fmt.Println("p2", string(player2Moves))

		gameOver = false
		time.Sleep(1 * time.Second)
	}

	player1.Wait()
	player2.Wait()
}
