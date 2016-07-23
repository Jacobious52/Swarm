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
			str += string(m)
		}
	}
	fmt.Fprintln(dst, str)
}

func pprint(board [][]byte) {
	for _, n := range board {
		for _, m := range n {
			fmt.Printf("%c ", m)
		}
		fmt.Println()
	}
}

func read(src io.Reader) string {
	moves := ""
	fmt.Fscanln(src, &moves)
	return moves
}

func initBoard(w, h, c int) [][]byte {
	board := make([][]byte, w)
	for i := range board {
		board[i] = make([]byte, h)
		for j := range board[i] {
			board[i][j] = '0'
		}
	}

	c1, c2 := c, c

	rand.Seed(time.Now().UTC().UnixNano())

	for c1 > 0 && c2 > 0 {
		for i := 0; i < w; i++ {
			for j := 0; j < h; j++ {
				if c1 > 0 && rand.Intn(100) < 5 {
					board[i][j] = yellow
					c1--
				} else if c2 > 0 && rand.Intn(100) < 5 {
					board[i][j] = blue
					c2--
				}
			}
		}
	}
	return board
}

func readMove(b [][]byte, w io.Writer, r io.Reader, c chan string) {
	print(b, w)
	c <- read(r)
}

type cell struct {
	ox, oy int
	x, y   int
	c      byte
}

func position(move byte) (int, int) {
	switch move {
	case '0':
		return 0, -1
	case '1':
		return 1, -1
	case '2':
		return 1, 0
	case '3':
		return 1, 1
	case '4':
		return 0, 1
	case '5':
		return -1, 1
	case '6':
		return -1, 0
	case '7':
		return -1, -1
	}
	return 0, 0
}

func clamp(x, y, w, h int) (int, int) {
	if x < 0 {
		x = 0
	} else if x > w-1 {
		x = w - 1
	}

	if y < 0 {
		y = 0
	} else if y > h-1 {
		y = h - 1
	}

	return x, y
}

func performMoves(board [][]byte, moves string, color byte, cellChan chan []cell) {
	// sanity check for player move count compared to board count
	moveSet := moves
	moveIndex := 0
	for i := range board {
		for j := range board[i] {
			if board[i][j] == color {
				moveIndex++
			}
		}
	}

	cells := make([]cell, moveIndex)

	if moveIndex != len(moveSet) {
		fmt.Fprintln(os.Stderr, "warning", string(color), "incorrect number of moves. given", len(moveSet), "should be", moveIndex, "forfiet move!")
		moveIndex = 0
		for i := range board {
			for j := range board[i] {
				if board[i][j] == color {
					cells[moveIndex] = cell{i, j, i, j, color}
					moveIndex++
				}
			}
		}
		cellChan <- cells
		return
	}

	moveIndex = 0

	for i, a := range board {
		for j, b := range a {
			if b == color {
				y, x := position(moveSet[moveIndex])
				nx, ny := clamp(i+x, j+y, len(board), len(a))
				cells[moveIndex] = cell{i, j, nx, ny, color}
				moveIndex++
			}
		}
	}

	cellChan <- cells
}

type cellpair struct {
	c *cell
	d bool
}

func updateBoard(board [][]byte, p1, p2 []cell) ([][]byte, int, bool) {
	// clear board
	for i := range board {
		for j := range board[i] {
			if board[i][j] == yellow || board[i][j] == blue {
				board[i][j] = '0'
			}
		}
	}

	died := false

	// check for blue collisions
	for i, c1 := range p1 {
		for j, c2 := range p2 {
			if c1.x == c2.x && c1.y == c2.y {
				if rand.Intn(10) < 5 && i < len(p1) {
					p1 = append(p1[:i], p1[i+1:]...)
					i--
					died = true
				} else {
					p2 = append(p2[:j], p2[j+1:]...)
					j--
					died = true
				}
			}
		}
	}

	// update
	for _, c := range p1 {
		board[c.x][c.y] = c.c
	}

	for _, c := range p2 {
		board[c.x][c.y] = c.c
	}

	// delete cells if captured

	count1, count2 := len(p1), len(p2)

	deleteChan := make(chan cellpair, count1+count2)

	// check for captures in p1
	for _, c := range p1 {
		go surrounded(board, c, blue, deleteChan)
	}

	// check for p2
	for _, c := range p2 {
		go surrounded(board, c, yellow, deleteChan)
	}

	for i := 0; i < count1+count2; i++ {
		cp := <-deleteChan
		if cp.d == true {
			board[cp.c.x][cp.c.y] = '0'
		}
	}

	winner := 0
	if float32(count1)/float32(count2) < 0.15 {
		winner = 1
	} else if float32(count2)/float32(count1) < 0.15 {
		winner = -1
	}

	return board, winner, died
}

type stack []cell

func (q stack) push(b cell) stack {
	return append(q, b)
}

func (q stack) pop() (stack, cell) {
	return q[:len(q)-1], q[len(q)-1]
}

func (q stack) empty() bool {
	return len(q) == 0
}

// check if position on board is surrounded by colour.
// returns true for surrounded, false for not surrounded
func surrounded(b [][]byte, start cell, blockingColor byte, deleteChan chan cellpair) {
	var visited []cell
	var fringe stack

	fringe = fringe.push(start)

	adj := [4][2]int{{0, -1}, {0, 1}, {-1, 0}, {0, 1}}

	for fringe.empty() == false {
		var current cell
		fringe, current = fringe.pop()
		visited = append(visited, current)

		if (current.x == 0 || current.y == 0 ||
			current.x == len(b)-1 || current.y == len(b[0])-1) &&
			current.c != blockingColor {

			deleteChan <- cellpair{&current, false}

			return
		}

		for _, p := range adj {
			x, y := current.x+p[0], current.y+p[1]
			if x >= 0 && y >= 0 && x < len(b) && y < len(b[0]) {
				next := cell{x, y, x, y, b[x][y]}

				in := false
				for _, v := range visited {
					if next.x == v.x && next.y == v.y {
						in = true
						break
					}
				}

				if in != true && b[x][y] != blockingColor {
					fringe = fringe.push(next)
				}
			}
		}
	}

	deleteChan <- cellpair{nil, false}
	return
}

func play() {
	if len(os.Args) != 6 {
		fmt.Println("not enough args", len(os.Args))
		os.Exit(1)
	}

	w, _ := strconv.Atoi(os.Args[1])
	h, _ := strconv.Atoi(os.Args[2])
	c, _ := strconv.Atoi(os.Args[3])
	p1, p2 := os.Args[4], os.Args[5]

	// init the players, p1 is yellow, p2 is blue
	player1 := exec.Command(p1, strconv.Itoa(w), strconv.Itoa(h), string(yellow))
	player2 := exec.Command(p2, strconv.Itoa(w), strconv.Itoa(h), string(blue))

	// create board
	board := initBoard(w, h, c)

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

	out, _ := os.Create("game.txt")
	defer out.Close()

	fmt.Fprintln(out, w, h)

	// start game loop
	gameOver := false
	moves, lastDied := 0, 0
	for !gameOver {

		player1Chan := make(chan string)
		player2Chan := make(chan string)

		go readMove(board, player1In, player1Out, player1Chan)
		go readMove(board, player2In, player2Out, player2Chan)

		player1MoveChan := make(chan []cell)
		player2MoveChan := make(chan []cell)

		// wait for input to be read and then async performMoves
		for i := 0; i < 2; i++ {
			select {
			case player1Moves := <-player1Chan:
				go performMoves(board, player1Moves, yellow, player1MoveChan)
			case player2Moves := <-player2Chan:
				go performMoves(board, player2Moves, blue, player2MoveChan)
			}
		}

		// wait for both to be done
		player1Cells := <-player1MoveChan
		player2Cells := <-player2MoveChan

		moves++

		var winner int
		var died bool
		board, winner, died = updateBoard(board, player1Cells, player2Cells)

		print(board, out)

		if winner == 1 {
			fmt.Fprintf(out, "%c", blue)
			fmt.Fprintln(os.Stderr, "player 1 wins")
			gameOver = true
		} else if winner == -1 {
			fmt.Fprintf(out, "%c", yellow)
			fmt.Fprintln(os.Stderr, "player 2 wins")
			gameOver = true
		}

		if died == true {
			lastDied = moves
		} else if died == false && moves-lastDied > 50 {
			fmt.Fprintf(out, "%c", 'S')
			fmt.Fprintln(os.Stderr, "stalemate")
			gameOver = true
		}
	}

	player1.Process.Kill()
	player2.Process.Kill()
}

func main() {
	play()
}
