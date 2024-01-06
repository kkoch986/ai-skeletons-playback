package main

import (
	"os"

	"github.com/kkoch986/ai-skeletons-playback/command"
)

func main() {
	app := command.App()
	_ = app.Run(os.Args)
}
