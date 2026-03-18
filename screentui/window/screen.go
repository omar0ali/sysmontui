package window

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmontui/screentui/entity"
)

type screen struct {
	screen       tcell.Screen
	style        tcell.Style
	quit         chan struct{}
	screenOption *ScreenOption
}

func New(so *ScreenOption) (*screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	return &screen{
		screen:       s,
		style:        tcell.StyleDefault.Foreground(tcell.ColorDefault),
		quit:         make(chan struct{}),
		screenOption: so,
	}, nil
}

func (s *screen) keyEvents(events func(ev *tcell.EventKey)) {
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

func (s *screen) refreshScreen(update func(delta float64), render func()) {
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

func (s *screen) run(
	update func(elapsed float64),
	render func(),
	events func(ev *tcell.EventKey),
) {
	go s.keyEvents(events)          // read input on a go routine
	s.refreshScreen(update, render) // show screen and refresh
}

func (s *screen) Run(entity *entity.Entity) {
	s.run(func(delta float64) {
		// update
		for _, j := range entity.GetEntities(entity.GetScene()) {
			j.Update(delta)
		}
	}, func() {
		// render
		for _, j := range entity.GetEntities(entity.GetScene()) {
			j.Render(s)
		}
	}, func(e *tcell.EventKey) {
		// current
		if e.Key() == tcell.KeyCtrlQ {
			s.Exit()
		}
		// from entities
		for _, j := range entity.GetEntities(entity.GetScene()) {
			j.Events(e)
		}

	})
}

func (s *screen) Color(styleColor color.Color) {
	s.style = tcell.StyleDefault.Foreground(styleColor).Background(tcell.ColorReset)
}

func (s *screen) SetContent(x, y int, r rune) {
	s.screen.SetContent(x, y, r, nil, s.style)
}

func (s *screen) Size() (width int, height int) {
	return s.screen.Size()
}
func (s *screen) Exit() {
	s.screen.Fini()
	close(s.quit)
}
