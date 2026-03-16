package main

import (
	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/screentui"
)

func main() {

	s, err := screen.New()
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
