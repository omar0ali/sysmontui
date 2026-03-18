package interfaces

import (
	"github.com/gdamore/tcell/v3"
)

type EntityActions interface {
	Update(d float64)
	Render(ScreenControl)
	Events(ev *tcell.EventKey)
}
