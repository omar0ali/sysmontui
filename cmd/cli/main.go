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

	entity := entity.Init() // create a new entity (where a collection of entities are collected)

	// add entities you want to display (each scene will have their own list of entities)
	// each entity must contain (Init, Update, Render, Events) functions / actions

	entity.AddPermEntity(
		controls.Init(entity),
		cpuinfo.Init(),
	) // each entity must have init

	entity.AddEntity("0", home.Init())    // non perm scene
	entity.AddEntity("1", meminfo.Init()) // non perm scene

	// set scene (optional) Always last entity added will be set as the current scene.

	entity.SetScene("0") // set current scene to be displayed / rendered

	s.Run(entity) // run
}
