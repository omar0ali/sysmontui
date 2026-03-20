package controls

import (
	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type Controls struct {
	defaultListLimit int
	listOfControls   []string
	sceneControl     interfaces.SceneControl
}

func Init(s interfaces.SceneControl) *Controls {
	listOfControls := []string{
		"[ESC] Exit",
		"[0] Home",
		"[1] Memory Info",
		"[2] CPU Info",
		"[3] Processors",
	}

	return &Controls{
		listOfControls:   listOfControls,
		defaultListLimit: len(listOfControls),
		sceneControl:     s,
	}
}

func (c *Controls) MenuAddToList(s string) {
	c.listOfControls = append(c.listOfControls, s)
}

func (c *Controls) MenuResetList() {
	if len(c.listOfControls) > c.defaultListLimit {
		c.listOfControls = c.listOfControls[:c.defaultListLimit]
	}
}

func (c *Controls) Update(d float64) {}
func (c *Controls) Render(s interfaces.ScreenControl) {
	w, h := s.Size()
	window.LineHorizontal(s, w, 3, tcell.RuneHLine)  // for ui
	window.LineVertical(s, h-3, 30, tcell.RuneVLine) // for ui

	start := 4
	for y, text := range c.listOfControls {
		window.Text(s, screentui.P(2, float64(start+y)), text)
	}
}

func (c *Controls) Events(ev *tcell.EventKey) {
	s := ev.Str()
	if len(s) == 1 && s[0] >= '0' && s[0] <= '4' {
		c.MenuResetList()

		page := int(s[0] - '0')
		switch page {
		case 0:
			c.sceneControl.SetScene("0")
		case 1:
			c.MenuAddToList("-------------------")
			c.MenuAddToList("[U] Update")
			c.MenuAddToList("[T] Toggle (MB, GB)")
			c.sceneControl.SetScene("1")
			// page 1
		case 2:
			c.sceneControl.SetScene("2")
			// page 2
		case 3:
			c.sceneControl.SetScene("3")
			// page 3
		}
	}
}
