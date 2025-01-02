package log

// todo: later
type Log interface {
	Println(v ...any)
	Printf(format string, v ...any)
}
