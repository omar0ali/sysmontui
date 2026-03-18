package cpuinfo

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmon/sysmon"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type CpuInfo struct {
	cpuinfo     *sysmon.CpuInfo
	cpuinfoName screentui.Point
	cpuinfoLC   screentui.Point
	cpuinfoPC   screentui.Point
}

func Init() *CpuInfo {
	info, err := sysmon.ReadCpuInfo()
	if err != nil {
		panic(err)
	}
	return &CpuInfo{
		cpuinfo:     info,
		cpuinfoName: screentui.P(0, 2),
		cpuinfoLC:   screentui.P(0, 3),
		cpuinfoPC:   screentui.P(0, 4),
	}
}

func (c *CpuInfo) Update(d float64) {}

func (c *CpuInfo) Render(s interfaces.ScreenControl) {

	s.Color(color.Blue)
	w, _ := s.Size()
	window.DrawBox(s, screentui.P((float64(w)/2)-25/2, 1), 25, 5)
	s.Color(color.White)
	window.TextCenter(s, c.cpuinfoName, fmt.Sprintf("Name: %s", c.cpuinfo.ModelName))
	window.TextCenter(s, c.cpuinfoLC, fmt.Sprintf("Logical CPU'S: %d", c.cpuinfo.LogicalCPUs))
	window.TextCenter(s, c.cpuinfoPC, fmt.Sprintf("Physical Cores: %d", c.cpuinfo.PhysicalCores))
}

func (c *CpuInfo) Events(ev *tcell.EventKey) {}
