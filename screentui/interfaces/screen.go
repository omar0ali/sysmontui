package interfaces

import "github.com/gdamore/tcell/v3/color"

type ScreenControl interface {
	Color(color.Color)
	DefaultColor()
	SetContent(x, y int, r rune)
	Size() (int, int)
	Exit()
}
