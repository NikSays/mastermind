package gui

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/niksays/mastermind/game"
)

// End previous game,
// start game setup,
// create new grid for solution
func newGame(gui *GuiElements, gs *GameState) bool {
	GameEnd(gui, gs)
	if gui.newGameTypePreset.GetActive() {
		gui.presetCodeDialog.Show()
	} else {
		startRandomGame(gui, gs)
	}

	gs.currentAttemptGrid = newAttemptGrid()
	gui.gameField.Attach(gs.currentAttemptGrid, 0, gs.attempt, 1, 1)
	gui.newGameDialog.Hide()
	return true
}

// Replace circle `pos` from grid with empty label.
// Set code element `pos` to -1
func ClearCode(pos *int, grid *gtk.Grid, code *[4]int) {
	if *pos <= 0 {
		return
	}
	*pos--
	child, err := grid.GetChildAt(*pos, 0)
	if err != nil {
		log.Panic("Can't get child")
	}
	child.ToWidget().Destroy()
	grid.Remove(child)
	label, err := gtk.LabelNew("")
	if err != nil {
		log.Panic("Can't create label")
	}
	code[*pos] = -1
	grid.Attach(label, *pos, 0, 1, 1)
	grid.ShowAll()
}

// Set element `pos` of `code` to `color`.
// Create circle with `color` in `gridPtr`.
// Activate `ok` button if `pos` reached 4
func SetClicked(color int, code *[4]int, pos *int, gridPtr **gtk.Grid, ok *gtk.Button) func(button *gtk.DrawingArea, e *gdk.Event) bool {
	return func(button *gtk.DrawingArea, e *gdk.Event) bool {
		if *pos >= 4 {
			return true
		}

		event := gdk.EventButtonNewFromEvent(e)
		if event.Type() != gdk.EVENT_BUTTON_PRESS {
			return true
		}

		child, err := (*gridPtr).GetChildAt(*pos, 0)
		if err != nil {
			log.Panic("Can't get child")
		}
		(*gridPtr).Remove(child)

		da, err := gtk.DrawingAreaNew()
		if err != nil {
			log.Panic("Can't create drawing area")
		}
		da.SetSizeRequest(50, 50)
		da.Connect("draw", SetColor(Colors[color][0], Colors[color][1], Colors[color][2]))
		(*gridPtr).Attach(da, *pos, 0, 1, 1)

		code[*pos] = color
		*pos++
		if *pos == 4 {
			ok.SetSensitive(true)
		}

		(*gridPtr).ShowAll()
		return true
	}
}

// Set color of drawing area
func SetColor(R float64, G float64, B float64) func(da *gtk.DrawingArea, cr *cairo.Context) bool {
	return func(da *gtk.DrawingArea, cr *cairo.Context) bool {
		cr.SetSourceRGB(R, G, B)
		cr.Arc(25, 25, 25, 0, 2*math.Pi)
		cr.Fill()
		return true
	}
}

// Create grid filled with empty labels
func newAttemptGrid() *gtk.Grid {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Panic("Can't create grid")
	}
	grid.SetColumnSpacing(5)
	for i := 0; i < 4; i++ {
		label, err := gtk.LabelNew("")
		if err != nil {
			log.Panic("Can't create label")
		}
		grid.Add(label)
	}
	return grid
}

func drawCorrectPattern(gui *GuiElements, gs *GameState) {
	gui.loseDialogCorrect.GetChildren().Foreach(func(item interface{}) {
		gui.loseDialogCorrect.Remove(item.(gtk.IWidget))
	})
	for i := 0; i < 4; i++ {
		c := gs.code[i]
		addCircleToGrid(gui.loseDialogCorrect, Colors[c])
	}
	gui.loseDialogCorrect.ShowAll()
}

func startRandomGame(gui *GuiElements, gs *GameState) {
	for i := 0; i < 4; i++ {
		gs.code[i] = rand.Intn(6)
	}
	fmt.Printf("Code: %v\n", gs.code)
	gui.solutionInput.SetVisible(true)
}

// Remove gaming field
func GameEnd(gui *GuiElements, gs *GameState) {
	gui.solutionInput.SetVisible(false)
	gui.gameField.GetChildren().FreeFull(func(item interface{}) {
		gui.gameField.Remove(item.(gtk.IWidget))
	})
	// Clear circles in preset dialog
	for i := 0; i < 4; i++ {
		ClearCode(&gs.currentCodePos, gui.presetCodeInput, &gs.code)
	}
	gs.Default()
}

func addCircleToGrid(grid *gtk.Grid, color [3]float64) {
	da, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Panic("Can't create DA")
	}
	da.Connect("draw", SetColor(color[0], color[1], color[2]))
	da.SetSizeRequest(50, 50)
	grid.Add(da)
}

// Return hints according to the algorithm
func createHintGrid(gs *GameState) *gtk.Grid {
	hintGrid, err := gtk.GridNew()
	if err != nil {
		log.Panic("Can't create grid")
	}
	hintGrid.SetColumnSpacing(5)
	hit, near := game.GiveHint(gs.solution, gs.code)
	for i := 0; i < hit; i++ {
		addCircleToGrid(hintGrid, Colors[0])
	}
	for i := 0; i < near; i++ {
		addCircleToGrid(hintGrid, Colors[5])
	}
	return hintGrid
}
