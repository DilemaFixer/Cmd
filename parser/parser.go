package parser

import (
	"fmt"
	"os"
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

func ParseOSArgs() (*ParsedInput, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("Parsing err: empty args")
	}
	return ParseArgs(os.Args[1:])
}

func ParseArgs(args []string) (*ParsedInput, error) {
	parts := trimNonEmpty(args)
	if len(parts) == 0 {
		return nil, fmt.Errorf("Parsing err: empty args")
	}

	first := parts[0]
	if isFlag(first) {
		return nil, fmt.Errorf("First word must be command , not flag %s", first)
	}
	if first == "--" {
		return nil, fmt.Errorf("First word must be command , not flag --")
	}

	res := NewParserInput(first)

	parts = parts[1:]
	if len(parts) == 0 {
		return res, nil
	}

	flagsStarted := false
	for _, tok := range parts {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		if tok == "--" {
			return nil, fmt.Errorf("Invalid flag: --")
		}

		if isFlag(tok) {
			flagsStarted = true
			flag, err := parseFlag(tok)
			if err != nil {
				return nil, err
			}
			res.InputFlags = append(res.InputFlags, flag)
			continue
		}

		// не флаг → это сабкоманда
		if flagsStarted {
			return nil, fmt.Errorf("Subcommand command %s can't go after flag", tok)
		}
		res.Subcommands = append(res.Subcommands, tok)
	}

	return res, nil
}
