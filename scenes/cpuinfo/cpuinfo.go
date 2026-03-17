package cpuinfo

import (
	"fmt"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmon/sysmon"
	"github.com/omar0ali/sysmontui/screentui"
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

func (c *CpuInfo) Render(s screentui.ScreenControl) {

	s.Color(color.Blue)
	w, _ := s.Size()
	screentui.DrawBox(s, screentui.P((float64(w)/2)-25/2, 1), 25, 5)
	s.Color(color.White)
	screentui.TextCenter(s, c.cpuinfoName, fmt.Sprintf("Name: %s", c.cpuinfo.ModelName))
	screentui.TextCenter(s, c.cpuinfoLC, fmt.Sprintf("Logical CPU'S: %d", c.cpuinfo.LogicalCPUs))
	screentui.TextCenter(s, c.cpuinfoPC, fmt.Sprintf("Physical Cores: %d", c.cpuinfo.PhysicalCores))
}

func (c *CpuInfo) Events(ev *tcell.EventKey) {}
