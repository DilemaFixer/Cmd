package parser 

import(
	"strings"
	"fmt"
)

type ParserInput struct {
	Command string
	Subcommands []string
	InputFlags []InputFlag
}

type InputFlag struct {
	Name string 
	Value string
}

func (f InputFlag)HaveInputValue() bool {
	return f.Value != ""
}

func ParseInput(input string) (*ParserInput, error) {
	input = strings.TrimSpace(input)
	if !validateInput(input) {
		return nil, fmt.Errorf("Parsing err: empty or only whitespace in string")
	}
		
	parts := strings.Split(input, " ")
	partsCount := len(parts)

	if partsCount == 0 {
		return nil, fmt.Errorf("Parsing err: empty or only whitespace in string")
	}
	
	result := &ParserInput {
		Command: parts[0],
		Subcommands: make([]string, 0),
		InputFlags: make([]InputFlag, 0),
	}

	parts = parts[1:]
	partsCount--

	if partsCount != 0 {
		var canBeSubcommand bool = true
		for _, str := range parts {
			if strings.HasPrefix(str, "--"){
				canBeSubcommand = false
				str = str[2:]
				var flag InputFlag
				if strings.Contains(str, "=") {
					strParts := strings.Split(str, "=")
					println(strParts[1])
					if strings.TrimSpace(strParts[1]) == "" {
						return nil, fmt.Errorf("Empty seter for %s", strParts[1])
					}
					flag = InputFlag {
						Name: strParts[0],
						Value: strParts[1],
					}
				} else {
					flag = InputFlag {
						Name: str,
						Value: "",
					}
				}
				
				result.InputFlags = append(result.InputFlags, flag)
				continue
			}

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

