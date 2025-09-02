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

func makeBrokenParserInput() *prs.ParserInput {
	return &prs.ParserInput{
		Command:     "",
		Subcommands: []string{},
		InputFlags: []prs.InputFlag{
			{Name: "badInt", Value: "notanumber"},
			{Name: "badFloat", Value: "notfloat"},
			{Name: "badBool", Value: "maybe"},
			{Name: "empty", Value: ""},
		},
	}
}

// --- Command tests ---

func TestGetCommand_ReturnsCommand(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if ctx.GetCommand() != "run" {
		t.Fatalf("expected 'run', got %s", ctx.GetCommand())
	}
}

func TestGetCommand_ReturnsEmpty_OnBrokenInput(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if ctx.GetCommand() != "" {
		t.Fatalf("expected empty command, got %s", ctx.GetCommand())
	}
}

func TestIsCommandEqual_MatchesCorrectly(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if !ctx.IsCommandEqual("run") {
		t.Fatalf("expected true for correct command")
	}
	if ctx.IsCommandEqual("stop") {
		t.Fatalf("expected false for wrong command")
	}
}

func TestIsCommandEqual_ReturnsFalse_OnEmptyCommand(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if ctx.IsCommandEqual("run") {
		t.Fatalf("expected false for IsCommandEqual on empty command")
	}
}

// --- Subcommand tests ---

func TestIsSubcommandExist_FindsExisting(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if !ctx.IsSubcommandExist("build") {
		t.Fatalf("expected subcommand 'build' to exist")
	}
}

func TestIsSubcommandExist_ReturnsFalseForMissing(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if ctx.IsSubcommandExist("clean") {
		t.Fatalf("expected subcommand 'clean' not to exist")
	}
}

func TestGetSubcommandsAsArr_ReturnsAll(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if len(ctx.GetSubcommandsAsArr()) != 2 {
		t.Fatalf("expected 2 subcommands, got %d", len(ctx.GetSubcommandsAsArr()))
	}
}

func TestGetSubcommandsAsArr_EmptyOnBrokenInput(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if len(ctx.GetSubcommandsAsArr()) != 0 {
		t.Fatalf("expected no subcommands")
	}
}

// --- Flags existence ---

func TestIsFlagExist_FindsFlag(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if !ctx.IsFlagExist("count") {
		t.Fatalf("expected flag 'count' to exist")
	}
}

func TestIsFlagExist_ReturnsFalseForMissing(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if ctx.IsFlagExist("missing") {
		t.Fatalf("expected missing flag not to exist")
	}
}

func TestIsFlagHaveValue_ReturnsTrueForFilled(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if !ctx.IsFlagHaveValue("count") {
		t.Fatalf("expected flag 'count' to have value")
	}
}

func TestIsFlagHaveValue_ReturnsFalseForEmpty(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if ctx.IsFlagHaveValue("empty") {
		t.Fatalf("expected flag 'empty' to have no value")
	}
}

func TestIsFlagHaveValue_ReturnsFalseForMissing(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if ctx.IsFlagHaveValue("missing") {
		t.Fatalf("expected missing flag not to have value")
	}
}

// --- Value conversions ---

func TestGetValueAsInt32_ParsesCorrectly(t *testing.T) {
	ctx := NewContext(makeParserInput())
	val, err := ctx.GetValueAsInt32("count")
	if err != nil || val != 10 {
		t.Fatalf("expected 10, got %d, err=%v", val, err)
	}
}

func TestGetValueAsInt32_ErrorsOnMissingOrEmpty(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if _, err := ctx.GetValueAsInt32("missing"); err == nil {
		t.Fatalf("expected error for missing flag")
	}
	if _, err := ctx.GetValueAsInt32("empty"); err == nil {
		t.Fatalf("expected error for empty flag")
	}
}

func TestGetValueAsIntAndInt64_ParsesCorrectly(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if val, err := ctx.GetValueAsInt("count"); err != nil || val != 10 {
		t.Fatalf("expected 10, got %d, err=%v", val, err)
	}
	if val64, err := ctx.GetValueAsInt64("count"); err != nil || val64 != 10 {
		t.Fatalf("expected 10, got %d, err=%v", val64, err)
	}
}

func TestGetValueAsIntAndInt64_ErrorsOnInvalid(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if _, err := ctx.GetValueAsInt("badInt"); err == nil {
		t.Fatalf("expected error for invalid int")
	}
	if _, err := ctx.GetValueAsInt64("badInt"); err == nil {
		t.Fatalf("expected error for invalid int64")
	}
}

func TestGetValueAsFloat32AndFloat64_ParsesCorrectly(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if val32, err := ctx.GetValueAsFloat32("rate"); err != nil || val32 != 3.14 {
		t.Fatalf("expected 3.14, got %f, err=%v", val32, err)
	}
	if val64, err := ctx.GetValueAsFloat64("rate"); err != nil || val64 != 3.14 {
		t.Fatalf("expected 3.14, got %f, err=%v", val64, err)
	}
}

func TestGetValueAsFloat32AndFloat64_ErrorsOnInvalid(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if _, err := ctx.GetValueAsFloat32("badFloat"); err == nil {
		t.Fatalf("expected error for invalid float32")
	}
	if _, err := ctx.GetValueAsFloat64("badFloat"); err == nil {
		t.Fatalf("expected error for invalid float64")
	}
}

func TestGetValueAsBool_ParsesCorrectly(t *testing.T) {
	ctx := NewContext(makeParserInput())
	val, err := ctx.GetValueAsBool("enabled")
	if err != nil || val != true {
		t.Fatalf("expected true, got %v, err=%v", val, err)
	}
}

func TestGetValueAsBool_ErrorsOnInvalidOrEmpty(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if _, err := ctx.GetValueAsBool("empty"); err == nil {
		t.Fatalf("expected error for empty flag")
	}
	if _, err := ctx.GetValueAsBool("badBool"); err == nil {
		t.Fatalf("expected error for invalid bool")
	}
}

func TestGetValueAsString_ReturnsValue(t *testing.T) {
	ctx := NewContext(makeParserInput())
	val, err := ctx.GetValueAsString("name")
	if err != nil || val != "test" {
		t.Fatalf("expected 'test', got %s, err=%v", val, err)
	}
}

func TestGetValueAsString_ErrorsOnMissing(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if _, err := ctx.GetValueAsString("missing"); err == nil {
		t.Fatalf("expected error for missing flag")
	}
}

func TestGetValueOrDefault_ReturnsValueOrFallback(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if ctx.GetValueOrDefault("name", "default") != "test" {
		t.Fatalf("expected 'test'")
	}
	if ctx.GetValueOrDefault("missing", "default") != "default" {
		t.Fatalf("expected 'default'")
	}
	if ctx.GetValueOrDefault("empty", "default") != "default" {
		t.Fatalf("expected 'default'")
	}
}

func TestGetValueOrDefault_ReturnsFallbackOnBroken(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if ctx.GetValueOrDefault("missing", "fallback") != "fallback" {
		t.Fatalf("expected 'fallback'")
	}
}

// --- Collections ---

func TestFlagsCollections_ReturnsAllFlags(t *testing.T) {
	ctx := NewContext(makeParserInput())
	if len(ctx.GetFlagsAsMap()) != 5 {
		t.Fatalf("expected 5 flags in map")
	}
	if len(ctx.GetFlagsAsArr()) != 5 {
		t.Fatalf("expected 5 flags in array")
	}
	if len(ctx.GetFlagsKeysAsArr()) != 5 {
		t.Fatalf("expected 5 keys")
	}
	if len(ctx.GetFlagsValuesAsArr()) != 5 {
		t.Fatalf("expected 5 values")
	}
}

func TestFlagsCollections_ReturnsBrokenFlags(t *testing.T) {
	ctx := NewContext(makeBrokenParserInput())
	if len(ctx.GetFlagsAsMap()) != 4 {
		t.Fatalf("expected 4 broken flags in map")
	}
	if len(ctx.GetFlagsKeysAsArr()) != 4 {
		t.Fatalf("expected 4 keys in broken input")
	}
}
