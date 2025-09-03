package parser

import (
	"strings"
)

type InputFlag struct {
	Name  string
	Value string
}

func (f InputFlag) HaveInputValue() bool {
	return f.Value != ""
}

func isFlag(str string) bool {
	return strings.HasPrefix(str, "--")
}
