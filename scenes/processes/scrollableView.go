package processes

import (
	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
)

const (
	PADDING_FOOTER = 22
)

type scrollWindow struct {
	start, end, max, currentIndex int
}

func (sc *scrollWindow) Events(p *ProcessesScene, ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Str() == "j" {
			if p.scrollWindow.currentIndex < p.scrollWindow.max-1 &&
				p.scrollWindow.start+p.scrollWindow.currentIndex < len(p.processes)-1 {
				p.scrollWindow.currentIndex++
				return
			}

			// otherwise scroll window
			if p.scrollWindow.end < len(p.processes)-1 {
				p.scrollWindow.start++
				p.scrollWindow.end++
			}
		} else if ev.Str() == "k" {
			if p.scrollWindow.currentIndex > 0 {
				p.scrollWindow.currentIndex--
				return
			}

			// otherwise scroll window up
			if p.scrollWindow.start > 0 {
				p.scrollWindow.start--
				p.scrollWindow.end--
			}
		}
	}
}

// important: used for the scrollable area where the processes are displayed
// must set both (end, max)
func (sc *scrollWindow) Render(p *ProcessesScene, s interfaces.ScreenControl) {
	_, h := s.Size()
	// get the latest hight of the screen -
	// used for the scrollable area in processes page
	sc.max = h - PADDING_FOOTER

	// this fix a problem when .start is at a position that > than what it
	// its displayed. When searching, it might crash.
	if sc.start > len(p.processes) {
		sc.start = 0
	}

	sc.end = sc.start + sc.max
	if len(p.processes) < sc.max {
		sc.end = len(p.processes)
	}
}
