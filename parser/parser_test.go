package parser

import (
	"reflect"
	"testing"
)

func TestParseInput_WithCommandSubcommandsAndAllKindOfFlags_ReturnsParsedInput(t *testing.T) {
	input := "command subcommand second_subcommand --bool_flag --int_flag=10 --string_flag='hello world'"
	parsedInput, err := ParseInput(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := ParsedInput{
		Command:     "command",
		Subcommands: []string{"subcommand", "second_subcommand"},
		InputFlags: []InputFlag{
			{Name: "bool_flag", Value: ""},
			{Name: "int_flag", Value: "10"},
			{Name: "string_flag", Value: "hello world"},
		},
	}

	if !reflect.DeepEqual(*parsedInput, want) {
		t.Fatalf("expected %v, got %v", want, parsedInput)
	}
}

func TestParseInput_WithEmptyString_ReturnError(t *testing.T) {
	_, err := ParseInput("")

	if err == nil {
		t.Fatalf("expected parsing error")
	}

	if err.Error() != "Parsing err: empty or only whitespace in string" {
		t.Fatalf("expected error msg 'Parsing err: empty or only whitespace in string' , got %s", err.Error())
	}
}

func TestParseInput_WithSpaceOnlyString_ReturnError(t *testing.T) {
	_, err := ParseInput("   ")

	if err == nil {
		t.Fatalf("expected parsing error")
	}

	if err.Error() != "Parsing err: empty or only whitespace in string" {
		t.Fatalf("expected error msg 'Parsing err: empty or only whitespace in string' , got %s", err.Error())
	}
}

func TestParseInput_WithOnlyCommand_ReturnCommandOnly(t *testing.T) {
	parsedInput, err := ParseInput("command")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	want := ParsedInput{
		Command:     "command",
		Subcommands: []string{},
		InputFlags:  []InputFlag{},
	}

	if !reflect.DeepEqual(*parsedInput, want) {
		t.Fatalf("expected %v, got %v", want, parsedInput)
	}
}

func TestParseInput_WithFlagOnly_ReturnError(t *testing.T) {
	_, err := ParseInput("--flag")
	if err == nil {
		t.Fatalf("expected error")
	}

	if err.Error() != "First word must be command , not flag --flag" {
		t.Fatalf("unexpected error , expect '%s' , get '%s'", "First word must be command , not flag --flag", err.Error())
	}
}

func TestParseInput_WithSubcommandAfterFlag_ReturnError(t *testing.T) {
	_, err := ParseInput("command --flag sucommand")
	if err == nil {
		t.Fatalf("expected error")
	}

	if err.Error() != "Subcommand command sucommand can't go after flag" {
		t.Fatalf("unexpected error, expect '%s', get '%s'", "Subcommand command sucommand can't go after flag", err.Error())
	}
}

func TestParseFlag_WithEmptyString_ReturnError(t *testing.T) {
	if _, err := parseFlag(""); err != nil && err.Error() != "empty input string" {
		t.Fatalf("unexpected error , get error: %s", err.Error())
	}
}

func TestParseFlag_WithEmptyString_ReturnsError(t *testing.T) {
	_, err := parseFlag("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseFlag_WithOnlyName_ReturnsFlagWithEmptyValue(t *testing.T) {
	got, err := parseFlag("--verbose")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := InputFlag{Name: "verbose", Value: ""}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestParseFlag_WithNameAndValue_ReturnsFlagWithValue(t *testing.T) {
	got, err := parseFlag("--threads=4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := InputFlag{Name: "threads", Value: "4"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestParseFlag_WithNameAndEmptyValue_ReturnsError(t *testing.T) {
	_, err := parseFlag("--mode=")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseFlag_WithoutPrefix_ReturnsFlag(t *testing.T) {
	got, err := parseFlag("count=10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := InputFlag{Name: "count", Value: "10"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
