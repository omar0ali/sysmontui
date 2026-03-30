package cpustat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmon/pkg"
	"github.com/omar0ali/sysmontui/scenes/options"
	"github.com/omar0ali/sysmontui/scenes/perm/controls"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type CpuStat struct {
	mu       sync.RWMutex
	cpuUsage []float64
	avg      float64
	logsFunc controls.LogsControl
}

func Init(logsFunc controls.LogsControl, ctx context.Context, op options.Options) *CpuStat {

	cpustats := &CpuStat{
		logsFunc: logsFunc,
	}

	cpustats.logsFunc("Reading cpu stats...")

	go func(ctx context.Context, cs *CpuStat, op options.Options) {
		prev, _ := pkg.ReadCpuStat()
		ticker := time.NewTicker(time.Second * time.Duration(op.Interval))
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				curr, err := pkg.ReadCpuStat()
				if err != nil {
					panic(err)
				}
				delta := pkg.DeltaCPUStats(prev, curr)

				var usages []float64
				var sum float64

				for _, d := range delta {
					total := d.User + d.Nice + d.System + d.Idle + d.Iowait + d.Irq + d.SoftIrq
					idle := d.Idle + d.Iowait
					usage := float64(total-idle) / float64(total) * 100

					usages = append(usages, usage)
					sum += usage
				}

				// changing values - locking
				cs.mu.Lock()
				cs.cpuUsage = usages
				cs.avg = sum / float64(len(usages))
				cs.mu.Unlock()

				prev = curr // update prev for the next iteration
			}
		}
	}(ctx, cpustats, op)

	return cpustats
}

func (c *CpuStat) Update(d float64) {}
func (c *CpuStat) Render(s interfaces.ScreenControl) {
	s.Color(color.White)
	window.Text(s, screentui.P(1, 1), "Page: CPU Status") // title

	paddingBetweenText := 15
	startXPos := 32

	startY := 4
	for y, d := range c.cpuUsage {
		window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(y+startY)), paddingBetweenText, []string{
			fmt.Sprintf("CPU%d usage:", y),
			fmt.Sprintf("%.1f%%", d),
		})
	}
}
func (c *CpuStat) Events(ev tcell.Event) {}
