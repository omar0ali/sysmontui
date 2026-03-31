package interfaces

import (
	"github.com/gdamore/tcell/v3"
)

type EntityActions interface {
	Update(d float64)
	Render(ScreenControl)
	Events(ev tcell.Event)
}

type SceneControl interface {
	GetScene() string
	SetScene(s string) bool
}

type MenuController interface {
	Lock()
	Unlock()
}

func MustMenuController(imc MenuController) MenuController {
	if imc == nil {
		panic("Missing Controls")
	}
	return imc
}
