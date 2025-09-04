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

	if isFlag(parts[0]) {
		return nil, fmt.Errorf("First word must be command , not flag %s", parts[0])
	}

	result := NewParserInput(parts[0])
	parts = removeFirst(parts, 1)

	if partsCount == 0 {
		return result, nil
	}

	var flagsQueueStarted bool = false
	for _, str := range parts {
		if isFlag(str) {
			flagsQueueStarted = true
			flag, err := parseFlag(str)
			if err != nil {
				return nil, err
			}
			result.InputFlags = append(result.InputFlags, flag)
		} else {
			if flagsQueueStarted {
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
		strParts := strings.SplitN(str, "=", 2)
		name := strings.TrimSpace(strParts[0])
		value := strings.TrimSpace(strParts[1])

		if value == "" {
			return flag, fmt.Errorf("Empty setter for %s", name)
		}

		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		flag = InputFlag{
			Name:  name,
			Value: value,
		}
	} else {
		flag = InputFlag{
			Name:  str,
			Value: "",
		}
	}

	return flag, nil
}
