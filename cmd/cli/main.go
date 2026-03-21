package main

import (
	"time"

	"github.com/omar0ali/sysmontui/scenes/home"
	"github.com/omar0ali/sysmontui/scenes/meminfo"
	"github.com/omar0ali/sysmontui/scenes/perm/controls"
	"github.com/omar0ali/sysmontui/scenes/perm/cpuinfo"
	"github.com/omar0ali/sysmontui/screentui/entity"
	"github.com/omar0ali/sysmontui/screentui/window"
)

func main() {
	// setup
	s, err := window.New(&window.ScreenOption{
		Ticker: time.Second / 5,
	})

	if err != nil {
		panic(err)
	}

	entities := entity.Init() // create a new entity (where a collection of entities are collected)

	// add entities to display (each scene will have their own list of entities)
	// each entity must contain (Init, Update, Render, Events) functions / actions

	// permanent entities (always shown on all scenes)
	entities.AddPermEntity(
		controls.Init(entities),
		cpuinfo.Init(),
	)

	// logs func can be used anywhere, when its required
	logsFunc := entity.GetAllEntitiesByType[*controls.Controls](
		entities.GetPermEntities(),
	)[0].LogsAddToList

	entities.AddEntity("0", home.Init())            // non perm scene
	entities.AddEntity("1", meminfo.Init(logsFunc)) // non perm scene

	entities.SetScene("0") // set current scene to be displayed / rendered

	s.Run(entities) // run
}
