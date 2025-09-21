package main

import (
	"fmt"
	"math"
	"time"

	"github.com/nsf/termbox-go"
)

type Ball struct {
	x, y   float64
	vx, vy float64
}

type Paddle struct {
	x, y    int
	targetY float64
	speed   float64
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
	config        GameConfig
	theme         Theme
}

func (g *Game) init(config GameConfig, theme Theme) {
	g.config = config
	g.theme = theme
	g.width, g.height = termbox.Size()
	g.ball = Ball{
		x:  float64(g.width) / 2,
		y:  float64(g.height) / 2,
		vx: 0.6,
		vy: 0.3,
	}
	g.leftPaddle = Paddle{
		x:     2,
		y:     g.height/2 - g.config.PaddleHeight/2,
		speed: 1.0,
	}
	g.rightPaddle = Paddle{
		x:       g.width - 3,
		y:       g.height/2 - g.config.PaddleHeight/2,
		targetY: float64(g.height/2 - g.config.PaddleHeight/2),
		speed:   g.config.BotSpeed,
	}
	g.leftScore = 0
	g.rightScore = 0
	g.running = true
	g.lastUpdate = time.Now()
}

func (g *Game) reset() {
	g.width, g.height = termbox.Size()
	g.ball = Ball{
		x:  float64(g.width) / 2,
		y:  float64(g.height) / 2,
		vx: 0.6,
		vy: 0.3,
	}
	g.leftPaddle.y = g.height/2 - g.config.PaddleHeight/2
	g.rightPaddle.x = g.width - 3
	g.rightPaddle.y = g.height/2 - g.config.PaddleHeight/2
	g.rightPaddle.targetY = float64(g.height/2 - g.config.PaddleHeight/2)
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
		g.ball.y <= float64(g.leftPaddle.y+g.config.PaddleHeight+1) &&
		g.ball.vx < 0 {
		g.ball.vx = math.Abs(g.ball.vx)
		g.ball.x = float64(g.leftPaddle.x + 2)
		paddleCenter := float64(g.leftPaddle.y + g.config.PaddleHeight/2)
		hitOffset := (g.ball.y - paddleCenter) / float64(g.config.PaddleHeight/2)
		g.ball.vy += hitOffset * 0.4
		g.ball.vx *= 1.01
		g.ball.vy *= 1.01
	}

	if g.ball.x >= float64(g.rightPaddle.x-1) &&
		g.ball.y >= float64(g.rightPaddle.y-1) &&
		g.ball.y <= float64(g.rightPaddle.y+g.config.PaddleHeight+1) &&
		g.ball.vx > 0 {
		g.ball.vx = -math.Abs(g.ball.vx)
		g.ball.x = float64(g.rightPaddle.x - 2)
		paddleCenter := float64(g.rightPaddle.y + g.config.PaddleHeight/2)
		hitOffset := (g.ball.y - paddleCenter) / float64(g.config.PaddleHeight/2)
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

	if g.ball.vx > 0 && distanceToBall < g.config.BotReactionDistance {
		predictedY := g.predictBallY()
		imperfection := 2.0
		if distanceToBall > 15 {
			predictedY += math.Sin(float64(time.Now().UnixNano())*1e-9) * imperfection
		}
		g.rightPaddle.targetY = predictedY - float64(g.config.PaddleHeight/2)
	} else {
		centerY := float64(g.height/2 - g.config.PaddleHeight/2)
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
		} else if currentY > float64(g.height-g.config.PaddleHeight-1) {
			currentY = float64(g.height - g.config.PaddleHeight - 1)
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
	g.ball.vy = 0.3
	g.rightPaddle.targetY = float64(g.height/2 - g.config.PaddleHeight/2)
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
	if g.leftPaddle.y < g.height-g.config.PaddleHeight-1 {
		moveAmount := g.leftPaddle.speed * deltaTime * 40
		newY := float64(g.leftPaddle.y) + moveAmount
		if newY > float64(g.height-g.config.PaddleHeight-1) {
			newY = float64(g.height - g.config.PaddleHeight - 1)
		}
		g.leftPaddle.y = int(math.Round(newY))
	}
}

func (g *Game) draw() {
	termbox.Clear(g.theme.TextColor, g.theme.BgColor)

	for x := 0; x < g.width; x++ {
		termbox.SetCell(x, 0, '═', g.theme.BorderColor, g.theme.BgColor)
		termbox.SetCell(x, g.height-1, '═', g.theme.BorderColor, g.theme.BgColor)
	}
	for y := 0; y < g.height; y++ {
		termbox.SetCell(0, y, '║', g.theme.BorderColor, g.theme.BgColor)
		termbox.SetCell(g.width-1, y, '║', g.theme.BorderColor, g.theme.BgColor)
	}
	termbox.SetCell(0, 0, '╔', g.theme.BorderColor, g.theme.BgColor)
	termbox.SetCell(g.width-1, 0, '╗', g.theme.BorderColor, g.theme.BgColor)
	termbox.SetCell(0, g.height-1, '╚', g.theme.BorderColor, g.theme.BgColor)
	termbox.SetCell(g.width-1, g.height-1, '╝', g.theme.BorderColor, g.theme.BgColor)

	for y := 2; y < g.height-1; y += 3 {
		termbox.SetCell(g.width/2, y, '┊', g.theme.DividerColor, g.theme.BgColor)
	}

	for i := 0; i < g.config.PaddleHeight; i++ {
		if g.leftPaddle.y+i >= 1 && g.leftPaddle.y+i < g.height-1 {
			termbox.SetCell(g.leftPaddle.x, g.leftPaddle.y+i, '█', g.theme.PlayerColor, g.theme.BgColor)
		}
		if g.rightPaddle.y+i >= 1 && g.rightPaddle.y+i < g.height-1 {
			termbox.SetCell(g.rightPaddle.x, g.rightPaddle.y+i, '█', g.theme.BotColor, g.theme.BgColor)
		}
	}

	bx, by := int(math.Round(g.ball.x)), int(math.Round(g.ball.y))
	if bx >= 1 && bx < g.width-1 && by >= 1 && by < g.height-1 {
		termbox.SetCell(bx, by, '●', g.theme.BallColor, g.theme.BgColor)
	}

	scoreText := fmt.Sprintf("PLAYER: %02d  │  BOT: %02d", g.leftScore, g.rightScore)
	for i, ch := range scoreText {
		termbox.SetCell(g.width/2-len(scoreText)/2+i, 1, ch, g.theme.TextColor, g.theme.BgColor)
	}

	termbox.Flush()
}

func runGameLoop(game *Game) {
	keyStates := make(map[termbox.Key]bool)

	go func() {
		for game.running {
			ev := termbox.PollEvent()
			switch ev.Type {
			case termbox.EventResize:
				game.width, game.height = ev.Width, ev.Height
				game.rightPaddle.x = game.width - 3
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
		<-ticker.C
		now := time.Now()
		deltaTime := now.Sub(game.lastUpdate).Seconds()

		if keyStates[termbox.KeyArrowUp] {
			game.moveLeftPaddleUp(deltaTime)
			keyStates[termbox.KeyArrowUp] = false
		}
		if keyStates[termbox.KeyArrowDown] {
			game.moveLeftPaddleDown(deltaTime)
			keyStates[termbox.KeyArrowDown] = false
		}

		game.update()
		game.draw()

		if game.leftScore >= 10 || game.rightScore >= 10 {
			game.running = false
			if showGameOverScreen(game) {
				game.reset()
			}
		}
	}
}

func showGameOverScreen(g *Game) bool {
	termbox.Clear(g.theme.TextColor, g.theme.BgColor)
	w, h := g.width, g.height

	var winner string
	var color termbox.Attribute
	if g.leftScore >= 10 {
		winner = "YOU WIN!"
		color = g.theme.PlayerColor
	} else {
		winner = "BOT WINS!"
		color = g.theme.BotColor
	}

	lines := []string{
		"╔══════════════════════════╗",
		"║        GAME OVER         ║",
		"╠══════════════════════════╣",
		fmt.Sprintf("║   %s   ║", winner),
		fmt.Sprintf("║   Final Score %02d - %02d   ║", g.leftScore, g.rightScore),
		"║                          ║",
		"║ Press R to restart        ║",
		"║ Press ESC to exit         ║",
		"╚══════════════════════════╝",
	}

	for i, line := range lines {
		for j, ch := range line {
			col := g.theme.TextColor
			if i == 3 {
				col = color | termbox.AttrBold
			}
			termbox.SetCell(w/2-len(line)/2+j, h/2-len(lines)/2+i, ch, col, g.theme.BgColor)
		}
	}
	termbox.Flush()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyEsc:
				return false
			}
			switch ev.Ch {
			case 'r', 'R':
				return true
			}
		}
	}
}
