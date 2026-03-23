package processes

import "sort"

// sorting processes
type sortBy int

const (
	sortByName sortBy = iota
	sortByPID
	sortByCPU
)

func sortProcesses(sortType sortBy, desc bool, processes []*process) {
	less := func(i, j int) bool {
		return processes[i].Name < processes[j].Name
	}

	switch sortType {
	case sortByPID:
		less = func(i, j int) bool {
			return processes[i].PID < processes[j].PID
		}
	case sortByCPU:
		less = func(i, j int) bool {
			return processes[i].CPUPercent < processes[j].CPUPercent
		}
	}

	sort.Slice(processes, func(i, j int) bool {
		if desc {
			return less(j, i) // flip
		}
		return less(i, j)
	})
}

func (s sortBy) String() string {
	switch s {
	case sortByPID:
		return "PID"
	case sortByCPU:
		return "CPU"
	default:
		return "Name"
	}
}
