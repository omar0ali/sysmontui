package meminfo

import (
	"context"
	"fmt"
	"sync"

	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmon/pkg"
	"github.com/omar0ali/sysmontui/scenes"
	"github.com/omar0ali/sysmontui/scenes/options"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type MemInfo struct {
	meminfo *pkg.MemInfo
	mu      sync.RWMutex
	unit    pkg.Unit
	unitStr string
}

func Init(op options.Settings) *MemInfo {
	meminfo := &MemInfo{
		unit:    pkg.MB,
		unitStr: "MB",
	}

	scenes.Log("Read memory info...")

	// important to init info to avoid nil
	if info, err := pkg.ReadMemInfo(meminfo.unit); err == nil {
		meminfo.mu.Lock()
		meminfo.meminfo = info
		meminfo.mu.Unlock()
	}

	go func(ctx context.Context, mi *MemInfo, op options.Settings) {
		ticker := time.NewTicker(time.Second * time.Duration(op.Interval))
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				memInfo, err := pkg.ReadMemInfo(mi.unit)
				if err != nil {
					panic(err)
				}

				mi.mu.Lock()
				mi.meminfo = memInfo
				mi.mu.Unlock()
			}
		}
	}(op.Context, meminfo, op)

	return meminfo
}

func (m *MemInfo) Update(d float64) {}
func (m *MemInfo) Render(s interfaces.ScreenControl) {
	s.Color(color.White)
	window.Text(s, screentui.P(1, 1), "Page: Memory Info") // title

	// start writing the text line by line
	// first column
	paddingBetweenText := 15
	startXPos := 32
	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), 4), paddingBetweenText, []string{
		"Total RAM:",
		fmt.Sprintf("%d %s", m.meminfo.Total, m.unitStr),
	})
	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), 5), paddingBetweenText, []string{
		"Available RAM:",
		fmt.Sprintf("%d %s", m.meminfo.Available, m.unitStr),
	})
	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), 6), paddingBetweenText, []string{
		"Used RAM:",
		fmt.Sprintf("%d %s", m.meminfo.Total-m.meminfo.Available, m.unitStr),
	})
	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), 7), paddingBetweenText, []string{
		"Free RAM:",
		fmt.Sprintf("%d %s", m.meminfo.Free, m.unitStr),
	})

	// next column
	startXPos += startXPos
	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), 4), paddingBetweenText, []string{
		"Cached:",
		fmt.Sprintf("%d %s", m.meminfo.Cached, m.unitStr),
	})
	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), 5), paddingBetweenText, []string{
		"Buffers:",
		fmt.Sprintf("%d %s", m.meminfo.Buffers, m.unitStr),
	})
	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), 6), paddingBetweenText, []string{
		"SwapFree:",
		fmt.Sprintf("%d %s", m.meminfo.SwapFree, m.unitStr),
	})
	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), 7), paddingBetweenText, []string{
		"SwapTotal:",
		fmt.Sprintf("%d %s", m.meminfo.SwapTotal, m.unitStr),
	})
}

func (m *MemInfo) Events(ev tcell.Event) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		if e.Str() == "T" {
			if m.unit == pkg.MB {
				m.unit = pkg.GB
				m.unitStr = "GB"
				scenes.Log("Toggle to GB")
			} else {
				m.unit = pkg.MB
				m.unitStr = "MB"
				scenes.Log("Toggle to MB")
			}
		}
	}
}
