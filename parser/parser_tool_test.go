package parser

import (
	"reflect"
	"testing"
)

// --- ValidateInput ---

func TestValidateInput_WithEmptyString_ReturnFalse(t *testing.T) {
	if validateInput("") {
		t.Fatalf("expected false for empty string")
	}
}

func TestValidateInput_WithSpaceOnlyString_ReturnFalse(t *testing.T) {
	if validateInput("      ") {
		t.Fatalf("expected false for space only string")
	}
}

func TestValidateInput_WithValidString_ReturnTrue(t *testing.T) {
	if !validateInput("test") {
		t.Fatalf("expected true for 'test' valid string")
	}
}

// --- CutInput ---

func TestCutInput_WithEmptyString_ReturnEmptyStringAndZero(t *testing.T) {
	if value, leng := cutInput(""); len(value) != 0 && leng != 0 {
		t.Fatalf("expected nil arr and zero leng for empty string, but have arr len:%d | returned len:%d", len(value), leng)
	}
}

func TestCutInput_WithSpaceOnlyString_ReturnEmptyStringAndZero(t *testing.T) {
	if value, leng := cutInput("    "); len(value) != 0 && leng != 0 {
		t.Fatalf("expected nil arr and zero leng for space only string, but have arr len:%d | returned len:%d", len(value), leng)
	}
}

func TestCutInput_WithValidString_ReturnArrWithOneItem(t *testing.T) {
	var value []string
	var leng int

	if value, leng = cutInput("value"); len(value) != 1 && leng != 1 {
		t.Fatalf("expected arr with one item and 1 leng for valid string, but have arr len:%d | returned len:%d", len(value), leng)
	}

	if value[0] != "value" {
		t.Fatalf("expected arr with item 'value' , but have %s", value[0])
	}
}

func TestCutInput_WithValidString_ReturnArrWithTwoItems(t *testing.T) {
	expected := []string{
		"value1",
		"value2",
	}

	var values []string
	var leng int

	if values, leng = cutInput("value1 value2"); len(values) != 2 && leng != 2 {
		t.Fatalf("expected arr with one item and 2 leng for valid string, but have arr len:%d | returned len:%d", len(values), leng)
	}

	for i, value := range expected {
		if values[i] != value {
			t.Fatalf("expected item %s at position %d , but have %s", value, i, values[i])
		}
	}
}

// --- RemoveFirst ---

func TestRemoveFirst_WithEmptySlice_ReturnsEmptySlice(t *testing.T) {
	input := []int{}
	got := removeFirst(input, 1)
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestRemoveFirst_WithCountGreaterThanLength_ReturnsOriginalSlice(t *testing.T) {
	input := []int{1, 2, 3}
	got := removeFirst(input, 5)
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestRemoveFirst_WithCountEqualLength_ReturnsEmptySlice(t *testing.T) {
	input := []int{1, 2, 3}
	got := removeFirst(input, 3)
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestRemoveFirst_WithCountLessThanLength_ReturnsRemainingSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := removeFirst(input, 2)
	want := []int{3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
