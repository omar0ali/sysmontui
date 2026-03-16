package main

import (
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/screentui"
)

func main() {
	s, err := screentui.New(&screentui.ScreenOption{
		Ticker:  time.Second,
		ShowFPS: true,
	})

	if err != nil {
		panic(err)
	}

	s.Run(func(delta float64) {
		// update
	}, func() {
		// render
	}, func(e *tcell.EventKey) {
		if e.Key() == tcell.KeyCtrlQ {
			s.Exit()
		}
	})
}
