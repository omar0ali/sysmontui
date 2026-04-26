package options

import (
	"context"

	"github.com/omar0ali/sysmontui/scenes/processes/parts/logsui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
)

type Settings struct {
	Interval       int
	MenuController interfaces.MenuController
	Context        context.Context
	LogsUIControl  logsui.LogsView
}
