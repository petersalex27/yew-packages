package errors

import (
	"github.com/petersalex27/yew-packages/str"
	"strconv"
)

type Locatable interface {
	str.Stringable
	SetPath(path string) Locatable
	SetLineChar(line, char int) Locatable
	GetLocation() (path string, line, char int)
}

func updateLine(loc Locatable, line int) Locatable {
	_, _, char := loc.GetLocation()
	return loc.SetLineChar(line, char)
} 

func updateChar(loc Locatable, char int) Locatable {
	_, line, _ := loc.GetLocation()
	return loc.SetLineChar(line, char)
} 

func updatePath(loc Locatable, path string) Locatable {
	return loc.SetPath(path)
} 

type empty_location struct {}

func (loc empty_location) GetLocation() (path string, line, char int) {
	return "", 0, 0
}

func (loc empty_location) SetPath(path string) Locatable {
	return Location{path: path, }
}

func (loc empty_location) SetLineChar(line, char int) Locatable {
	return Location{line: line, char: char,}
}

func (loc empty_location) String() string {
	return ""
}


type Location struct {
	path string
	line, char int
}

func (loc Location) GetLocation() (path string, line, char int) {
	return loc.path, loc.line, loc.char
}

func (loc Location) SetPath(path string) Locatable {
	return Location{path: path, line: loc.line, char: loc.char,}
}

func (loc Location) SetLineChar(line, char int) Locatable {
	return Location{path: loc.path, line: line, char: char,}
}

// "[path:line:char]", "[path:line:]", "[path::char]", "[:line:char]", 
// "[path::]", "[:line:]", "[::char]", ""
func (loc Location) String() string {
	line, char := strconv.Itoa(loc.line), strconv.Itoa(loc.char)
	path := loc.path
	if path == "" && loc.line == 0 && loc.char == 0 {
		return ""
	}

	path = "[" + path + ":"


	if loc.line == 0 {
		line = ""
	}
	line = line + ":"

	if loc.char == 0 {
		char = "" 
	}
	char = char + "]"

	return path + line + char
}