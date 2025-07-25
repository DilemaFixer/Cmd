package parser

import (
	"fmt"
	"strings"
)

type ParserInput struct {
	Command     string
	Subcommands []string
	InputFlags  []InputFlag
}

type InputFlag struct {
	Name  string
	Value string
}

func NewParserInput(command string) *ParserInput {
	return &ParserInput{
		Command:     command,
		Subcommands: make([]string, 0),
		InputFlags:  make([]InputFlag, 0),
	}
}

func (f InputFlag) HaveInputValue() bool {
	return f.Value != ""
}

func ParseInput(input string) (*ParserInput, error) {
	input = strings.TrimSpace(input)
	if !validateInput(input) {
		return nil, fmt.Errorf("Parsing err: empty or only whitespace in string")
	}

	parts, partsCount := cutInput(input)
	if partsCount == 0 {
		return nil, fmt.Errorf("Parsing err: empty or only whitespace in string")
	}

	result := NewParserInput(parts[0])
	parts = removeFirst(parts, 1)

	if partsCount == 0 {
		return result, nil
	}

	var canBeSubcommand bool = true
	for _, str := range parts {
		if isFlag(str) {
			canBeSubcommand = false
			flag, err := parseFlag(str)
			if err != nil {
				return nil, err
			}
			result.InputFlags = append(result.InputFlags, flag)
		} else {
			if !canBeSubcommand {
				return nil, fmt.Errorf("Subcommand command %s can't go after flag", str)
			}
			result.Subcommands = append(result.Subcommands, str)
		}
	}

	return result, nil
}

func validateInput(input string) bool {
	return input != "" && input != " "
}

func cutInput(str string) ([]string, int) {
	parts := strings.Split(str, " ")
	return parts, len(parts)
}

func removeFirst[T any](parts []T, count int) []T {
	if partsLen := len(parts); partsLen == 0 || partsLen < count {
		return parts
	}
	return parts[count:]
}

func isFlag(str string) bool {
	return strings.HasPrefix(str, "--")
}

func parseFlag(str string) (InputFlag, error) {
	str = str[2:]
	var flag InputFlag
	if strings.Contains(str, "=") {
		strParts := strings.Split(str, "=")
		if strings.TrimSpace(strParts[1]) == "" {
			return flag, fmt.Errorf("Empty seter for %s", strParts[1])
		}
		flag = InputFlag{
			Name:  strParts[0],
			Value: strParts[1],
		}
	} else {
		flag = InputFlag{
			Name:  str,
			Value: "",
		}
	}
	return flag, nil
}
