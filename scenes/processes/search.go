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

func (s *search) Events(p *ProcessesScene, ev tcell.Event) {
	closeSearchWith := func(search string) {
		p.scrollWindow.currentIndex = 0
		p.searchFor = search
		p.search = nil
		p.mController.Unlock()
	}
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Str() {
		case "/":
			p.Logs("Search Canceled / Reset")
			closeSearchWith("")
			return
		}
		switch ev.Key() {
		case tcell.KeyEnter:
			p.Logs("Searching for: " + p.search.text)
			closeSearchWith(p.search.text)
			return
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if len(s.text) > 0 {
				s.text = s.text[:len(s.text)-1]
			}
		default:
			s.text += ev.Str()
		}
	}
}
