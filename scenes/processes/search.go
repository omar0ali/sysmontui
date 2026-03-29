package processes

import (
	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type search struct {
	text string
}

func (s *search) Render(sc interfaces.ScreenControl) {
	window.Text(sc, screentui.P(33, 1), "Search: "+s.text)
	window.Text(sc, screentui.P(33, 2), "[Enter] Search - [/] Cancel | Reset")
}

func (s *search) Events(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if len(s.text) > 0 {
				s.text = s.text[:len(s.text)-1]
			}
		default:
			s.text += ev.Str()
		}
	}
}
