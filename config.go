package main

type GameConfig struct {
	PaddleHeight        int
	BotSpeed            float64
	BotReactionDistance float64
}

var Easy = GameConfig{
	PaddleHeight:        6,
	BotSpeed:            0.8,
	BotReactionDistance: 25,
}

var Hard = GameConfig{
	PaddleHeight:        5,
	BotSpeed:            1.1,
	BotReactionDistance: 40,
}
