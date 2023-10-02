package errors

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/petersalex27/yew-packages/errors/notification"
)

type Notification struct {
	notification.Locatable
	subtype subtype
	msg string
	source string
	ptr Pointable
}

func (n Notification) GetTail() string {
	if n.source == "" {
		return ""
	}

	path, line, char := n.Locatable.GetLocation()
	var tail string = "  "
	width := getLineStrWidth(path, line)
	padding, _ := (pointer_shared{paddingLeft: width-line}).Strings()
	tail = tail + padding + num(line) + " | "
	if n.ptr == nil {
		n.ptr = Pointer{}
	}

	if char < 0 {
		char = 0
	}
	padLen := len(tail)+(char-1)
	return "\n" + tail + n.source + "\n" + n.ptr.setPadding(padLen).String()
}

func notificationTemplate() Notification {
	return Notification{
		Locatable: notification.EmptyLocation(),
		subtype: "",
		ptr: nil,
		source: "",
	}
}


var specTable = map[rune]func(*Notification, any, int){
	'l': func(e *Notification, line any, position int) {
		res, ok := line.(int)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier l")
		}
		e.Locatable = notification.UpdateLine(e.Locatable, res)
	},
	'c': func(e *Notification, char any, position int) {
		res, ok := char.(int)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier c")
		}
		e.Locatable = notification.UpdateChar(e.Locatable, res)
	},
	's': func(e *Notification, source any, position int) {
		res, ok := source.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier s")
		}
		e.source = res
	},
	't': func(e *Notification, subtype_ any, position int) {
		res, ok := subtype_.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier t")
		}
		e.subtype = subtype(res)
	},
	'p': func(e *Notification, ptrMsg any, position int) {
		res, ok := ptrMsg.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier p")
		}
		if e.ptr == nil {
			e.ptr = Pointer{}
		}
		e.ptr = e.ptr.setTail(res)
	},
	'f': func(e *Notification, path any, position int) {
		res, ok := path.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier f")
		}
		e.Locatable = e.SetPath(res)
	},
	'r': func(e *Notification, rng any, position int) {
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
	'm': func(e *Notification, msg any, position int) {
		res, ok := msg.(string)
		if !ok {
			panic("argument " + num(position) + " does not match the format specifier m")
		}
		e.msg = res
	},
}

func Fnotify(format string, args ...any) Notification {
	e := New().Notification
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

func (n Notification) Notify(of string) string {
	loc := ""
	if n.Locatable != nil {
		loc = n.Locatable.String()
	}
	
	if loc != "" {
		loc = loc + " "
	}
	
	return fmt.Sprintf("%s%s%s%s", 
			color.RedString("%s%s", loc, of),
			n.subtype.subTypeString(), 
			color.RedString(": %s", n.msg),
			n.GetTail())
}