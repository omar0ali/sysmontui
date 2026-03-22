# sysmontui
[Sysmon](https://github.com/omar0ali/sysmon) is a library I created for the [sysmontui](#) application. The goal is to build a terminal-based system monitor, similar to `htop`, but with a unique user interface.

### Features
- Read system metrics from `/proc` using [https://github.com/omar0ali/sysmon](sysmon) library.
- CPU usage information
- Memory statistics
- List of processes


### Checklist:
- [x] Home Page
- [x] CPU Status Info Page
- [x] Memory Info Page
- [x] Showing list of processes
- [ ] Support sorting processes by different fields (name, PID, etc.)
- [ ] Search processes
- [ ] UI Polishing & Clean up

### Requirements
Linux only (uses /proc)

### Screenshots
#### CPU Status
![Sysmontui screenshot cpustatus](assets/screenshots/cpustatus.png)
#### Memory Info
![Sysmontui screenshot meminfo](assets/screenshots/meminfo.png)
#### Running Processes
![Sysmontui screenshot processes](assets/screenshots/processes.png)

## Installation

### From source
```bash
git clone https://github.com/omar0ali/sysmontui.git
cd sysmontui
go build -o build/sysmontui cmd/cli/main.go
```

### Using Go (recommended) - Linux only
```bash
go install github.com/omar0ali/sysmontui/cmd/sysmontui@latest
```

### Status
Work in progress
