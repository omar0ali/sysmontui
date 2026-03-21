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
	listOfLogs       []string
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

func (c *Controls) LogsAddToList(s string) {
	c.listOfLogs = append([]string{s}, c.listOfLogs...)
	if len(c.listOfLogs) > 4 {
		c.listOfLogs = c.listOfLogs[:4]
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
	// horizontal line for the title
	window.LineHorizontal(s, w, 3, tcell.RuneHLine) // for ui
	// left side vertical line
	window.LineVertical(s, h-3, 30, tcell.RuneVLine) // for ui

	// list of controls options
	startY := 4
	for y, text := range c.listOfControls {
		window.Text(s, screentui.P(2, float64(startY+y)), text)
	}

	// footer line for logs
	window.Text(s, screentui.P(33, float64(h)-14), "-- Logs --")
	window.LineHorizontalWithStartAndEnd(s, 31, w, h-13, tcell.RuneHLine) // for ui
	startY = h - 12
	for y, text := range c.listOfLogs {
		window.Text(s, screentui.P(34, float64(startY+y)), text)
	}
	// set start on the first line
	s.SetContent(32, h-12, '*')
}

func (c *Controls) Events(ev *tcell.EventKey) {
	s := ev.Str()
	if len(s) == 1 && s[0] >= '0' && s[0] <= '4' {
		c.MenuResetList()

		page := int(s[0] - '0')
		switch page {
		case 0:
			c.sceneControl.SetScene("0")
			c.LogsAddToList("Home Page")
		case 1:
			c.MenuAddToList("-------------------")
			c.MenuAddToList("[U] Update")
			c.MenuAddToList("[T] Toggle (MB, GB)")
			c.sceneControl.SetScene("1")
			c.LogsAddToList("Memory Page")
			// page 1
		case 2:
			c.sceneControl.SetScene("2")
			c.LogsAddToList("2")
			// page 2
		case 3:
			c.sceneControl.SetScene("3")
			c.LogsAddToList("3")
			// page 3
		}
	}
}
