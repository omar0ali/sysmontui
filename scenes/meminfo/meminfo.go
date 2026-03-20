package meminfo

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/omar0ali/sysmon/pkg"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type MemInfo struct {
	meminfo *pkg.MemInfo
	unit    pkg.Unit
	unitStr string
}

func Init() *MemInfo {
	meminfo, err := pkg.ReadMemInfo(pkg.MB)
	if err != nil {
		panic(err)
	}
	return &MemInfo{
		meminfo: meminfo,
		unit:    pkg.MB,
		unitStr: "MB",
	}
}

func (m *MemInfo) Update(d float64) {}
func (m *MemInfo) Render(s interfaces.ScreenControl) {
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

func (m *MemInfo) Events(ev *tcell.EventKey) {
	if ev.Str() == "T" {
		if m.unit == pkg.MB {
			m.unit = pkg.GB
			m.unitStr = "GB"
		} else {
			m.unit = pkg.MB
			m.unitStr = "MB"
		}
	}
	if ev.Str() == "U" || ev.Str() == "T" { // update current stats
		meminfo, err := pkg.ReadMemInfo(m.unit)
		if err != nil {
			panic(err)
		}
		m.meminfo = meminfo
	}
}
