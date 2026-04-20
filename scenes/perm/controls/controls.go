package controls

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type LogsControl func(string)

type Controls struct {
	defaultListLimit int
	listOfControls   []string
	listOfLogs       []string
	sceneControl     interfaces.SceneControl

	lockMenu bool
}

func Init(s interfaces.SceneControl) *Controls {
	listOfControls := []string{
		"[1] Memory Info",
		"[2] CPU Status",
		"[3] Running Processes",
		"[ESC] Exit",
	}

	controls := &Controls{
		listOfControls:   listOfControls,
		defaultListLimit: len(listOfControls),
		sceneControl:     s,
	}

	// Memory View / Page
	// Additional menu
	controls.MenuAddToList("-------------------")
	controls.MenuAddToList("[T] Toggle (MB, GB)")

	// Logs
	controls.LogsAddToList("Memory View")

	return controls
}

func (c *Controls) LogsAddToList(s string) {
	limit := 9
	c.listOfLogs = append([]string{s}, c.listOfLogs...)
	if len(c.listOfLogs) > limit {
		c.listOfLogs = c.listOfLogs[:limit]
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
	s.Color(color.YellowGreen)
	w, h := s.Size()
	// horizontal line for the title
	window.LineHorizontal(s, w, 3, tcell.RuneHLine) // for ui
	// left side vertical line
	window.LineVertical(s, h-3, 30, tcell.RuneVLine) // for ui

	s.Color(color.White)

	// list of controls options
	startY := 4
	for y, text := range c.listOfControls {
		s.Color(color.White)                                     // default color
		if c.sceneControl.GetScene() == fmt.Sprintf("%d", y+1) { // starting from 1
			s.Color(color.YellowGreen)
		}
		window.Text(s, screentui.P(2, float64(startY+y)), text)
	}

	// footer line for logs
	window.LineHorizontalWithStartAndEnd(s, 31, w, h-15, tcell.RuneHLine) // for ui
	window.Text(s, screentui.P(33, float64(h)-14), "-- Logs --")
	window.LineHorizontalWithStartAndEnd(s, 31, w, h-13, tcell.RuneHLine) // for ui

	startY = h - 12
	for y, text := range c.listOfLogs {
		window.Text(s, screentui.P(34, float64(startY+y)), text)
	}
	// set start on the first line
	s.SetContent(32, h-12, '*')
}

func (c *Controls) Events(ev tcell.Event) {
	if c.lockMenu {
		return
	}
	switch ev := ev.(type) {
	case *tcell.EventKey:
		s := ev.Str()
		if len(s) == 1 && s[0] >= '1' && s[0] <= '3' {
			c.MenuResetList()

			page := int(s[0] - '0')
			switch page {
			case 1:
				if c.sceneControl.SetScene("1") {
					c.LogsAddToList("Memory View")
				}
				c.MenuAddToList("-------------------")
				c.MenuAddToList("[T] Toggle (MB, GB)")
			case 2:
				if c.sceneControl.SetScene("2") {
					c.LogsAddToList("CPU Status View")
				}
			case 3:
				if c.sceneControl.SetScene("3") {
					c.LogsAddToList("Running Processes View")
				}
				c.MenuAddToList("-------------------")
				c.MenuAddToList("[j] Scroll Up")
				c.MenuAddToList("[k] Scroll Down")
				c.MenuAddToList("[t] Toggle Asc/Desc")
				c.MenuAddToList("[T] Toggle Sort")
				c.MenuAddToList("[/] Search By Name")
				c.MenuAddToList("[K] Kill Process")
			}
		}
	}
}

func (c *Controls) Lock() {
	c.lockMenu = true
}

func (c *Controls) Unlock() {
	c.lockMenu = false
}
