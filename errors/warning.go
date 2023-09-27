package errors

type Warning struct{ Notification }

func Fwarn(format string, args ...any) Warning {
	return Warning{Fnotify(format, args...)}
}

func (w Warning) Warn() string {
	return w.Notify("Warning")
}
