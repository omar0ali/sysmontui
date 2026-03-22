package cpuinfo

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmon/pkg"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type CpuInfo struct {
	cpuinfo *pkg.CpuInfo
}

func Init() *CpuInfo {
	info, err := pkg.ReadCpuInfo()
	if err != nil {
		panic(err)
	}
	return &CpuInfo{
		cpuinfo: info,
	}
}

func (c *CpuInfo) Update(d float64) {}

func (c *CpuInfo) Render(s interfaces.ScreenControl) {
	s.Color(color.White)
	w, h := s.Size()
	window.LineHorizontal(s, w, h-3, tcell.RuneHLine)

	details := []string{
		c.cpuinfo.ModelName,
		fmt.Sprintf("Logical CPUs: %d", c.cpuinfo.LogicalCPUs),
		fmt.Sprintf("Physical Cores: %d", c.cpuinfo.PhysicalCores),
	}

	spaceInBetween := w / (len(details) + 1)

	start := spaceInBetween
	for _, text := range details {
		centerOfText := len(text) / 2 // center the text
		window.Text(s, screentui.P(float64(start-centerOfText), float64(h-2)), text)
		start += spaceInBetween
	}
}

func (c *CpuInfo) Events(ev tcell.Event) {}
