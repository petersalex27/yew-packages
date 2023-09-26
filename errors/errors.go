package errors

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/petersalex27/yew-packages/util"
	"github.com/fatih/color"
)

type Err struct {
	Locatable
	subtype string
	msg string
	source string
	ptr Pointable
}

var defaultErrLock sync.Mutex
var defaultErrorMessage string = "unknown error"

var pathMaxLineLenLock sync.Mutex
var pathMaxLineLen = make(map[string]int)
func lineInPath(path string, line int) {
	if path == "" || line <= 0 {
		return
	}

	pathMaxLineLenLock.Lock()
	defer pathMaxLineLenLock.Unlock()

	lineLen, found := pathMaxLineLen[path]
	if !found {
		pathMaxLineLen[path] = line
	} else {
		pathMaxLineLen[path] = util.Max(line, lineLen)
	}
}

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

func (e Err) SetLocation(path string, line, char int) Err {
	e.Locatable = e.Locatable.SetPath(path).SetLineChar(line, char)
	return e
}

func (e Err) SetSubtype(ty string) Err {
	e.subtype = ty
	return e
}

func (e Err) subTypeString() string {
	if e.subtype == "" {
		return ""
	}
	return color.MagentaString(" (%s)", e.subtype)
}

var specTable = map[rune]func(*Err, any, int){
	'l': func(e *Err, line any, position int) {
		res, ok := line.(int)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier l")
		}
		e.Locatable = updateLine(e.Locatable, res)
	},
	'c': func(e *Err, char any, position int) {
		res, ok := char.(int)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier c")
		}
		e.Locatable = updateChar(e.Locatable, res)
	},
	's': func(e *Err, source any, position int) {
		res, ok := source.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier s")
		}
		e.source = res
	},
	't': func(e *Err, subtype any, position int) {
		res, ok := subtype.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier t")
		}
		e.subtype = res
	},
	'p': func(e *Err, ptrMsg any, position int) {
		res, ok := ptrMsg.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier p")
		}
		if e.ptr == nil {
			e.ptr = Pointer{}
		}
		e.ptr = e.ptr.setTail(res)
	},
	'f': func(e *Err, path any, position int) {
		res, ok := path.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier f")
		}
		e.Locatable = e.SetPath(res)
	},
	'r': func(e *Err, rng any, position int) {
		res, ok := rng.(int)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier f")
		}
		if e.ptr == nil {
			e.ptr = PointerRange{rngLen: res}
		} else {
			e.ptr = PointerRange{rngLen: res, pointer_shared: e.ptr.getShared()}
		}
	},
	'm': func(e *Err, msg any, position int) {
		res, ok := msg.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier m")
		}
		e.msg = res
	},
}

func Ferr(format string, args ...any) Err {
	e := New()
	format = strings.TrimSpace(format)
	if len(args) < len(format) {
		panic("foramt has more specifiers than arguments")
	}

	for i, f := range format {
		setFunc, found := specTable[f]
		if !found {
			panic("illegal format specifier " + string(f))
		}
		setFunc(&e, args[i], i+1)
	}
	return e
}

func New() Err {
	defaultErrLock.Lock()
	defer defaultErrLock.Unlock()

	return Err{
		Locatable: empty_location{},
		subtype: "",
		msg: defaultErrorMessage,
		ptr: nil,
		source: "",
	}
}

func (e Err) GetTail() string {
	if e.source == "" {
		return ""
	}

	path, line, char := e.Locatable.GetLocation()
	var tail string = "  "
	width := getLineStrWidth(path, line)
	padding, _ := (pointer_shared{paddingLeft: width-line}).Strings()
	tail = tail + padding + num(line) + " | "
	if e.ptr == nil {
		e.ptr = Pointer{}
	}

	if char < 0 {
		char = 0
	}
	padLen := len(tail)+(char-1)
	return "\n" + tail + e.source + "\n" + e.ptr.setPadding(padLen).String()
}
// [path:line:char] Error (Subtype): error message here
//   line | source code here
//                 ^^^^ tail message
func (e Err) Error() string {
	loc := e.Locatable.String()
	
	return fmt.Sprintf("%s%s%s%s", 
			color.RedString("%s Error", loc),
			e.subTypeString(), 
			color.RedString(": %s", e.msg),
			e.GetTail())
}

