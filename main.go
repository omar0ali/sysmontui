package main

import (
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/scenes/cpuinfo"
	"github.com/omar0ali/sysmontui/screentui"
)

func main() {
	s, err := screentui.New(&screentui.ScreenOption{
		Ticker: time.Second / 5,
	})

	if err != nil {
		panic(err)
	}

	// setup

	screentui.AddEntity("cpuinfo", cpuinfo.Init())

	// set scene

	screentui.CurrentScene = "cpuinfo"

	// run

	s.Run(func(delta float64) {
		// update
		for _, j := range screentui.GetEntities(screentui.CurrentScene) {
			j.Update(delta)
		}
	}, func() {
		// render
		for _, j := range screentui.GetEntities(screentui.CurrentScene) {
			j.Render(s)
		}
	}, func(e *tcell.EventKey) {
		// current
		if e.Key() == tcell.KeyCtrlQ {
			s.Exit()
		}
		// from entities
		for _, j := range screentui.GetEntities(screentui.CurrentScene) {
			j.Events(e)
		}

	})
}
