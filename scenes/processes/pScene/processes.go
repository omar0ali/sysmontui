package pScene

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/omar0ali/sysmon/pkg"
	"github.com/omar0ali/sysmontui/scenes"
	"github.com/omar0ali/sysmontui/scenes/perm/effects"
	"github.com/omar0ali/sysmontui/scenes/processes"
	"github.com/omar0ali/sysmontui/scenes/processes/parts/killui"
	"github.com/omar0ali/sysmontui/scenes/processes/parts/searchui"
	"github.com/omar0ali/sysmontui/scenes/settings"
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
	"github.com/omar0ali/sysmontui/screentui/window"
)

type ProcessesScene struct {
	mu sync.RWMutex

	processes    []*processes.Process
	ScrollWindow ScrollWindow

	// order
	sortedBy sortBy

	// actions ui
	SearchState *searchui.State
	Search      *searchui.Search
	kill        *killui.Kill

	selectedProcess *processes.Process

	MenuController interfaces.MenuController
}

func Init(op settings.Settings) *ProcessesScene {

	pScene := &ProcessesScene{
		processes: []*processes.Process{},
		sortedBy:  sortByName,

		ScrollWindow: ScrollWindow{
			LogsUiView: op.LogsUIControl,
		},

		SearchState: &searchui.State{
			IsLoading: true,
		},
	}

	// this used to lock access to to change current page when typing is required.
	pScene.MenuController = interfaces.MustMenuController(op.MenuController)

	scenes.Log("Reading processes...")

	go func(ctx context.Context, pScene *ProcessesScene, op settings.Settings) {
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

				pScene.mu.Lock()
				pScene.processes = pScene.processes[:0] // reset list

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
					search := strings.TrimSpace(strings.ToLower(pScene.SearchState.StrSearch))

					if search != "" && !strings.Contains(name, search) {
						continue
					}

					// --
					pScene.processes = append(pScene.processes, &processes.Process{
						Name:       name,
						PID:        pid,
						CPUPercent: cpuPercent,
						MEMUsage:   memUsage,
					})

					counter++
					// update baseline
					prevProcCPU[pid] = curr
				}
				pScene.mu.Unlock()

				prevCPU = currCPU
			}
			// disable loading when done. // status will show 'No Processes' if the list is empty
			pScene.SearchState.IsLoading = false
		}
	}(op.Context, pScene, op)

	return pScene
}

func (p *ProcessesScene) Update(d float64) {
	// update status | only if its true
	if p.SearchState.IsLoading {
		p.processes = p.processes[:0]
	}
}

func (p *ProcessesScene) Render(s interfaces.ScreenControl) {
	if p.Search != nil {
		p.Search.Render(s)
	}
	if p.kill != nil {
		p.kill.Render(s)
	}

	p.ScrollWindow.Render(p, s)

	s.Color(color.White)
	window.Text(s, screentui.P(1, 1), "Page: Running Processes") // title

	searchON := ""

	if p.SearchState.StrSearch != "" {
		searchON = "| [q] to reset search"
	}

	paddingBetweenText := 25
	startXPos := 32
	startYPos := 4
	window.Text(s,
		screentui.P(float64(startXPos), float64(startYPos)),
		fmt.Sprintf("Total Processes: %d | Displaying (%d - %d) %s", len(p.processes),
			p.ScrollWindow.Start, p.ScrollWindow.End, searchON,
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
		if p.SearchState.IsLoading {
			p.SearchState.Status = "Loading " + string(effects.Spinner(7, 150))
		} else {
			p.SearchState.Status = "No Processes"
		}
		window.Text(s,
			screentui.P(float64(startXPos), float64(startYPos+3)),
			p.SearchState.Status,
		)
		return
	}

	// sort processes by name default
	sortProcesses(p.sortedBy, p.processes)

	startProcessIndex := p.ScrollWindow.CurrentIndex + startYPos + 3

	// this will avoid index out of bound, triggered when hiding logs
	lastIndex := min(p.ScrollWindow.End, len(p.processes))

	// displaying list of processes
	for y, proc := range p.processes[p.ScrollWindow.Start:lastIndex] {
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
	if p.Search != nil {
		p.Search.Events(p, ev)
		return
	}

	// kill a process events
	if p.kill != nil {
		p.kill.Events(p, ev)
		return
	}

	// scroll window events
	p.ScrollWindow.Events(p, ev)

	// processes events
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Str() == "T" {
			p.sortedBy = nextSort(p.sortedBy)
			scenes.Log(fmt.Sprintf("Sorted By %s", p.sortedBy.String()))
		} else if ev.Str() == "t" {
			var order string
			if toggleDesc() {
				order = "Descending"
			} else {
				order = "Ascending"
			}
			scenes.Log("Order: " + order)
		} else if ev.Str() == "/" {
			p.MenuController.Lock()
			p.Search = &searchui.Search{}
			scenes.Log("Search ON")
		} else if ev.Str() == "K" {
			if len(p.processes) < 1 {
				scenes.Log("Cannot perform this action now.")
				return
			}
			p.MenuController.Lock()
			p.kill = &killui.Kill{
				Process: p.selectedProcess,
			}
			scenes.Log("Process Selected: " + p.selectedProcess.Name + " - Waiting action...")
		} else if ev.Str() == "q" {
			// clear / reset current invoked search
			if p.SearchState.StrSearch != "" {
				p.SearchState.IsLoading = true
				p.ScrollWindow.CurrentIndex = 0
				scenes.Log("Search cleared...")
				p.SearchState.StrSearch = ""
			}
			p.MenuController.Unlock()
		}
	}
}

func (p *ProcessesScene) CloseingKillUI() {
	p.kill = nil
	p.MenuController.Unlock()
	p.ScrollWindow.CurrentIndex = 0
}

func (p *ProcessesScene) ClosingSearchWith(search string) {
	if search == "" {
		if p.SearchState.StrSearch != "" {
			scenes.Log("Search Canceled / Reset")
			p.SearchState.IsLoading = true
		}
	}
	p.ScrollWindow.CurrentIndex = 0
	p.SearchState.StrSearch = search
	p.Search = nil
	p.MenuController.Unlock()
}

func (p *ProcessesScene) SetLoading() {
	p.SearchState.IsLoading = true
}
