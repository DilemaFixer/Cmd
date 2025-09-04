package parser

import (
	"strings"
)

func validateInput(input string) bool {
	input = strings.TrimSpace(input)
	return input != ""
}

func cutInput(str string) ([]string, int) {
	if strings.TrimSpace(str) == "" {
		return nil, 0
	}

	var parts []string
	var buf strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false

	for _, r := range str {
		switch r {
		case ' ':
			if inSingleQuotes || inDoubleQuotes {
				buf.WriteRune(r)
			} else {
				if buf.Len() > 0 {
					parts = append(parts, buf.String())
					buf.Reset()
				}
			}
		case '"':
			inDoubleQuotes = !inDoubleQuotes
		case '\'':
			inSingleQuotes = !inSingleQuotes
		default:
			buf.WriteRune(r)
		}
	}

	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}

	return parts, len(parts)
}

func removeFirst[T any](parts []T, count int) []T {
	if partsLen := len(parts); partsLen == 0 || partsLen < count {
		return parts
	}
	return parts[count:]
}
