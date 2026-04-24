package main

import (
	"context"
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/omar0ali/sysmontui/scenes/cpustat"
	"github.com/omar0ali/sysmontui/scenes/meminfo"
	"github.com/omar0ali/sysmontui/scenes/options"
	"github.com/omar0ali/sysmontui/scenes/perm/controls"
	"github.com/omar0ali/sysmontui/scenes/perm/cpuinfo"
	"github.com/omar0ali/sysmontui/scenes/processes"
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

	defer func() {
		if r := recover(); r != nil {
			s.Exit()
			log.Printf("Fatal error: %v\nStack trace:\n%s", r, debug.Stack())
		}
	}()

	// check if running linux
	if runtime.GOOS != "linux" {
		panic("Not running on linux")
	}

	entities := entity.Init() // create a new entity (where a collection of entities are collected)

	// this used to close all the go routines that are running
	ctx, cancel := context.WithCancel(context.Background())

	// add entities to display (each scene will have their own list of entities)
	// each entity must contain (Init, Update, Render, Events) functions / actions

	// permanent entities (always shown on all scenes)
	entities.AddPermEntity(
		controls.Init(entities),
		cpuinfo.Init(),
	)

	// control a Permanent Scene Menu
	control := entity.GetAllEntitiesByType[*controls.Controls](
		entities.GetPermEntities(),
	)[0]

	// logs func can be used anywhere, when its required
	logsFunc := control.LogsAddToList

	// non perm scenes
	// entities.AddEntity("0", home.Init()) // removed
	entities.AddEntity("1", meminfo.Init(logsFunc, ctx, options.Options{Interval: 2}))
	entities.AddEntity("2", cpustat.Init(logsFunc, ctx, options.Options{Interval: 2}))
	entities.AddEntity("3", processes.Init(logsFunc, ctx, options.Options{
		Interval:       3,
		MenuController: control,
	}))

	entities.SetScene("1") // set current scene to be displayed / rendered

	s.Run(entities, cancel) // run
}
