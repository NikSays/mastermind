package gui

import (
	"fmt"
	"log"
	"reflect"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// Create GUI struct
// and bind basic function, like closing
func BuildGui() (gui *GuiElements) {
	gtk.Init(nil)
	nativeBuilder, _ := gtk.BuilderNew()
	builder := ErrorHandlingBuilder{nativeBuilder}
	err := builder.AddFromString(glade)
	if err != nil {
		log.Fatal("Couldn't open gui.glade: ", err)
	}

	gui = &GuiElements{}
	gui.window = builder.MustGetObject("window").(*gtk.Window)
	gui.window.Connect("destroy", func() bool {
		fmt.Println("close")
		gtk.MainQuit()
		return true
	})

	gui.newGameButton = builder.MustGetObject("new_game_button").(*gtk.Button)
	gui.newGameButton.Connect("clicked", func() bool {
		gui.newGameDialog.Show()
		return true
	})

	gui.newGameDialog = builder.MustGetObject("new_game_dialog").(*gtk.Dialog)
	gui.newGameTypePreset = builder.MustGetObject("type_select_preset").(*gtk.RadioButton)
	gui.newGameDialogOk = builder.MustGetObject("type_select_ok").(*gtk.Button)
	gui.newGameDialogCancel = builder.MustGetObject("type_select_cancel").(*gtk.Button)
	gui.newGameDialog.Connect("delete-event", func() bool {
		gui.newGameDialog.Hide()
		return true
	})
	gui.newGameDialogCancel.Connect("clicked", func() bool {
		gui.newGameDialog.Hide()
		return true
	})

	gui.presetCodeDialog = builder.MustGetObject("preset_dialog").(*gtk.Dialog)
	gui.presetCodeInput = builder.MustGetObject("preset_code").(*gtk.Grid)
	gui.presetCodeBackspace = builder.MustGetObject("preset_delete_button").(*gtk.Button)
	gui.presetCodeDialogOk = builder.MustGetObject("preset_ok").(*gtk.Button)
	gui.presetCodeDialogCancel = builder.MustGetObject("preset_cancel").(*gtk.Button)
	for i := 0; i < 6; i++ {
		id := fmt.Sprintf("preset_input_%d", i+1)
		button := builder.MustGetObject(id).(*gtk.DrawingArea)
		// Make clickable
		button.AddEvents(int(gdk.BUTTON_PRESS_MASK))
		button.Connect("draw", SetColor(Colors[i][0], Colors[i][1], Colors[i][2]))
		gui.presetCodeButtons[i] = button
	}

	gui.gameField = builder.MustGetObject("solution_grid").(*gtk.Grid)
	gui.solutionInput = builder.MustGetObject("solution_input").(*gtk.Box)
	gui.solutionSubmit = builder.MustGetObject("submit_button").(*gtk.Button)
	gui.solutionBackspace = builder.MustGetObject("delete_button").(*gtk.Button)
	for i := 0; i < 6; i++ {
		id := fmt.Sprintf("input_color_%d", i+1)
		button := builder.MustGetObject(id).(*gtk.DrawingArea)
		// Make clickable
		button.AddEvents(int(gdk.BUTTON_PRESS_MASK))
		button.Connect("draw", SetColor(Colors[i][0], Colors[i][1], Colors[i][2]))
		gui.solutionInputButtons[i] = button
	}

	gui.winDialog = builder.MustGetObject("win").(*gtk.Window)
	gui.winDialogOk = builder.MustGetObject("win_ok").(*gtk.Button)
	gui.winDialog.Connect("delete-event", func() bool {
		gui.winDialog.Hide()
		return true
	})
	gui.winDialogOk.Connect("clicked", func() bool {
		gui.winDialog.Hide()
		return true
	})

	gui.loseDialog = builder.MustGetObject("lose").(*gtk.Window)
	gui.loseDialogCorrect = builder.MustGetObject("lose_correct_pattern").(*gtk.Grid)
	gui.loseDialogOk = builder.MustGetObject("lose_ok").(*gtk.Button)
	gui.loseDialog.Connect("delete-event", func() bool {
		gui.loseDialog.Hide()
		return true
	})
	gui.loseDialogOk.Connect("clicked", func() bool {
		gui.loseDialog.Hide()
		return true
	})

	return gui
}

func StartGame(gui *GuiElements, gs *GameState) {
	gs.Default()

	gui.newGameDialogOk.Connect("clicked", func() bool {
		newGame(gui, gs)
		return true
	})

	for i, button := range gui.presetCodeButtons {
		button.Connect("button-press-event", SetClicked(i, &gs.code, &gs.currentCodePos, &gui.presetCodeInput, gui.presetCodeDialogOk))
	}
	for i, button := range gui.solutionInputButtons {
		button.Connect("button-press-event", SetClicked(i, &gs.solution, &gs.currentSolutionPos, &gs.currentAttemptGrid, gui.solutionSubmit))
	}
	gui.solutionBackspace.Connect("clicked", func() bool {
		ClearCode(&gs.currentSolutionPos, gs.currentAttemptGrid, &gs.solution)
		gui.solutionSubmit.SetSensitive(false)
		return true
	})

	gui.solutionSubmit.Connect("clicked", func() bool {
		gui.solutionSubmit.SetSensitive(false)
		// Check game end
		if reflect.DeepEqual(gs.solution, gs.code) {
			GameEnd(gui, gs)
			gui.winDialog.Show()
			return true
		}
		if gs.attempt == 5 {
			drawCorrectPattern(gui, gs)
			GameEnd(gui, gs)
			gui.loseDialog.Show()
			return true
		}
		// Show hint
		hintGrid := createHintGrid(gs)
		gui.gameField.Attach(hintGrid, 1, gs.attempt, 1, 1)
		// Shift to new solution
		gs.newSolution()
		gs.currentAttemptGrid = newAttemptGrid()
		gui.gameField.Attach(gs.currentAttemptGrid, 0, gs.attempt, 1, 1)
		gui.gameField.ShowAll()
		return true
	})

	gui.presetCodeDialog.Connect("delete-event", func() bool {
		gui.presetCodeDialog.Hide()
		gui.presetCodeDialogOk.SetSensitive(false)
		return true
	})
	gui.presetCodeDialogCancel.Connect("clicked", func() bool {
		gui.presetCodeDialog.Hide()
		gui.presetCodeDialogOk.SetSensitive(false)
		return true
	})
	gui.presetCodeDialogOk.Connect("clicked", func() bool {
		fmt.Printf("Code: %v\n", gs.code)
		gui.presetCodeDialog.Hide()
		gui.presetCodeDialogOk.SetSensitive(false)
		gui.solutionInput.SetVisible(true)
		return true
	})

	gui.presetCodeBackspace.Connect("clicked", func() bool {
		ClearCode(&gs.currentCodePos, gui.presetCodeInput, &gs.code)
		gui.presetCodeDialogOk.SetSensitive(false)
		return true
	})

	gui.window.Show()
	gtk.Main()
}
