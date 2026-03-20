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
}

func Init() *Controls {
	listOfControls := []string{
		"Controls: ",
		"[ESC] Exit",
		"[1] CPU Info",
		"[2] Memory Info",
		"|",
	}

	return &Controls{
		listOfControls:   listOfControls,
		defaultListLimit: len(listOfControls),
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
	window.LineHorizontal(s, w, h-3, tcell.RuneHLine)

	spaceInBetween := w / (len(c.listOfControls) + 1)
	start := spaceInBetween
	for _, text := range c.listOfControls {
		centerOfText := len(text) / 2 // center the text
		window.Text(s, screentui.P(float64(start-centerOfText), float64(h-2)), text)
		start += spaceInBetween
	}
}

func (c *Controls) Events(ev *tcell.EventKey) {
	s := ev.Str()
	if len(s) == 1 && s[0] >= '1' && s[0] <= '3' {
		c.MenuResetList()

		page := int(s[0] - '0')
		switch page {
		case 1:
			c.MenuAddToList("[R] Remove")
			c.MenuAddToList("[J] Jump")
			// page 1
		case 2:
			c.MenuAddToList("Test 1")
			c.MenuAddToList("Test")
			// page 2
		case 3:
			// page 3
		}
	}
}
