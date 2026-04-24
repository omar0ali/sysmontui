package processes

type Process struct {
	Name       string
	PID        int
	CPUPercent float64
	MEMUsage   uint64
}
