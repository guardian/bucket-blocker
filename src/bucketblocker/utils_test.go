package main

import (
	"testing"
)

func TestSplitAndTrim(t *testing.T) {
	// Test case 1: Empty string should return an empty slice
	input1 := ""
	expected1 := 0
	result1 := len(splitAndTrim(input1))
	if result1 != expected1 {
		t.Errorf("Expected %v, but got %v", expected1, result1)
	}

	// Test case 2: Space separated strings should not be split, but should be trimmed
	input2 := "  hello  world  "
	expected2 := []string{"hello  world"}
	result2 := splitAndTrim(input2)
	if len(result2) != len(expected2) {
		t.Errorf("Expected %v, but got %v", expected2, result2)
	}

	// Test case 3: Comma separated strings should be split
	input3 := "a,b,c,d"
	expected3 := []string{"a", "b", "c", "d"}
	result3 := splitAndTrim(input3)
	if len(result3) != len(expected3) {
		t.Errorf("Expected %v, but got %v", expected3, result3)
	}
}

func TestContains(t *testing.T) {
	// Test case 1: Empty slice
	slice1 := []string{}
	item1 := "hello"
	expected1 := false
	result1 := contains(slice1, item1)
	if result1 != expected1 {
		t.Errorf("Expected %v, but got %v", expected1, result1)
	}

	// Test case 2: Slice with matching item
	slice2 := []string{"apple", "banana", "cherry"}
	item2 := "banana"
	expected2 := true
	result2 := contains(slice2, item2)
	if result2 != expected2 {
		t.Errorf("Expected %v, but got %v", expected2, result2)
	}

	// Test case 3: Slice without matching item
	slice3 := []string{"apple", "banana", "cherry"}
	item3 := "orange"
	expected3 := false
	result3 := contains(slice3, item3)
	if result3 != expected3 {
		t.Errorf("Expected %v, but got %v", expected3, result3)
	}
}
