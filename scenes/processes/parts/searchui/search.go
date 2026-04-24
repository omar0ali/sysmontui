package searchui

import (
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/scenes"
	"github.com/omar0ali/sysmontui/scenes/processes/parts"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type Search struct {
	Text string
}

func (s *Search) Render(sc interfaces.ScreenControl) {
	txt := s.Text
	if (time.Now().UnixMilli()/500)%2 == 1 {
		txt += string(tcell.RuneBlock)
	}
	window.Text(sc, screentui.P(32, 0), "Filter")
	window.Text(sc, screentui.P(32, 1), txt)
	window.Text(sc, screentui.P(32, 2), "[Enter] Search - [/, q] Cancel | Reset")
}

func (s *Search) Events(c parts.SearchController, ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Str() {
		case "/", "q":
			c.ClosingSearchWith("") // clear search
			return
		}
		switch ev.Key() {
		case tcell.KeyEnter:
			scenes.Log("Searching for: " + s.Text)
			c.SetLoading()
			c.ClosingSearchWith(s.Text)
			return
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if len(s.Text) > 0 {
				s.Text = s.Text[:len(s.Text)-1]
			}
		case tcell.KeyCtrlW:
			s.Text = ""
		default:
			s.Text += ev.Str()
		}
	}
}
