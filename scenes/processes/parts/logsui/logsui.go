package logsui

const showLogsSize int = 22
const hideLogsSize int = 12

type LogsView interface {
	IsVisible() bool
	GetPaddingSize() int
	ToggleLogsView()
}

type logsViewState struct {
	showLogs       bool
	currentPadding int
}

func New() *logsViewState {
	return &logsViewState{
		showLogs:       true,
		currentPadding: showLogsSize,
	}

}

func (l *logsViewState) IsVisible() bool {
	return l.showLogs
}

func (l *logsViewState) GetPaddingSize() int {
	return l.currentPadding
}

func (l *logsViewState) ToggleLogsView() {
	l.showLogs = !l.showLogs
	if l.showLogs {
		l.currentPadding = showLogsSize
	} else {
		l.currentPadding = hideLogsSize
	}
}
