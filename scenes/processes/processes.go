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

var (
	start = 0
	end   = 0
)

type process struct {
	Name       string
	PID        int
	CPUPercent float64
}

type Processes struct {
	mu            sync.RWMutex
	LogsAddToList controls.LogsAddToList
	Processes     []*process
}

func Init(logsFunc controls.LogsAddToList, ctx context.Context, op options.Options) *Processes {
	processes := &Processes{
		LogsAddToList: logsFunc,
		Processes:     []*process{},
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
	if end == 0 {
		_, h := s.Size()
		end = (h / 2) - 6 // limit
	}

	s.Color(color.White)
	window.Text(s, screentui.P(1, 1), "Page: Running Processes") // title

	paddingBetweenText := 40
	startXPos := 32
	startYPos := 4
	window.Text(s,
		screentui.P(float64(startXPos), float64(startYPos)),
		fmt.Sprintf("Total Processes: %d | Displaying (%d - %d)", len(p.Processes), start, end),
	)

	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+1)), paddingBetweenText,
		[]string{"Name", "PID", "CPU"},
	)

	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+2)), paddingBetweenText,
		[]string{"-------------", "-------------", "------------"},
	)

	if len(p.Processes) == 0 {
		window.Text(s, screentui.P(float64(startXPos), float64(startYPos+3)), "Loading...")
		return
	}

	// sort processes by name

	sort.Slice(p.Processes, func(i, j int) bool {
		return p.Processes[i].Name < p.Processes[j].Name
	})

	for y, proc := range p.Processes[start:end] {
		window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+y+3)), paddingBetweenText, []string{
			proc.Name,
			fmt.Sprintf("%d", proc.PID),
			fmt.Sprintf("%.2f%%", proc.CPUPercent),
		})
	}
}

func (p *Processes) Events(ev *tcell.EventKey) {
	if ev.Str() == "j" {
		if end < len(p.Processes) {
			start++
			end++
		}
		p.LogsAddToList("Scrolling down\t[j]")
	} else if ev.Str() == "k" {
		if start > 0 {
			start--
			end--
		}
		p.LogsAddToList("Scrolling up\t[k]")
	}
}
