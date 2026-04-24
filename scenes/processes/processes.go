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
	"github.com/omar0ali/sysmontui/scenes/perm/effects"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type process struct {
	Name       string
	PID        int
	CPUPercent float64
	MEMUsage   uint64
}

type ProcessesScene struct {
	mu sync.RWMutex

	Logs controls.LogsControl

	processes    []*process
	scrollWindow scrollWindow

	// order
	sortedBy sortBy

	// actions
	searchState *SearchState
	search      *search
	kill        *kill

	selectedProcess *process

	mController interfaces.MenuController
}

func Init(logsFunc controls.LogsControl, ctx context.Context, op options.Options) *ProcessesScene {

	processes := &ProcessesScene{
		Logs:      logsFunc,
		processes: []*process{},
		sortedBy:  sortByName,

		scrollWindow: scrollWindow{},

		searchState: &SearchState{
			isLoading: true,
		},
	}

	// this used to lock access to to change current page when typing is required.
	processes.mController = interfaces.MustMenuController(op.MenuController)

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

					memUsage := p.Status.VmRSS / 1024 // to get MB

					// -- apply search
					name := strings.ToLower(p.Stat.Comm)
					search := strings.TrimSpace(strings.ToLower(processes.searchState.strSearch))

					if search != "" && !strings.Contains(name, search) {
						continue
					}

					// --
					processes.processes = append(processes.processes, &process{name, pid, cpuPercent, memUsage})
					counter++
					// update baseline
					prevProcCPU[pid] = curr
				}
				processes.mu.Unlock()

				prevCPU = currCPU
			}
			// disable loading when done. // status will show 'No Processes' if the list is empty
			processes.searchState.isLoading = false
		}
	}(ctx, processes, op)

	return processes
}

func (p *ProcessesScene) Update(d float64) {
	// update status | only if its true
	if p.searchState.isLoading {
		p.processes = p.processes[:0]
	}
}

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

	paddingBetweenText := 25
	startXPos := 32
	startYPos := 4
	window.Text(s,
		screentui.P(float64(startXPos), float64(startYPos)),
		fmt.Sprintf("Total Processes: %d | Displaying (%d - %d)", len(p.processes),
			p.scrollWindow.start, p.scrollWindow.end,
		),
	)

	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+1)), paddingBetweenText,
		[]string{"Name", "PID", "CPU", "MemB"},
	)

	window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+2)), paddingBetweenText,
		[]string{"-------------", "-------------", "------------", "------------"},
	)

	// displaying status on screen
	if len(p.processes) == 0 {
		if p.searchState.isLoading {
			p.searchState.status = "Loading " + string(effects.Spinner(7, 150))
		} else {
			p.searchState.status = "No Processes"
		}
		window.Text(s,
			screentui.P(float64(startXPos), float64(startYPos+3)),
			p.searchState.status,
		)
		return
	}

	// sort processes by name default
	sortProcesses(p.sortedBy, p.processes)

	startProcessIndex := p.scrollWindow.currentIndex + startYPos + 3

	// this will avoid index out of bound, triggered when hiding logs
	lastIndex := min(p.scrollWindow.end, len(p.processes))

	// displaying list of processes
	for y, proc := range p.processes[p.scrollWindow.start:lastIndex] {
		var name string
		if startProcessIndex == startYPos+y+3 {
			name = fmt.Sprintf("[K] %s", proc.Name)
			p.selectedProcess = proc
			s.Color(color.YellowGreen)
		} else {
			name = fmt.Sprintf("%s", proc.Name)
			s.DefaultColor()
		}

		window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+y+3)), paddingBetweenText, []string{
			name,
			fmt.Sprintf("%d", proc.PID),
			fmt.Sprintf("%.2f%%", proc.CPUPercent),
			fmt.Sprintf("%d MB", proc.MEMUsage),
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
			var order string
			if toggleDesc() {
				order = "Descending"
			} else {
				order = "Ascending"

			}
			p.Logs("Order: " + order)
		} else if ev.Str() == "/" {
			p.mController.Lock()
			p.search = &search{}
			p.Logs("Search ON")
		} else if ev.Str() == "K" {
			if len(p.processes) < 1 {
				p.Logs("Cannot perform this action now.")
				return
			}
			p.mController.Lock()
			p.kill = &kill{
				process: p.selectedProcess,
			}
			p.Logs("Process Selected: " + p.selectedProcess.Name + " - Waiting action...")
		} else if ev.Str() == "q" {
			// clear / reset current invoked search
			if p.searchState.strSearch != "" {
				p.searchState.isLoading = true
				p.scrollWindow.currentIndex = 0
				p.Logs("Search cleared...")
				p.searchState.strSearch = ""
			}
			p.mController.Unlock()
		}
	}
}
