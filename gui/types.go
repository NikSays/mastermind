package gui

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type GuiElements struct {
	window *gtk.Window

	newGameButton *gtk.Button

	newGameDialog *gtk.Dialog
	// Set when preset code. Unset when random code
	newGameTypePreset   *gtk.RadioButton
	newGameDialogOk     *gtk.Button
	newGameDialogCancel *gtk.Button

	presetCodeDialog       *gtk.Dialog
	presetCodeInput        *gtk.Grid
	presetCodeButtons      [6]*gtk.DrawingArea
	presetCodeBackspace    *gtk.Button
	presetCodeDialogOk     *gtk.Button
	presetCodeDialogCancel *gtk.Button

	gameField            *gtk.Grid
	solutionInput        *gtk.Box
	solutionInputButtons [6]*gtk.DrawingArea
	solutionBackspace    *gtk.Button
	solutionSubmit       *gtk.Button
	winDialog            *gtk.Window
	winDialogOk          *gtk.Button

	loseDialog               *gtk.Window
	loseDialogCorrectPattern [4]*gtk.DrawingArea
	loseDialogOk             *gtk.Button
}

type GameState struct {
	currentCodePos     int
	code               [4]int
	attempt            int
	currentAttemptGrid *gtk.Grid
	currentSolutionPos int
	solution           [4]int
}

func (gs *GameState) Default() {
	gs.currentCodePos = 0
	gs.code = [4]int{-1, -1, -1, -1}
	gs.attempt = 0
	gs.currentAttemptGrid = nil
	gs.currentSolutionPos = 0
	gs.solution = [4]int{-1, -1, -1, -1}
}

func (gs *GameState) newSolution() {
	gs.attempt++
	gs.currentSolutionPos = 0
	gs.solution = [4]int{-1, -1, -1, -1}
}

type ErrorHandlingBuilder struct{ *gtk.Builder }

// Helper that panics if object was not found
// Reduces iferr boilerplate when preset objects are loaded once
// and are not destroyed
func (b *ErrorHandlingBuilder) MustGetObject(id string) glib.IObject {
	object, err := b.GetObject(id)
	if err != nil {
		log.Fatalf("Couldn't get element %s: %s", id, err)
	}
	return object
}
