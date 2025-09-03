package parser

import (
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

func TestIsFlag_WithEmptyString_ReturnFalse(t *testing.T) {
	if isFlag("") {
		t.Fatalf("expected false on empty string")
	}
}

func TestIsFlag_WithSpaceOnlyString_ReturnFalse(t *testing.T) {
	if isFlag("       ") {
		t.Fatalf("expected false on space only string")
	}
}

func TestIsFlag_WithoutFlagPrefix_ReturnFalse(t *testing.T) {
	if isFlag("notFlag") {
		t.Fatalf("expected false on word without flag prefix")
	}
}

func TestIfFlag_WordWithFlagPrefix_ReturnTrue(t *testing.T) {
	if !isFlag("--isFlag") {
		t.Fatalf("expected true on word with flag prefix")
	}
}
