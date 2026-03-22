package home

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type Home struct{}

func Init() *Home {
	return &Home{}
}

func (h *Home) Update(d float64) {}
func (h *Home) Render(s interfaces.ScreenControl) {
	s.Color(color.White)
	window.Text(s, screentui.P(1, 1), "Page: Home") // title
	startXPos := 32
	window.Text(s, screentui.P(float64(startXPos), 4), "Sysmon TUI Application.")
}
func (h *Home) Events(ev tcell.Event) {}
