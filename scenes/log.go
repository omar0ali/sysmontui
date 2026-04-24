package scenes

var logFunc = func(string) {} // default no-op (never nil)

func SetLogger(f func(string)) {
	if f != nil {
		logFunc = f
	}
}

func Log(msg string) {
	logFunc(msg)
}
