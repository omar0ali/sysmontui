# sysmontui
[Sysmon](https://github.com/omar0ali/sysmon) is a library I created for the [sysmontui](#) application. The goal is to build a terminal-based system monitor, similar to `htop`, but with a unique user interface.

### Features
- Read system metrics from `/proc` using [https://github.com/omar0ali/sysmon](sysmon) library.
- CPU usage information
- Memory statistics
- List of processes


### Checklist:
- [x] Showing CPU Status
- [x] Showing Memory Info
- [x] Showing list of processes
- [x] Support sorting processes by different fields (name, PID, etc.) and ordering (Ascending, Descending)
- [x] Search processes by name
- [x] Kill process
- [X] Show and Hide logs using `[l]` key

### Requirements
Linux only (uses /proc)

### Screenshots
#### CPU Status
![Sysmontui screenshot cpustatus](assets/screenshots/cpustatus.png)
#### Memory Info
![Sysmontui screenshot meminfo](assets/screenshots/meminfo.png)
#### Running Processes
![Sysmontui screenshot processes](assets/screenshots/processes.png)
#### Sorting By (name, pid, cpu usage) and ordering (asc, desc)
![Sysmontui screenshot processes sorting 1](assets/screenshots/sort(1).png)
![Sysmontui screenshot processes sorting 2](assets/screenshots/sort(2).png)
#### Search processes by name
![Sysmontui screenshot search](assets/screenshots/search.png)
#### Kill process - sending SIGTERM
![Sysmontui screenshot sigterm](assets/screenshots/sigterm.png)

## Installation

### From source
```bash
git clone https://github.com/omar0ali/sysmontui.git
cd sysmontui
go build -o build/sysmontui cmd/sysmontui/main.go

#Run
./build/sysmontui
```

### Using Go (recommended) - Linux only
```bash
go install github.com/omar0ali/sysmontui/cmd/sysmontui@latest
```

### Status
Work in progress

## Third-Party Licenses
This project uses [tcell](https://github.com/gdamore/tcell) (Apache License 2.0):
https://github.com/gdamore/tcell
