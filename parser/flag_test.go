package parser

import (
	"reflect"
	"testing"
)

func makeFlagWithData() *InputFlag {
	return &InputFlag{
		Name:  "Name",
		Value: "Value",
	}
}

func makeEmptyFlag() *InputFlag {
	return &InputFlag{
		Name:  "",
		Value: "",
	}
}

func TestHaveInputValue_WithEmptyFlag_ReturnsFalse(t *testing.T) {
	flag := makeEmptyFlag()
	if flag.HaveInputValue() {
		t.Fatalf("Expected false for empty flag")
	}
}

func TestHaveInputValue_WithInputFlag_ReturnsTrue(t *testing.T) {
	flag := makeFlagWithData()
	if !flag.HaveInputValue() {
		t.Fatalf("Expected true for not empty input flag")
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
