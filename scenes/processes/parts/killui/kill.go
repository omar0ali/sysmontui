package killui

import (
	"fmt"
	"os"
	"syscall"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmontui/scenes"
	"github.com/omar0ali/sysmontui/scenes/perm/effects"
	"github.com/omar0ali/sysmontui/scenes/processes"
	"github.com/omar0ali/sysmontui/scenes/processes/parts"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type Kill struct {
	Process *processes.Process
}

func (k *Kill) Render(sc interfaces.ScreenControl) {
	sc.Color(color.Red)
	window.Text(sc, screentui.P(32, 0), string(effects.Spinner(9, 500)))
	sc.DefaultColor()
	window.Text(sc, screentui.P(34, 0), "Are you sure?")
	window.Text(sc, screentui.P(32, 1), "Name: "+k.Process.Name+" PID: "+fmt.Sprint(k.Process.PID))
	window.Text(sc, screentui.P(32, 2), "[Enter] Kill - [/, q] Cancel")
}

func (k *Kill) Events(c parts.KillController, ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Str() {
		case "/", "q":
			scenes.Log("Canceled")
			c.CloseingKillUI()
			return
		}
		switch ev.Key() {
		case tcell.KeyEnter:
			c.SetLoading()
			scenes.Log("SIGTERM: " + k.Process.Name)
			err := k.SIGTERM()
			if err != nil {
				scenes.Log(err.Error())
				c.CloseingKillUI()
			}
			c.CloseingKillUI()
			return
		}
	}
}

func (k *Kill) SIGTERM() error {
	if k.Process == nil {
		return fmt.Errorf("Error: pid not set.")
	}
	// find process
	process, err := os.FindProcess(k.Process.PID)
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
