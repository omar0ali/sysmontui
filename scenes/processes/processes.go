package processes

import (
	"context"
	"fmt"
	"strings"
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

type process struct {
	Name       string
	PID        int
	CPUPercent float64
}

type ProcessesScene struct {
	mu sync.RWMutex

	Logs controls.LogsControl

	processes    []*process
	scrollWindow scrollWindow

	// order
	sortedBy sortBy
	desc     bool

	// actions
	search    *search
	kill      *kill
	searchFor string

	selectedProcess *process

	mController interfaces.MenuController
}

func Init(logsFunc controls.LogsControl, ctx context.Context, op options.Options) *ProcessesScene {
	processes := &ProcessesScene{
		Logs:         logsFunc,
		processes:    []*process{},
		sortedBy:     sortByName,
		scrollWindow: scrollWindow{},
	}

	// this used to lock access to to change current page when typing is required.
	if op.MenuController == nil {
		panic("Missing Controls")
	} else {
		processes.mController = op.MenuController
	}

	processes.Logs("Reading processes...")

	go func(ctx context.Context, processes *ProcessesScene, op options.Options) {
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

				// First run
				// initialize baseline
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
				processes.processes = processes.processes[:0] // reset list

				counter := 1
				for pid, p := range procs {
					prev, ok := prevProcCPU[pid]
					if !ok {
						// init baseline for new process
						prevProcCPU[pid] = p.Stat.UTime + p.Stat.STime
						continue
					}

					curr := p.Stat.UTime + p.Stat.STime
					procDelta := curr - prev

					cpuPercent := float64(procDelta) / float64(totalDelta) * 100

					// -- apply search
					name := strings.ToLower(p.Stat.Comm)
					search := strings.TrimSpace(strings.ToLower(processes.searchFor))

					if search != "" && !strings.Contains(name, search) {
						continue
					}

					// --
					processes.processes = append(processes.processes, &process{name, pid, cpuPercent})
					counter++
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

func (p *ProcessesScene) Update(d float64) {}

func (p *ProcessesScene) Render(s interfaces.ScreenControl) {
	if p.search != nil {
		p.search.Render(s)
	}
	if p.kill != nil {
		p.kill.Render(s)
	}

	p.scrollWindow.Render(p, s)

	s.Color(color.White)
	window.Text(s, screentui.P(1, 1), "Page: Running Processes") // title

	paddingBetweenText := 30
	startXPos := 32
	startYPos := 4
	window.Text(s,
		screentui.P(float64(startXPos), float64(startYPos)),
		fmt.Sprintf("Total Processes: %d | Displaying (%d - %d)", len(p.processes),
			p.scrollWindow.start, p.scrollWindow.end,
		),
	)

	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+1)), paddingBetweenText,
		[]string{"Name", "PID", "CPU"},
	)

	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+2)), paddingBetweenText,
		[]string{"-------------", "-------------", "------------"},
	)

	if len(p.processes) == 0 {
		status := "Loading..."
		if p.searchFor != "" {
			status = "No Processes"
		}
		window.Text(s,
			screentui.P(float64(startXPos), float64(startYPos+3)),
			status,
		)
		return
	}

	// sort processes by name default
	sortProcesses(p.sortedBy, p.desc, p.processes)

	// displaying list of processes
	currentProcessIndex := p.scrollWindow.currentIndex + startYPos + 3
	for y, proc := range p.processes[p.scrollWindow.start:p.scrollWindow.end] {
		var name string
		if currentProcessIndex == startYPos+y+3 {
			name = fmt.Sprintf("[K] %s", proc.Name)
			p.selectedProcess = proc
			s.Color(color.YellowGreen)
		} else {
			name = fmt.Sprintf("%s", proc.Name)
			s.Color(color.White)
		}

		window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+y+3)), paddingBetweenText, []string{
			name,
			fmt.Sprintf("%d", proc.PID),
			fmt.Sprintf("%.2f%%", proc.CPUPercent),
		})
	}
}

func (p *ProcessesScene) Events(ev tcell.Event) {
	// search processes events
	if p.search != nil {
		p.search.Events(p, ev)
		return
	}

	// kill a process events
	if p.kill != nil {
		p.kill.Events(p, ev)
		return
	}

	// scroll window events
	p.scrollWindow.Events(p, ev)

	// processes events
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Str() == "T" {
			p.sortedBy = nextSort(p.sortedBy)
			p.Logs(fmt.Sprintf("Sorted By %s", p.sortedBy.String()))
		} else if ev.Str() == "t" {
			p.desc = !p.desc
			order := "Ascending"
			if p.desc {
				order = "Descending"
			}
			p.Logs("Order: " + order)
		} else if ev.Str() == "/" {
			p.mController.Lock()
			p.search = &search{}
			p.Logs("Search ON")
		} else if ev.Str() == "K" {
			p.mController.Lock()
			p.kill = &kill{
				process: p.selectedProcess,
			}
			p.Logs("Process Selected: " + p.selectedProcess.Name + " - Waiting action...")
		}
	}
}
