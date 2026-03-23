package controls

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type LogsAddToList func(string)

type Controls struct {
	defaultListLimit int
	listOfControls   []string
	listOfLogs       []string
	sceneControl     interfaces.SceneControl
}

func Init(s interfaces.SceneControl) *Controls {
	listOfControls := []string{
		"[0] Home",
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

	controls.LogsAddToList("Home Page")

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
		s.Color(color.White)
		if c.sceneControl.GetScene() == fmt.Sprintf("%d", y) {
			s.Color(color.YellowGreen)
			window.Text(s, screentui.P(2, float64(startY+y)), text)
			continue
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
	switch ev := ev.(type) {
	case *tcell.EventKey:
		s := ev.Str()
		if len(s) == 1 && s[0] >= '0' && s[0] <= '3' {
			c.MenuResetList()

			page := int(s[0] - '0')
			switch page {
			case 0:
				c.sceneControl.SetScene("0")
				c.LogsAddToList("Home Page")
			case 1:
				c.MenuAddToList("-------------------")
				c.MenuAddToList("[T] Toggle (MB, GB)")
				c.sceneControl.SetScene("1")
				c.LogsAddToList("Memory Page")
			case 2:
				c.sceneControl.SetScene("2")
				c.LogsAddToList("CPU Status Page")
			case 3:
				c.sceneControl.SetScene("3")
				c.LogsAddToList("Running Processes Page")
				c.MenuAddToList("-------------------")
				c.MenuAddToList("[j] Scroll Up")
				c.MenuAddToList("[k] Scroll Down")
				c.MenuAddToList("[t] Toggle Asc/Desc")
				c.MenuAddToList("[T] Toggle Sort")
			}
		}
	}
}
