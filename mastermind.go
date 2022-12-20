// A gui mastermind game implementation
package main

import (
	"math/rand"
	"time"

	"github.com/niksays/mastermind/gui"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	guiElements := gui.BuildGui()
	gs := &gui.GameState{}
	gui.StartGame(guiElements, gs)
}
