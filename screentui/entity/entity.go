package entity

import (
	"github.com/omar0ali/sysmontui/screentui/interfaces"
)

type Entity struct {
	entities     map[string][]interfaces.EntityActions
	permEntities []interfaces.EntityActions
	currentScene string
}

// all entities should have init func

func Init() *Entity {
	return &Entity{
		entities:     map[string][]interfaces.EntityActions{},
		permEntities: []interfaces.EntityActions{},
	}
}

func (e *Entity) SetScene(s string) {
	e.currentScene = s
}

func (e *Entity) GetScene() string {
	return e.currentScene
}

func (e *Entity) GetEntities(scene string) []interfaces.EntityActions {
	return e.entities[scene]
}

func GetAllEntitiesByType[T any](entities []interfaces.EntityActions) []T {
	var result []T
	for _, v := range entities {
		if val, ok := v.(T); ok {
			result = append(result, val)
		}
	}
	return result
}

func (e *Entity) GetPermEntities() []interfaces.EntityActions {
	return e.permEntities
}

func (e *Entity) AddPermEntity(entity ...interfaces.EntityActions) {
	e.permEntities = append(e.permEntities, entity...)
}

func (e *Entity) AddEntity(scene string, entity ...interfaces.EntityActions) {
	e.entities[scene] = append(e.entities[scene], entity...)
	e.currentScene = scene // should also set the last scene as default
}

func (e *Entity) ClearEntities(scene string) {
	clear(e.entities[scene])
}

func (e *Entity) ClearPermEntities() {
	e.permEntities = e.permEntities[:0]
}
