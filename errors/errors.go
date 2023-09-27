package errors

import (
	"strconv"
	"sync"
)

type Err struct {Notification}

var defaultErrLock sync.Mutex
var defaultErrorMessage string = "unknown error"

var pathMaxLineLenLock sync.Mutex
var pathMaxLineLen = make(map[string]int)

func getLineStrWidth(path string, fallback int) int {
	var lineLen int
	var found bool

	pathMaxLineLenLock.Lock()
	if path == "" {
		found = false
	} else {
		lineLen, found = pathMaxLineLen[path]
	}
	pathMaxLineLenLock.Unlock()

	var use int = lineLen
	if !found {
		if fallback < 0 {
			panic("fallback < 0")
		}
		use = fallback
	}
	return numWidth(use)
} 

func num(n int) string {
	return strconv.Itoa(n)
}

func numWidth(n int) int {
	return len(num(n))
}

// this sets a global variable
func SetDefaultError(new string) (old string) {
	defaultErrLock.Lock()
	defer defaultErrLock.Unlock()

	old = defaultErrorMessage
	defaultErrorMessage = new

	return
}

func New() Err {
	defaultErrLock.Lock()
	defer defaultErrLock.Unlock()

	e := notificationTemplate()
	e.msg = defaultErrorMessage
	return Err{e}
}

func Ferr(format string, args ...any) Err {
	return Err{Fnotify(format, args...)}
}

func (e Err) Error() string {
	return e.Notify("Error")
}