package screen

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v3"
)

var style tcell.Style

type Screen struct {
	screen tcell.Screen
	style  tcell.Style
	quit   chan struct{}
}

func New() (*Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}
	return &Screen{
		screen: s,
		style:  tcell.StyleDefault.Foreground(tcell.ColorDefault),
		quit:   make(chan struct{}),
	}, nil
}

func (s *Screen) Run(update func(), render func(), events func(ev tcell.Event)) {
	ticker := time.NewTicker(time.Second / 29) // 30 fps
	defer ticker.Stop()

	// input
	go func() {
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
			events(ev)
		}
	}()

	// display

	last := time.Now()
	var fps float64

	for {
		select {
		case <-ticker.C:
			// Calculate FPS
			now := time.Now()
			elapsed := now.Sub(last).Seconds()
			last = now
			if elapsed > 0 {
				fps = 1 / elapsed
			}

			// Clear the screen
			s.screen.Clear()

			// Update logic
			update()

			// Render your content
			render()

			// Draw FPS at top-left corner
			fpsStr := fmt.Sprintf("FPS: %.1f", fps)
			for i, r := range fpsStr {
				s.SetContent(i, 0, r)
			}

			// Show the buffer
			s.screen.Show()

		case <-s.quit:
			return
		}
	}
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
