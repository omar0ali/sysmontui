package main

import (
	"context"
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/omar0ali/sysmontui/scenes"
	"github.com/omar0ali/sysmontui/scenes/cpustat"
	"github.com/omar0ali/sysmontui/scenes/meminfo"
	"github.com/omar0ali/sysmontui/scenes/perm/controls"
	"github.com/omar0ali/sysmontui/scenes/perm/cpuinfo"
	"github.com/omar0ali/sysmontui/scenes/processes/pScene"
	"github.com/omar0ali/sysmontui/scenes/processes/parts/logsui"
	"github.com/omar0ali/sysmontui/scenes/settings"
	"github.com/omar0ali/sysmontui/screentui/entity"
	"github.com/omar0ali/sysmontui/screentui/window"
)

func main() {
	// setup
	s, err := window.New(&window.ScreenOption{
		Ticker: time.Second / 10,
		// ShowFPS: true,
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

	logsuiView := logsui.New()

	// scenes setup
	controlScene := controls.Init(entities, logsuiView)

	// logs func can be used anywhere, when its required
	scenes.SetLogger(controlScene.LogsAddToList)

	// permanent entities (always shown on all scenes)
	entities.AddPermEntity(
		controlScene,
		cpuinfo.Init(),
	)

	// setting options
	ops := settings.New(
		2,
		controlScene, // control
		ctx,          // context
		logsuiView,   // logsUIView Control
	)

	// add entities to display (each scene will have their own list of entities)
	// each entity must contain (Init, Update, Render, Events) functions / actions

	// non perm scenes
	entities.AddEntity("1", meminfo.Init(ops))
	entities.AddEntity("2", cpustat.Init(ops))
	entities.AddEntity("3", pScene.Init(ops))

	entities.SetScene("1") // set current scene to be displayed / rendered

	s.Run(entities, cancel) // run
}
