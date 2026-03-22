package processes

import (
	"context"
	"fmt"
	"sort"
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

const (
	PADDING_FOOTER = 22
)

type scrollWindow struct {
	start, end, max int
}

type process struct {
	Name       string
	PID        int
	CPUPercent float64
}

type Processes struct {
	mu            sync.RWMutex
	LogsAddToList controls.LogsAddToList
	Processes     []*process
	scrollWindow  scrollWindow
	getScreenSize func() (int, int)
}

func Init(logsFunc controls.LogsAddToList, ctx context.Context, op options.Options) *Processes {
	processes := &Processes{
		LogsAddToList: logsFunc,
		Processes:     []*process{},
		scrollWindow:  scrollWindow{},
	}

	go func(ctx context.Context, processes *Processes, op options.Options) {
		ticker := time.NewTicker(time.Second * time.Duration(op.Interval))
		defer ticker.Stop()

		var prevCPU []*pkg.CPUStats
		prevProcCPU := map[int]uint64{}
		initialized := false

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				pids, err := pkg.GetPids()
				if err != nil {
					continue
				}

				procs := map[int]*pkg.Process{}
				for _, pid := range pids {
					p, err := pkg.NewProcess(pid)
					if err != nil {
						continue
					}
					procs[pid] = p
				}

				currCPU, err := pkg.ReadCpuStat()
				if err != nil {
					continue
				}

				// First run → just initialize baseline
				if !initialized {
					for pid, p := range procs {
						prevProcCPU[pid] = p.Stat.UTime + p.Stat.STime
					}
					prevCPU = currCPU
					initialized = true
					continue
				}

				pkg.RefreshProcesses(procs)

				delta := pkg.DeltaCPUStats(prevCPU, currCPU)

				totalDelta := uint64(0)
				for _, d := range delta {
					totalDelta += d.User + d.Nice + d.System + d.Idle +
						d.Iowait + d.Irq + d.SoftIrq
				}

				processes.mu.Lock()
				processes.Processes = processes.Processes[:0] // reset list

				for pid, p := range procs {
					prev, ok := prevProcCPU[pid]
					if !ok {
						continue
					}

					curr := p.Stat.UTime + p.Stat.STime
					procDelta := curr - prev

					cpuPercent := float64(procDelta) / float64(totalDelta) * 100

					processes.Processes = append(processes.Processes, &process{
						Name:       p.Stat.Comm,
						PID:        pid,
						CPUPercent: cpuPercent,
					})

					// update baseline
					prevProcCPU[pid] = curr
				}

				processes.mu.Unlock()

				prevCPU = currCPU
			}
		}
	}(ctx, processes, op)

	return processes
}

func (p *Processes) Update(d float64) {}

func (p *Processes) Render(s interfaces.ScreenControl) {
	// this will be used to get the size
	// when screen is resize triggered.
	if p.getScreenSize == nil {
		p.getScreenSize = s.Size
	}

	// important: used for the scrollable area where the processes are displayed
	// must set both (end, max)

	// check max must be set when
	if p.scrollWindow.max == 0 {
		_, h := p.getScreenSize()
		p.scrollWindow.max = h - PADDING_FOOTER
		// update end
		p.scrollWindow.end = p.scrollWindow.max
	}

	s.Color(color.White)
	window.Text(s, screentui.P(1, 1), "Page: Running Processes") // title

	paddingBetweenText := 40
	startXPos := 32
	startYPos := 4
	window.Text(s,
		screentui.P(float64(startXPos), float64(startYPos)),
		fmt.Sprintf("Total Processes: %d | Displaying (%d - %d)", len(p.Processes),
			p.scrollWindow.start, p.scrollWindow.end,
		),
	)

	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+1)), paddingBetweenText,
		[]string{"Name", "PID", "CPU"},
	)

	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+2)), paddingBetweenText,
		[]string{"-------------", "-------------", "------------"},
	)

	if len(p.Processes) == 0 {
		window.Text(s,
			screentui.P(float64(startXPos), float64(startYPos+3)),
			"Loading...",
		)
		return
	}

	// sort processes by name
	sort.Slice(p.Processes, func(i, j int) bool {
		return p.Processes[i].Name < p.Processes[j].Name
	})

	for y, proc := range p.Processes[p.scrollWindow.start:p.scrollWindow.end] {
		window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+y+3)), paddingBetweenText, []string{
			proc.Name,
			fmt.Sprintf("%d", proc.PID),
			fmt.Sprintf("%.2f%%", proc.CPUPercent),
		})
	}
}

func (p *Processes) Events(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventResize: // when screen resized
		_, h := p.getScreenSize()

		// get the latest hight of the screen -
		// used for the scrollable area in processes page
		p.scrollWindow.max = h - PADDING_FOOTER
		if p.scrollWindow.end == 0 {
			p.scrollWindow.end = p.scrollWindow.max
		}

		// update end
		p.scrollWindow.end = p.scrollWindow.max

	case *tcell.EventKey:
		if ev.Str() == "j" {
			if p.scrollWindow.end < len(p.Processes) {
				p.scrollWindow.start++
				p.scrollWindow.end++
			}
			p.LogsAddToList("Scrolling down\t[j]")
		} else if ev.Str() == "k" {
			if p.scrollWindow.start > 0 {
				p.scrollWindow.start--
				p.scrollWindow.end--
			}
			p.LogsAddToList("Scrolling up\t[k]")
		}
	}
}
