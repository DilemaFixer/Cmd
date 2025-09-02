package context

import (
	"testing"

	prs "github.com/DilemaFixer/Cmd/parser"
)

func makeParserInput() *prs.ParserInput {
	return &prs.ParserInput{
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
}

func TestGetCommandAndComparison(t *testing.T) {
	ctx := NewContext(makeParserInput())

	if ctx.GetCommand() != "run" {
		t.Fatalf("expected 'run', got %s", ctx.GetCommand())
	}

	if !ctx.IsCommandEqual("run") {
		t.Fatalf("expected true for command equal")
	}

	if ctx.IsCommandEqual("stop") {
		t.Fatalf("expected false for wrong command")
	}
}

func TestSubcommands(t *testing.T) {
	ctx := NewContext(makeParserInput())

	if !ctx.IsSubcommandExist("build") {
		t.Fatalf("expected subcommand 'build' to exist")
	}

	if ctx.IsSubcommandExist("clean") {
		t.Fatalf("expected subcommand 'clean' not to exist")
	}

	subcommands := ctx.GetSubcommandsAsArr()
	if len(subcommands) != 2 {
		t.Fatalf("expected 2 subcommands, got %d", len(subcommands))
	}
}

func TestFlagExistence(t *testing.T) {
	ctx := NewContext(makeParserInput())

	if !ctx.IsFlagExist("count") {
		t.Fatalf("expected flag 'count' to exist")
	}

	if ctx.IsFlagExist("missing") {
		t.Fatalf("expected flag 'missing' not to exist")
	}

	if !ctx.IsFlagHaveValue("count") {
		t.Fatalf("expected flag 'count' to have value")
	}

	if ctx.IsFlagHaveValue("empty") {
		t.Fatalf("expected flag 'empty' to have no value")
	}
}

func TestGetValueAsInt32(t *testing.T) {
	ctx := NewContext(makeParserInput())

	val, err := ctx.GetValueAsInt32("count")
	if err != nil || val != 10 {
		t.Fatalf("expected 10, got %d, err=%v", val, err)
	}

	_, err = ctx.GetValueAsInt32("missing")
	if err == nil {
		t.Fatalf("expected error for missing flag")
	}

	_, err = ctx.GetValueAsInt32("empty")
	if err == nil {
		t.Fatalf("expected error for empty flag")
	}
}

func TestGetValueAsInt64AndInt(t *testing.T) {
	ctx := NewContext(makeParserInput())

	val64, err := ctx.GetValueAsInt64("count")
	if err != nil || val64 != 10 {
		t.Fatalf("expected 10, got %d, err=%v", val64, err)
	}

	val, err := ctx.GetValueAsInt("count")
	if err != nil || val != 10 {
		t.Fatalf("expected 10, got %d, err=%v", val, err)
	}
}

func TestGetValueAsFloat32AndFloat64(t *testing.T) {
	ctx := NewContext(makeParserInput())

	val32, err := ctx.GetValueAsFloat32("rate")
	if err != nil || val32 != 3.14 {
		t.Fatalf("expected 3.14, got %f, err=%v", val32, err)
	}

	val64, err := ctx.GetValueAsFloat64("rate")
	if err != nil || val64 != 3.14 {
		t.Fatalf("expected 3.14, got %f, err=%v", val64, err)
	}
}

func TestGetValueAsBool(t *testing.T) {
	ctx := NewContext(makeParserInput())

	val, err := ctx.GetValueAsBool("enabled")
	if err != nil || val != true {
		t.Fatalf("expected true, got %v, err=%v", val, err)
	}

	_, err = ctx.GetValueAsBool("empty")
	if err == nil {
		t.Fatalf("expected error for empty flag")
	}
}

func TestGetValueAsString(t *testing.T) {
	ctx := NewContext(makeParserInput())

	val, err := ctx.GetValueAsString("name")
	if err != nil || val != "test" {
		t.Fatalf("expected 'test', got %s, err=%v", val, err)
	}

	_, err = ctx.GetValueAsString("missing")
	if err == nil {
		t.Fatalf("expected error for missing flag")
	}
}

func TestGetValueOrDefault(t *testing.T) {
	ctx := NewContext(makeParserInput())

	val := ctx.GetValueOrDefault("name", "default")
	if val != "test" {
		t.Fatalf("expected 'test', got %s", val)
	}

	val = ctx.GetValueOrDefault("missing", "default")
	if val != "default" {
		t.Fatalf("expected 'default', got %s", val)
	}

	val = ctx.GetValueOrDefault("empty", "default")
	if val != "default" {
		t.Fatalf("expected 'default', got %s", val)
	}
}

func TestFlagsCollections(t *testing.T) {
	ctx := NewContext(makeParserInput())

	flagsMap := ctx.GetFlagsAsMap()
	if len(flagsMap) != 5 {
		t.Fatalf("expected 5 flags, got %d", len(flagsMap))
	}

	flagsArr := ctx.GetFlagsAsArr()
	if len(flagsArr) != 5 {
		t.Fatalf("expected 5 flags array, got %d", len(flagsArr))
	}

	keysArr := ctx.GetFlagsKeysAsArr()
	if len(keysArr) != 5 {
		t.Fatalf("expected 5 keys, got %d", len(keysArr))
	}

	valuesArr := ctx.GetFlagsValuesAsArr()
	if len(valuesArr) != 5 {
		t.Fatalf("expected 5 values, got %d", len(valuesArr))
	}
}
