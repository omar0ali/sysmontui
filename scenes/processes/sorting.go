package processes

import "sort"

// sorting processes
type sortBy int

var (
	desc bool
)

const (
	sortByName sortBy = iota
	sortByPID
	sortByCPU
	sortByMEM
)

func toggleDesc() bool {
	desc = !desc
	return desc
}

func sortProcesses(sortType sortBy, processes []*process) {
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
	case sortByMEM:
		less = func(i, j int) bool {
			return processes[i].MEMUsage < processes[j].MEMUsage
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
	case sortByMEM:
		return "MEM"
	default:
		return "Name"
	}
}

func nextSort(by sortBy) sortBy {
	if by == sortByName {
		return sortByPID
	}
	if by == sortByPID {
		return sortByCPU
	}
	if by == sortByCPU {
		return sortByMEM
	}
	return sortByName
}
