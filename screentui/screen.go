package screentui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v3"
)

var style tcell.Style

type ScreenOption struct {
	Ticker  time.Duration
	ShowFPS bool
}

type Screen struct {
	screen       tcell.Screen
	style        tcell.Style
	quit         chan struct{}
	screenOption *ScreenOption
}

func New(so *ScreenOption) (*Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	return &Screen{
		screen:       s,
		style:        tcell.StyleDefault.Foreground(tcell.ColorDefault),
		quit:         make(chan struct{}),
		screenOption: so,
	}, nil
}

func (s *Screen) keyEvents(events func(ev *tcell.EventKey)) {
	for {
		ev := <-s.screen.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.screen.Sync()
		case *tcell.EventKey:

			switch ev.Key() {
			case tcell.KeyEscape:
				s.Exit()
			}
		}
		if ev, ok := ev.(*tcell.EventKey); ok {
			events(ev)
		}
	}
}

func (s *Screen) refreshScreen(update func(delta float64), render func()) {
	ticker := time.NewTicker(s.screenOption.Ticker) // 30 fps
	defer ticker.Stop()
	last := time.Now()
	for {
		select {
		case <-ticker.C:
			// get elapsed time
			now := time.Now()
			elapsed := now.Sub(last).Seconds()
			last = now

			update(elapsed)  // update logic
			s.screen.Clear() // clear screen
			render()         // render screen

			if s.screenOption.ShowFPS {
				fpsStr := fmt.Sprintf("FPS: %.1f", 1/elapsed)
				for i, r := range fpsStr {
					s.SetContent(i, 0, r)
				}
			}

			s.screen.Show() // show what has been rendered on screen
		case <-s.quit: // quit signal
			return
		}
	}

}

func (s *Screen) Run(
	update func(elapsed float64),
	render func(),
	events func(ev *tcell.EventKey),
) {
	go s.keyEvents(events)          // read input on a go routine
	s.refreshScreen(update, render) // show screen and refresh
}

func (s *Screen) Color(styleColor tcell.Color) {
	style = tcell.StyleDefault.Foreground(styleColor)
}

func (s *Screen) SetContent(x, y int, r rune) {
	s.screen.SetContent(x, y, r, nil, s.style)
}

func (s *Screen) Size() (width int, height int) {
	return s.screen.Size()
}
func (s *Screen) Exit() {
	s.screen.Fini()
	close(s.quit)
}
