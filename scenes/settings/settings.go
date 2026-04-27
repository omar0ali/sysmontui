package settings

import (
	"context"

	"github.com/omar0ali/sysmontui/scenes/processes/parts/logsui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
)

func New(interval int, mCtrl interfaces.MenuController, ctx context.Context, logsUI logsui.LogsView) Settings {
	return Settings{
		Interval:       interval,
		Context:        ctx,
		LogsUIControl:  logsUI,
		MenuController: mCtrl,
	}
}

type Settings struct {
	Interval       int
	MenuController interfaces.MenuController
	Context        context.Context
	LogsUIControl  logsui.LogsView
}
