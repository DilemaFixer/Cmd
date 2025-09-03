package router

import (
	ctx "github.com/DilemaFixer/Cmd/context"
	prs "github.com/DilemaFixer/Cmd/parser"

	"testing"
)

func makeParserInput() *prs.ParsedInput {
	input := &prs.ParsedInput{
		Command:     "run",
		Subcommands: []string{"build", "deploy"},
		InputFlags: []prs.InputFlag{
			{Name: "count", Value: "10"},
			{Name: "rate", Value: "3.14"},
			{Name: "enabled", Value: "true"},
			{Name: "empty", Value: ""},
			{Name: "name", Value: "test"},
		},
	}

	return input
}

func makeContextWithData() *ctx.Context {
	return ctx.NewContext(makeParserInput())
}

func TestBuildRoutingPath_WithFullDataContext_ReturnRoutingPath(t *testing.T) {
	expectedCount := 3
	expectedPath := []string{
		"run",
		"build",
		"deploy",
	}

	context := makeContextWithData()
	str, count := buildRoutingPath(context)

	if count != expectedCount {
		t.Fatalf("expected len %d , get : %d", expectedCount, count)
	}

	for i, expectedPathPart := range expectedPath {
		if str[i] != expectedPathPart {
			t.Fatalf("expected path part %s at %d , but have %s", expectedPathPart, i, str[i])
		}
	}
}
