package pScene

import (
	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/scenes/processes/parts/logsui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
)

type ScrollWindow struct {
	Start, End, Max, CurrentIndex int
}

func (sc *ScrollWindow) Events(p *ProcessesScene, ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Str() == "j" {
			// normal behavior when scrolling down
			if p.ScrollWindow.CurrentIndex < p.ScrollWindow.Max-1 &&
				p.ScrollWindow.Start+p.ScrollWindow.CurrentIndex < len(p.processes)-1 {
				p.ScrollWindow.CurrentIndex++
				return
			}

			// otherwise scroll window
			if p.ScrollWindow.End < len(p.processes) {
				p.ScrollWindow.Start++
				p.ScrollWindow.End++
			}
		} else if ev.Str() == "k" {
			if p.ScrollWindow.CurrentIndex > 0 {
				p.ScrollWindow.CurrentIndex--
				return
			}

			// otherwise scroll window up
			if p.ScrollWindow.Start > 0 {
				p.ScrollWindow.Start--
				p.ScrollWindow.End--
			}
		}
	}
}

// important: used for the scrollable area where the processes are displayed
// must set both (end, max)
func (sc *ScrollWindow) Render(p *ProcessesScene, s interfaces.ScreenControl) {
	// check current position of the highlighted process
	// this fix the position when user shows the logs back
	// it should fix the position of the selected item
	if p.ScrollWindow.CurrentIndex > p.ScrollWindow.Max {
		p.ScrollWindow.CurrentIndex = p.ScrollWindow.Max - 1
	}

	_, h := s.Size()

	var paddingFooter int
	if logsui.ShowLogs {
		paddingFooter = logsui.LOGS_PADDING
	} else {
		paddingFooter = logsui.LOGS_NO_PADDING
	}

	// get the latest hight of the screen -
	// used for the scrollable area in processes page
	sc.Max = h - paddingFooter

	// this fix a problem when .start is at a position that > than what it
	// its displayed. When searching, it might crash.
	if sc.Start > len(p.processes) {
		sc.Start = 0
	}

	sc.End = sc.Start + sc.Max
	if len(p.processes) < sc.Max {
		sc.End = len(p.processes)
	}
}
