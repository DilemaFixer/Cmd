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

	parts := strings.Split(str, " ")
	return parts, len(parts)
}

func removeFirst[T any](parts []T, count int) []T {
	if partsLen := len(parts); partsLen == 0 || partsLen < count {
		return parts
	}
	return parts[count:]
}
