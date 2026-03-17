package screentui

import "github.com/gdamore/tcell/v3"

var (
	entities     map[string][]Entity
	CurrentScene string
)

type Point struct {
	X, Y float64
}

// all entities should have init func

type Entity interface {
	Update(d float64)
	Render(ScreenControl)
	Events(ev *tcell.EventKey)
}

func GetEntities(scene string) []Entity {
	return entities[scene]
}

func AddEntity(scene string, entity Entity) {
	if entities == nil {
		entities = map[string][]Entity{}
	}
	entities[scene] = append(entities[scene], entity)
}

func ClearEntities(scene string) {
	clear(entities[scene])
}

func P(X, Y float64) Point {
	return Point{X, Y}
}
