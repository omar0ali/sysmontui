package main

import (
	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmontui/screen"
)

func main() {

	s, err := screen.New()
	if err != nil {
		panic(err)
	}

	s.Run(func() {
		// update
	}, func() {
		// render
	}, func(e tcell.Event) {
		if ev, ok := e.(*tcell.EventKey); ok {
			if ev.Key() == tcell.KeyCtrlQ {
				s.Exit()
			}
		}
	})
}
