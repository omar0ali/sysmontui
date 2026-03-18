package main

import (
	"time"

	"github.com/omar0ali/sysmontui/scenes/cpuinfo"
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

	entity.AddEntity("cpuinfo", cpuinfo.Init()) // each entity must have init

	// set scene (optional) Always last entity added will be set as the current scene.

	entity.SetScene("cpuinfo") // set current scene to be displayed / rendered

	// run

	s.Run(entity)
}
