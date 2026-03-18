package controls

import (
	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
)

type Controls struct{}

func Init() *Controls {
	return &Controls{}
}

func (c *Controls) Update(d float64)                  {}
func (c *Controls) Render(s interfaces.ScreenControl) {}
func (c *Controls) Events(ev *tcell.EventKey)         {}
