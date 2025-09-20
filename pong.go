package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	PADDLE_HEIGHT = 5
	BOT_SPEED     = 0.8   
	BOT_REACTION_DISTANCE = 25   
)

type Ball struct {
	x, y   float64
	vx, vy float64
}

type Paddle struct {
	x, y     int
	targetY  float64   
	speed    float64   
}

type Game struct {
	width, height int
	ball          Ball
	leftPaddle    Paddle
	rightPaddle   Paddle
	leftScore     int
	rightScore    int
	running       bool
	lastUpdate    time.Time
}

func (g *Game) init() {
	g.width, g.height = termbox.Size()
	g.ball = Ball{
		x:  float64(g.width) / 2,
		y:  float64(g.height) / 2,
		vx: 0.6,
		vy: 0.3,
	}
	g.leftPaddle = Paddle{
		x:     2,
		y:     g.height/2 - PADDLE_HEIGHT/2,
		speed: 1.0,
	}
	g.rightPaddle = Paddle{
		x:       g.width - 3,
		y:       g.height/2 - PADDLE_HEIGHT/2,
		targetY: float64(g.height/2 - PADDLE_HEIGHT/2),
		speed:   BOT_SPEED,
	}
	g.leftScore = 0
	g.rightScore = 0
	g.running = true
	g.lastUpdate = time.Now()
}

func (g *Game) update() {
	now := time.Now()
	deltaTime := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	 
	g.ball.x += g.ball.vx * deltaTime * 40   
	g.ball.y += g.ball.vy * deltaTime * 40

	 
	if g.ball.y <= 1 {
		g.ball.y = 1
		g.ball.vy = math.Abs(g.ball.vy)   
	}
	if g.ball.y >= float64(g.height-2) {
		g.ball.y = float64(g.height - 2)
		g.ball.vy = -math.Abs(g.ball.vy)   
	}

	 
	if g.ball.x <= float64(g.leftPaddle.x+1) && 
		g.ball.y >= float64(g.leftPaddle.y-1) && 
		g.ball.y <= float64(g.leftPaddle.y+PADDLE_HEIGHT+1) &&
		g.ball.vx < 0 {
		g.ball.vx = math.Abs(g.ball.vx)   
		g.ball.x = float64(g.leftPaddle.x + 2)
		 
		paddleCenter := float64(g.leftPaddle.y + PADDLE_HEIGHT/2)
		hitOffset := (g.ball.y - paddleCenter) / float64(PADDLE_HEIGHT/2)
		g.ball.vy += hitOffset * 0.4
		 
		g.ball.vx *= 1.01   
		g.ball.vy *= 1.01
	}

	 
	if g.ball.x >= float64(g.rightPaddle.x-1) && 
		g.ball.y >= float64(g.rightPaddle.y-1) && 
		g.ball.y <= float64(g.rightPaddle.y+PADDLE_HEIGHT+1) &&
		g.ball.vx > 0 {
		g.ball.vx = -math.Abs(g.ball.vx)   
		g.ball.x = float64(g.rightPaddle.x - 2)
		 
		paddleCenter := float64(g.rightPaddle.y + PADDLE_HEIGHT/2)
		hitOffset := (g.ball.y - paddleCenter) / float64(PADDLE_HEIGHT/2)
		g.ball.vy += hitOffset * 0.4
		 
		g.ball.vx *= 1.01   
		g.ball.vy *= 1.01
	}

	 
	if g.ball.x < -2 {
		g.rightScore++
		g.resetBall(1)  
	}
	if g.ball.x > float64(g.width+2) {
		g.leftScore++
		g.resetBall(-1)  
	}

	 
	g.updateBotPaddle(deltaTime)
}

func (g *Game) updateBotPaddle(deltaTime float64) {
	 
	distanceToBall := g.ball.x - float64(g.rightPaddle.x)
	
	 
	if g.ball.vx > 0 && distanceToBall < BOT_REACTION_DISTANCE {
		 
		predictedY := g.predictBallY()
		
		 
		imperfection :=  2.0
		if distanceToBall > 15 {
			predictedY += math.Sin(float64(time.Now().UnixNano())*1e-9) * imperfection
		}
		
		g.rightPaddle.targetY = predictedY - float64(PADDLE_HEIGHT/2)
	} else {
		 
		centerY := float64(g.height/2 - PADDLE_HEIGHT/2)
		currentY := float64(g.rightPaddle.y)
		g.rightPaddle.targetY = currentY + (centerY-currentY)*0.01
	}
	
	 
	currentY := float64(g.rightPaddle.y)
	diff := g.rightPaddle.targetY - currentY
	
	if math.Abs(diff) > 0.1 {
		moveSpeed := g.rightPaddle.speed * deltaTime * 40   
		if diff > 0 {
			currentY += moveSpeed
		} else {
			currentY -= moveSpeed
		}
		
		 
		if currentY < 1 {
			currentY = 1
		} else if currentY > float64(g.height-PADDLE_HEIGHT-1) {
			currentY = float64(g.height-PADDLE_HEIGHT-1)
		}
		
		g.rightPaddle.y = int(math.Round(currentY))
	}
}

func (g *Game) predictBallY() float64 {
	 
	ballX := g.ball.x
	ballY := g.ball.y
	ballVx := g.ball.vx
	ballVy := g.ball.vy
	
	paddleX := float64(g.rightPaddle.x)
	
	if ballVx <= 0 {
		return ballY
	}
	
	 
	steps := 0
	maxSteps := 1000  
	
	for ballX < paddleX && steps < maxSteps {
		ballX += ballVx * 0.1
		ballY += ballVy * 0.1
		
		 
		if ballY <= 1 {
			ballY = 1
			ballVy = math.Abs(ballVy)
		} else if ballY >= float64(g.height-2) {
			ballY = float64(g.height - 2)
			ballVy = -math.Abs(ballVy)
		}
		
		steps++
	}
	
	return ballY
}

func (g *Game) resetBall(direction float64) {
	g.ball.x = float64(g.width) / 2
	g.ball.y = float64(g.height) / 2
	g.ball.vx = 0.6 * direction   
	g.ball.vy = 0.3 * (1 - 2*math.Mod(float64(time.Now().UnixNano()), 2))  
	
	 
	g.rightPaddle.targetY = float64(g.height/2 - PADDLE_HEIGHT/2)
}

func (g *Game) moveLeftPaddleUp(deltaTime float64) {
	if g.leftPaddle.y > 1 {
		moveAmount := g.leftPaddle.speed * deltaTime * 40   
		newY := float64(g.leftPaddle.y) - moveAmount
		if newY < 1 {
			newY = 1
		}
		g.leftPaddle.y = int(math.Round(newY))
	}
}

func (g *Game) moveLeftPaddleDown(deltaTime float64) {
	if g.leftPaddle.y < g.height-PADDLE_HEIGHT-1 {
		moveAmount := g.leftPaddle.speed * deltaTime * 40   
		newY := float64(g.leftPaddle.y) + moveAmount
		if newY > float64(g.height-PADDLE_HEIGHT-1) {
			newY = float64(g.height-PADDLE_HEIGHT-1)
		}
		g.leftPaddle.y = int(math.Round(newY))
	}
}

func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	
	 
	for x := 0; x < g.width; x++ {
		termbox.SetCell(x, 0, '═', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, g.height-1, '═', termbox.ColorWhite, termbox.ColorDefault)
	}
	for y := 0; y < g.height; y++ {
		termbox.SetCell(0, y, '║', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(g.width-1, y, '║', termbox.ColorWhite, termbox.ColorDefault)
	}
	
	 
	termbox.SetCell(0, 0, '╔', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, 0, '╗', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(0, g.height-1, '╚', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, g.height-1, '╝', termbox.ColorWhite, termbox.ColorDefault)

	 
	for y := 2; y < g.height-1; y += 3 {
		termbox.SetCell(g.width/2, y, '┊', termbox.ColorBlue, termbox.ColorDefault)
	}

	 
	for i := 0; i < PADDLE_HEIGHT; i++ {
		 
		if g.leftPaddle.y+i >= 1 && g.leftPaddle.y+i < g.height-1 {
			char := '█'
			color := termbox.ColorGreen
			if i == 0 || i == PADDLE_HEIGHT-1 {
				color = termbox.ColorGreen | termbox.AttrBold
			}
			termbox.SetCell(g.leftPaddle.x, g.leftPaddle.y+i, char, color, termbox.ColorDefault)
		}
		
		 
		if g.rightPaddle.y+i >= 1 && g.rightPaddle.y+i < g.height-1 {
			char := '█'
			color := termbox.ColorRed
			if i == 0 || i == PADDLE_HEIGHT-1 {
				color = termbox.ColorRed | termbox.AttrBold
			}
			termbox.SetCell(g.rightPaddle.x, g.rightPaddle.y+i, char, color, termbox.ColorDefault)
		}
	}

	 
	bx, by := int(math.Round(g.ball.x)), int(math.Round(g.ball.y))
	if bx >= 1 && bx < g.width-1 && by >= 1 && by < g.height-1 {
		termbox.SetCell(bx, by, '●', termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault)
		
		 
		trailX := bx - int(math.Copysign(1, g.ball.vx))
		if trailX >= 1 && trailX < g.width-1 {
			termbox.SetCell(trailX, by, '○', termbox.ColorYellow, termbox.ColorDefault)
		}
	}

	 
	scoreText := fmt.Sprintf("PLAYER: %02d  │  BOT: %02d", g.leftScore, g.rightScore)
	for i, ch := range scoreText {
		color := termbox.ColorWhite | termbox.AttrBold
		if i >= 8 && i <= 9 {  
			color = termbox.ColorGreen | termbox.AttrBold
		} else if i >= 18 && i <= 19 {  
			color = termbox.ColorRed | termbox.AttrBold
		}
		termbox.SetCell(g.width/2-len(scoreText)/2+i, 1, ch, color, termbox.ColorDefault)
	}

	 
	controlText := "↑↓ Move Paddle  │  ESC Quit"
	for i, ch := range controlText {
		termbox.SetCell(g.width/2-len(controlText)/2+i, g.height+2, ch, termbox.ColorCyan, termbox.ColorDefault)
	}
	
	 
	speed := math.Sqrt(g.ball.vx*g.ball.vx + g.ball.vy*g.ball.vy)
	speedText := fmt.Sprintf("Speed: %.1f", speed)
	for i, ch := range speedText {
		termbox.SetCell(2+i, g.height+1, ch, termbox.ColorMagenta, termbox.ColorDefault)
	}

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal("Failed to initialize termbox:", err)
	}
	defer termbox.Close()

	 
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	instructions := []string{
		"╔═══════════════════════════════════════╗",
		"║           TERMINAL PONG               ║",
		"╠═══════════════════════════════════════╣",
		"║                                       ║",
		"║  You: GREEN paddle (left side)        ║",
		"║  Bot: RED paddle (right side)         ║",
		"║                                       ║",
		"║  Controls:                            ║",
		"║  ↑ Arrow Up    - Move paddle up       ║",
		"║  ↓ Arrow Down  - Move paddle down     ║", 
		"║  ESC           - Quit game            ║",
		"║                                       ║",
		"║  First to 10 points wins!             ║",
		"║  Ball gets faster with each hit!      ║",
		"║  Game adapts to your terminal size!   ║",
		"║                                       ║",
		"║         Press SPACE to start...       ║",
		"╚═══════════════════════════════════════╝",
	}

	for i, line := range instructions {
		for j, ch := range line {
			color := termbox.ColorWhite
			if i == 1 {
				color = termbox.ColorYellow | termbox.AttrBold
			} else if i >= 7 && i <= 10 {
				color = termbox.ColorCyan
			} else if i == 16 {
				color = termbox.ColorGreen | termbox.AttrBold
			}
			termbox.SetCell(w/2-len(line)/2+j, h/2-len(instructions)/2+i, ch, color, termbox.ColorDefault)
		}
	}
	termbox.Flush()

	 
	for {
		if ev := termbox.PollEvent(); ev.Type == termbox.EventKey && ev.Key == termbox.KeySpace {
			break
		}
	}

	game := &Game{}
	game.init()

	 
	keyStates := make(map[termbox.Key]bool)

	 
	go func() {
		for game.running {
			ev := termbox.PollEvent()
			switch ev.Type {
			case termbox.EventResize:
				game.width, game.height = ev.Width, ev.Height
				 
				game.rightPaddle.x = game.width - 3
				 
				if game.leftPaddle.y > game.height-PADDLE_HEIGHT-1 {
					game.leftPaddle.y = game.height - PADDLE_HEIGHT - 1
				}
				if game.rightPaddle.y > game.height-PADDLE_HEIGHT-1 {
					game.rightPaddle.y = game.height - PADDLE_HEIGHT - 1
				}
				 
				game.rightPaddle.targetY = float64(game.height/2 - PADDLE_HEIGHT/2)
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					game.running = false
				case termbox.KeyArrowUp:
					keyStates[termbox.KeyArrowUp] = true
				case termbox.KeyArrowDown:
					keyStates[termbox.KeyArrowDown] = true
				}
			}
		}
	}()

	ticker := time.NewTicker(20 * time.Millisecond)  
	defer ticker.Stop()

	for game.running {
		select {
		case <-ticker.C:
			now := time.Now()
			deltaTime := now.Sub(game.lastUpdate).Seconds()
			
			 
			if keyStates[termbox.KeyArrowUp] {
				game.moveLeftPaddleUp(deltaTime)
				 
				go func() {
					time.Sleep(50 * time.Millisecond)
					keyStates[termbox.KeyArrowUp] = false
				}()
			}
			if keyStates[termbox.KeyArrowDown] {
				game.moveLeftPaddleDown(deltaTime)
				go func() {
					time.Sleep(50 * time.Millisecond)
					keyStates[termbox.KeyArrowDown] = false
				}()
			}
			
			game.update()
			game.draw()
			
			 
			if game.leftScore >= 10 || game.rightScore >= 10 {
				 
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				
				var winner string
				var winnerColor termbox.Attribute
				if game.leftScore >= 10 {
					winner = "🎉 VICTORY! YOU WIN! 🎉"
					winnerColor = termbox.ColorGreen | termbox.AttrBold
				} else {
					winner = "💀 DEFEAT! BOT WINS! 💀"
					winnerColor = termbox.ColorRed | termbox.AttrBold
				}
				
				gameOverText := []string{
					"╔═══════════════════════════════════════╗",
					"║              GAME OVER                ║",
					"╠═══════════════════════════════════════╣",
					"║                                       ║",
					winner,
					"║                                       ║",
					fmt.Sprintf("║      Final Score: %02d - %02d             ║", game.leftScore, game.rightScore),
					"║                                       ║",
					"║         Press ESC to exit             ║",
					"╚═══════════════════════════════════════╝",
				}
				
				for i, line := range gameOverText {
					color := termbox.ColorWhite
					if i == 1 {
						color = termbox.ColorYellow | termbox.AttrBold
					} else if i == 4 {
						color = winnerColor
					} else if i == 6 {
						color = termbox.ColorCyan
					}
					for j, ch := range line {
						termbox.SetCell(game.width/2-len(line)/2+j, game.height/2-len(gameOverText)/2+i, ch, color, termbox.ColorDefault)
					}
				}
				
				for i, line := range gameOverText {
					color := termbox.ColorWhite
					if i == 1 {
						color = termbox.ColorYellow | termbox.AttrBold
					} else if i == 4 {
						color = winnerColor
					} else if i == 6 {
						color = termbox.ColorCyan
					}
					for j, ch := range line {
						termbox.SetCell(game.width/2-len(line)/2+j, game.height/2-len(gameOverText)/2+i, ch, color, termbox.ColorDefault)
					}
				}
				termbox.Flush()
				
				 
				for {
					if ev := termbox.PollEvent(); ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
						game.running = false
						break
					}
				}
			}
		}
	}
}
