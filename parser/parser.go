package parser

import (
	"fmt"
	"strings"
)

func ParseInput(input string) (*ParsedInput, error) {
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

func parseFlag(str string) (InputFlag, error) {
	str = strings.TrimSpace(str)
	if str == "" {
		return InputFlag{}, fmt.Errorf("empty input string")
	}

	str = strings.TrimPrefix(str, "--")
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
