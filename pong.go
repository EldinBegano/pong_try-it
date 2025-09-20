package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	width       = 40
	height      = 10
	paddleSize  = 3
	sleepMillis = 100
)

type Game struct {
	playerY int
	ballX   int
	ballY   int
	ballVX  int
	ballVY  int
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (g *Game) draw() {
	clearScreen()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x == 0 && y >= g.playerY && y < g.playerY+paddleSize {
				fmt.Print("|")
			} else if x == g.ballX && y == g.ballY {
				fmt.Print("O") 
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (g *Game) update() {
	g.ballX += g.ballVX
	g.ballY += g.ballVY

	if g.ballY <= 0 || g.ballY >= height-1 {
		g.ballVY *= -1
	}

	if g.ballX == 1 && g.ballY >= g.playerY && g.ballY < g.playerY+paddleSize {
		g.ballVX *= -1
	}

	if g.ballX <= 0 || g.ballX >= width-1 {
		g.ballX = width / 2
		g.ballY = height / 2
		g.ballVX = 1
		g.ballVY = 1
	}

	if g.ballY < g.playerY {
		g.playerY--
	} else if g.ballY > g.playerY+paddleSize-1 {
		g.playerY++
	}

	if g.playerY < 0 {
		g.playerY = 0
	}
	if g.playerY > height-paddleSize {
		g.playerY = height - paddleSize
	}
}

func main() {
	game := Game{
		playerY: height / 2,
		ballX:   width / 2,
		ballY:   height / 2,
		ballVX:  1,
		ballVY:  1,
	}

	for {
		game.draw()
		game.update()
		time.Sleep(sleepMillis * time.Millisecond)
	}
}
