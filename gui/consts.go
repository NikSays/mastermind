package gui

import (
	_ "embed"
)

// Glade GUI schema file
//
//go:embed gui.glade
var glade string

// 6 colors represented as
// RGB values between 0 and 1
var Colors = [6][3]float64{
	{0.929, 0.200, 0.231}, // Red
	{1.000, 0.639, 0.282}, // Orange
	{0.341, 0.890, 0.537}, // Yellow
	{0.384, 0.627, 0.917}, // Blue
	{0.752, 0.380, 0.796}, // Purple
	{0.964, 0.960, 0.956}, // White
}
