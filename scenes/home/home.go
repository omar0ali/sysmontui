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

	// padding required due to the left section
	padding := 40
	w, _ := s.Size()
	window.TextWrap(s, screentui.P(float64(startXPos), 6), w-padding,
		`System monitoring tool with a page based UI for CPU, memory,
		and processes. Sysmon is a library I created for the sysmontui application.
		The goal is to build a terminal-based system monitor,
		similar to htop, but with a unique user interface.`,
	)

	// header
	window.Text(s, screentui.P(float64(startXPos), 12), "Features:")

	// list of features
	window.Text(s, screentui.P(float64(startXPos), 13), "- CPU usage information")
	window.Text(s, screentui.P(float64(startXPos), 14), "- Memory statistics")
	window.Text(s, screentui.P(float64(startXPos), 15), "- List of processes +(Sort, Order)")

}
func (h *Home) Events(ev tcell.Event) {}
