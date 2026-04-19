package processes

import (
	"fmt"
	"os"
	"syscall"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmontui/scenes/perm/effects"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type kill struct {
	process *process
}

func (k *kill) Render(sc interfaces.ScreenControl) {
	sc.Color(color.Red)
	window.Text(sc, screentui.P(33, 0), string(effects.Spinner(9, 500)))
	sc.DefaultColor()
	window.Text(sc, screentui.P(35, 0), "Are you sure?")
	window.Text(sc, screentui.P(33, 1), "Name: "+k.process.Name+" PID: "+fmt.Sprint(k.process.PID))
	window.Text(sc, screentui.P(33, 2), "[Enter] Kill - [/, q] Cancel")
}

func (k *kill) Events(p *ProcessesScene, ev tcell.Event) {
	close := func() {
		p.kill = nil
		p.mController.Unlock()
		p.scrollWindow.currentIndex = 0
	}
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Str() {
		case "/", "q":
			p.Logs("Canceled")
			close()
			return
		}
		switch ev.Key() {
		case tcell.KeyEnter:
			isLoading = true
			p.Logs("SIGTERM: " + p.selectedProcess.Name)
			err := p.kill.SIGTERM()
			if err != nil {
				p.Logs(err.Error())
			}
			close()
			return
		}
	}
}

func (k *kill) SIGTERM() error {
	if k.process == nil {
		return fmt.Errorf("Error: pid not set.")
	}
	// find process
	process, err := os.FindProcess(k.process.PID)
	if err != nil {
		return fmt.Errorf("can't find the process: %s", err)
	}
	// Send SIGTERM
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("Error sending signal: %s", err)
	}
	return nil
}
