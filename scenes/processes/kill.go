package processes

import (
	"fmt"
	"os"
	"syscall"

	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type kill struct {
	process *process
}

func (k *kill) Render(sc interfaces.ScreenControl) {
	window.Text(sc, screentui.P(33, 1), "Name: "+k.process.Name+" PID: "+fmt.Sprint(k.process.PID))
	window.Text(sc, screentui.P(33, 2), "[Enter] Kill - [/] Cancel")
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
