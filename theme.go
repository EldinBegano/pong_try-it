package main


import "github.com/nsf/termbox-go"

type Theme struct {
	BorderColor  termbox.Attribute
	DividerColor termbox.Attribute
	PlayerColor  termbox.Attribute
	BotColor     termbox.Attribute
	BallColor    termbox.Attribute
	TextColor    termbox.Attribute
	BgColor      termbox.Attribute
}

var Themes = map[string]Theme{
	"theme-uno": {
		BorderColor:  termbox.ColorWhite,
		DividerColor: termbox.ColorWhite,
		PlayerColor:  termbox.ColorWhite,
		BotColor:     termbox.ColorWhite,
		BallColor:    termbox.ColorWhite | termbox.AttrBold,
		TextColor:    termbox.ColorWhite | termbox.AttrBold,
		BgColor:      termbox.ColorBlack,
	},
	"theme-to": {
		BorderColor:  termbox.ColorCyan,       
		DividerColor: termbox.ColorMagenta,     
		PlayerColor:  termbox.ColorBlue,       
		BotColor:     termbox.ColorYellow,  
		BallColor:    termbox.ColorGreen | termbox.AttrBold, 
		TextColor:    termbox.ColorWhite | termbox.AttrBold, 
		BgColor:      termbox.ColorBlack,   
	},
	"theme-tree": {
		BorderColor:  termbox.ColorWhite,
		DividerColor: termbox.ColorBlue | termbox.AttrBold,
		PlayerColor:  termbox.ColorBlue | termbox.AttrBold,
		BotColor:     termbox.ColorYellow | termbox.AttrBold,
		BallColor:    termbox.ColorWhite | termbox.AttrBold,
		TextColor:    termbox.ColorWhite | termbox.AttrBold,
		BgColor:      termbox.ColorBlack,
	},
	"theme-for": {
		BorderColor:  termbox.ColorCyan,
		DividerColor: termbox.ColorBlue | termbox.AttrBold,
		PlayerColor:  termbox.ColorCyan | termbox.AttrBold,
		BotColor:     termbox.ColorMagenta | termbox.AttrBold,
		BallColor:    termbox.ColorWhite | termbox.AttrBold,
		TextColor:    termbox.ColorCyan | termbox.AttrBold,
		BgColor:      termbox.ColorBlue,
	},
}