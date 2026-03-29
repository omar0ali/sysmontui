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
	mu sync.RWMutex

	LogsAddToList controls.LogsAddToList

	Processes    []*process
	scrollWindow scrollWindow

	sortedBy sortBy
	desc     bool

	search *search
	kill   *kill

	searchFor string

	mController interfaces.MenuController
}

func Init(logsFunc controls.LogsAddToList, ctx context.Context, op options.Options) *Processes {
	processes := &Processes{
		LogsAddToList: logsFunc,
		Processes:     []*process{},
		scrollWindow:  scrollWindow{},
		sortedBy:      sortByName,
	}

	if op.MenuController == nil {
		panic("Missing Controls Pause")
	} else {
		processes.mController = op.MenuController
	}

	processes.LogsAddToList("Reading processes...")

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
				processes.Processes = processes.Processes[:0] // reset list

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
					processes.Processes = append(processes.Processes, &process{name, pid, cpuPercent})
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

func (p *Processes) Update(d float64) {}

func (p *Processes) Render(s interfaces.ScreenControl) {
	if p.search != nil {
		p.search.Render(s)
	}
	if p.kill != nil {
		p.kill.Render(s)
	}

	// important: used for the scrollable area where the processes are displayed
	// must set both (end, max)
	p.scrollWindow.updateArea(p, s)

	s.Color(color.White)
	window.Text(s, screentui.P(1, 1), "Page: Running Processes") // title

	paddingBetweenText := 30
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
	sortProcesses(p.sortedBy, p.desc, p.Processes)
	for y, proc := range p.Processes[p.scrollWindow.start:p.scrollWindow.end] {
		window.ListOfTextsWithPadding(s, screentui.P(float64(startXPos), float64(startYPos+y+3)), paddingBetweenText, []string{
			proc.Name,
			fmt.Sprintf("%d", proc.PID),
			fmt.Sprintf("%.2f%%", proc.CPUPercent),
		})
	}
}

func (p *Processes) Events(ev tcell.Event) {
	// search processes events
	if p.search != nil {
		closeSearchWith := func(search string) {
			p.searchFor = search
			p.search = nil
			p.mController.Unlock()
		}
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Str() {
			case "/":
				p.LogsAddToList("Search Canceled / Reset")
				closeSearchWith("")
				return
			}
			switch ev.Key() {
			case tcell.KeyEnter:
				p.LogsAddToList("Searching for: " + p.search.text)
				closeSearchWith(p.search.text)
				return
			}
		}
		p.search.Events(ev)
		return
	}

	// kill a process events
	if p.kill != nil {
		close := func() {
			p.kill = nil
			p.mController.Unlock()
		}
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Str() {
			case "/":
				p.LogsAddToList("Canceled")
				close()
				return
			}
			switch ev.Key() {
			case tcell.KeyEnter:
				//TODO: kill process
				// p.LogsAddToList("SIGTERM: " + p.selectedProcess.Name)
				err := p.kill.SIGTERM()
				if err != nil {
					p.LogsAddToList(err.Error())
				}
				close()
				return
			}
		}
		return
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Str() == "j" {
			if p.scrollWindow.end < len(p.Processes) {
				p.scrollWindow.start++
				p.scrollWindow.end++
			}
		} else if ev.Str() == "k" {
			if p.scrollWindow.start > 0 {
				p.scrollWindow.start--
				p.scrollWindow.end--
			}
		} else if ev.Str() == "T" {
			p.sortedBy = NextSort(p.sortedBy)
			p.LogsAddToList(fmt.Sprintf("Sorted By %s", p.sortedBy.String()))
		} else if ev.Str() == "t" {
			p.desc = !p.desc
			order := "Ascending"
			if p.desc {
				order = "Descending"
			}
			p.LogsAddToList("Order: " + order)
		} else if ev.Str() == "/" {
			p.mController.Lock()
			p.search = &search{}
			p.LogsAddToList("Search ON")
		} else if ev.Str() == "K" {
			//TODO: kill process
		}
	}
}

func (sc *scrollWindow) updateArea(p *Processes, s interfaces.ScreenControl) {
	_, h := s.Size()
	// get the latest hight of the screen -
	// used for the scrollable area in processes page
	sc.max = h - PADDING_FOOTER

	// this fix a problem when .start is at a position that > than what it
	// its displayed. When searching, it might crash.
	if sc.start > len(p.Processes) {
		sc.start = 0
	}

	sc.end = sc.start + sc.max
	if len(p.Processes) < sc.max {
		sc.end = len(p.Processes)
	}
}
