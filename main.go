package main

import (
	"github.com/nsf/termbox-go"
	"log"
)

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal("Failed to initialize termbox:", err)
	}
	defer termbox.Close()

	showStartScreen()
	chosenConfig := showDifficultyMenu()
	chosenTheme := showThemeMenu()

	game := &Game{}
	game.init(chosenConfig, Themes[chosenTheme])
	runGameLoop(game)
}

func showStartScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()

	instructions := []string{
		"╔═══════════════════════════════════════╗",
		"║                TERMINAL PONG          ║",
		"╠═══════════════════════════════════════╣",
		"║                                       ║",
		"║  You: GREEN paddle (left side)        ║",
		"║  Bot: RED paddle (right side)         ║",
		"║                                       ║",
		"║  Controls:                            ║",
		"║  ↑ Arrow Up - Move paddle up          ║",
		"║  ↓ Arrow Down - Move paddle down      ║",
		"║  ESC - Quit game                      ║",
		"║                                       ║",
		"║  First to 10 points wins!             ║",
		"║  Ball gets faster with each hit!      ║",
		"║  Game adapts to your terminal size!   ║",
		"║                                       ║",
		"║  Press SPACE to start...              ║",
		"╚═══════════════════════════════════════╝",
	}

	for i, line := range instructions {
		for j, ch := range line {
			termbox.SetCell(w/2-len(line)/2+j, h/2-len(instructions)/2+i, ch, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	termbox.Flush()

	for {
		if ev := termbox.PollEvent(); ev.Type == termbox.EventKey && ev.Key == termbox.KeySpace {
			break
		}
	}
}

func showDifficultyMenu() GameConfig {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()

	options := []string{
		"╔══════════════════════════╗",
		"║     Select Difficulty    ║",
		"╠══════════════════════════╣",
		"║  1. Easy                 ║",
		"║  2. Impossible           ║",
		"╚══════════════════════════╝",
	}

	for i, line := range options {
		for j, ch := range line {
			termbox.SetCell(w/2-len(line)/2+j, h/2-len(options)/2+i, ch, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	termbox.Flush()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Ch {
			case '1':
				return Easy
			case '2':
				return Hard
			}
		}
	}
}

func showThemeMenu() string {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()

	options := []string{
		"╔══════════════════════════╗",
		"║      Select Theme        ║",
		"╠══════════════════════════╣",
		"║  1. Black and White      ║",
		"║  2. Catppuccin           ║",
		"║  3. Normal  				║",
		"╚══════════════════════════╝",
	}

	for i, line := range options {
		for j, ch := range line {
			termbox.SetCell(w/2-len(line)/2+j, h/2-len(options)/2+i, ch, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	termbox.Flush()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Ch {
			case '1':
				return "theme-uno"
			case '2':
				return "theme-to"
			case '3':
				return "theme-tree"
			}
		}
	}
}
